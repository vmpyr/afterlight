package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/vmpyr/afterlight/internal/store"
)

type AuthHandler struct {
	store *store.Store
}

func NewAuthHandler(s *store.Store) *AuthHandler {
	return &AuthHandler{store: s}
}

func (h *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()

	// Public routes
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(h.AuthMiddleware)
		r.Get("/me", h.GetCurrentUser)
		r.Post("/logout", h.Logout)
	})

	return r
}

// Handlers
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := IsValidPassword(req.Password); err != nil {
		http.Error(w, "Password does not meet complexity requirements: "+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.store.CreateUserTx(r.Context(), store.RegisterUserParams{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		// TODO: Handle specific errors (e.g., duplicate email)
		http.Error(w, "Registration failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(UserResponse{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		CurrentStatus: user.CurrentStatus,
		CreatedAt:     user.CreatedAt,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.store.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User does not exist", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	match, err := argon2id.ComparePasswordAndHash(req.Password, user.PasswordHash)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !match {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	// Cookies / Sessions
	if oldCookie, err := r.Cookie("session_token"); err == nil {
		_ = h.store.DeleteSession(r.Context(), oldCookie.Value)
	}

	token := uuid.New().String()
	expiresAt := time.Now().Add(30 * 24 * time.Hour) // 30 days

	_, err = h.store.CreateSession(r.Context(), store.CreateSessionParams{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   false, // TODO: Set to true in production with HTTPS
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(UserResponse{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		CurrentStatus: user.CurrentStatus,
		CreatedAt:     user.CreatedAt,
	})
}

func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value(UserKey).(*store.User)
	if userCtx == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(UserResponse{
		ID:            userCtx.ID,
		Name:          userCtx.Name,
		Email:         userCtx.Email,
		CurrentStatus: userCtx.CurrentStatus,
		CreatedAt:     userCtx.CreatedAt,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "No session cookie found", http.StatusBadRequest)
		return
	} else {
		_ = h.store.DeleteSession(r.Context(), cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out"))
}
