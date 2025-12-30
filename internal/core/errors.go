package core

import (
	"errors"
)

// All custom error definitions
var ErrUserNotFound = errors.New("user not found")
var ErrPasswordLength = errors.New("password must be at least 8 characters")
var ErrWeakPassword = errors.New("password must contain at least one uppercase letter, one lowercase letter, one digit, and one special character")
