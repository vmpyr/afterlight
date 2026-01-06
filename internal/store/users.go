package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
	"github.com/vmpyr/afterlight/internal/core"
)

func (s *Store) CreateUserTx(ctx context.Context, input core.RegisterRequest) (User, error) {
	if err := core.IsValidPassword(input.Password); err != nil {
		return User{}, fmt.Errorf("password does not meet complexity requirements: %w", err)
	}
	hash, err := argon2id.CreateHash(input.Password, argon2id.DefaultParams)
	if err != nil {
		return User{}, fmt.Errorf("hashing failed: %w", err)
	}

	userID := uuid.New().String()
	now := time.Now().UTC()

	// Start transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return User{}, err
	}
	defer tx.Rollback()

	qTx := s.Queries.WithTx(tx)
	user, err := qTx.CreateUser(ctx, CreateUserParams{
		ID:                 userID,
		Name:               input.Name,
		Email:              input.Email,
		PasswordHash:       hash,
		IsPaused:           false,
		CheckInInterval:    2592000, // 30 days
		TriggerIntervalNum: 4,
		BufferPeriod:       604800, // 7 days
		VerifierQuorum:     sql.NullInt64{Int64: 1, Valid: true},
		LastCheckIn:        now,
		CurrentStatus:      core.StatusAlive,
	})
	if err != nil {
		return User{}, err
	}

	_, err = qTx.CreateContactMethod(ctx, CreateContactMethodParams{
		ID:          uuid.New().String(),
		UserID:      sql.NullString{String: userID, Valid: true},
		Channel:     "EMAIL",
		Destination: input.Email,
		Metadata:    core.Metadata{},
		CreatedAt:   now,
	})
	if err != nil {
		return User{}, err
	}

	if err := tx.Commit(); err != nil {
		return User{}, err
	}

	return user, nil
}
