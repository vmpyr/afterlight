package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
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
	r.Post("/register", h.Register)
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

	resp := UserResponse{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		CurrentStatus: user.CurrentStatus,
		CreatedAt:     user.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
