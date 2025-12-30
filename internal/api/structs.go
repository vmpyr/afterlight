package api

import (
	"time"

	"github.com/vmpyr/afterlight/internal/core"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	Email         string          `json:"email"`
	CurrentStatus core.UserStatus `json:"current_status"`
	CreatedAt     time.Time       `json:"created_at"`
}
