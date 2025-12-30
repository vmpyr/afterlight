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

type RegisterUserParams struct {
	Name     string
	Email    string
	Password string
}

func (s *Store) CreateUserTx(ctx context.Context, input RegisterUserParams) (*User, error) {
	hash, err := argon2id.CreateHash(input.Password, argon2id.DefaultParams)
	if err != nil {
		return nil, fmt.Errorf("hashing failed: %w", err)
	}

	userID := uuid.New().String()
	now := time.Now().UTC()

	// Start transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	_, err = qTx.CreateContactMethod(ctx, CreateContactMethodParams{
		ID:        uuid.New().String(),
		UserID:    sql.NullString{String: userID, Valid: true},
		Channel:   "EMAIL",
		Target:    input.Email,
		Metadata:  core.Metadata{},
		CreatedAt: now,
	})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &User{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		CurrentStatus: user.CurrentStatus,
		CreatedAt:     user.CreatedAt,
	}, nil
}
