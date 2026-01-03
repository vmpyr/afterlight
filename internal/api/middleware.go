package api

import (
	"context"
	"net/http"
)

type ContextKey string

const UserKey ContextKey = "user"

func (h *AuthHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Unauthorized: No session cookie", http.StatusUnauthorized)
			return
		}

		user, err := h.store.GetUserBySessionToken(r.Context(), cookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserKey, &user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
