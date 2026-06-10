package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/myothiha97/ledgerflow/backend/internal/domain"
	"github.com/myothiha97/ledgerflow/backend/internal/store/gen"
)

// uniqueViolation is the Postgres SQLSTATE for a unique-constraint conflict.
const uniqueViolation = "23505"

// CreateUser inserts a user and maps a duplicate-email conflict to domain.ErrEmailTaken.
func (s *Store) CreateUser(ctx context.Context, name, email, passwordHash string) (domain.User, error) {
	row, err := s.queries.CreateUser(ctx, gen.CreateUserParams{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return domain.User{}, domain.ErrEmailTaken
		}
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}
	return toDomainUser(row), nil
}

// GetUserByEmail looks up a user by email, mapping "no rows" to domain.ErrNotFound.
func (s *Store) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	row, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, fmt.Errorf("get user by email: %w", err)
	}
	return toDomainUser(row), nil
}

// GetUserByID looks up a user by id, mapping "no rows" to domain.ErrNotFound.
func (s *Store) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	row, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, fmt.Errorf("get user by id: %w", err)
	}
	return toDomainUser(row), nil
}

// toDomainUser maps a generated row to the framework-free domain entity.
func toDomainUser(u gen.User) domain.User {
	return domain.User{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
	}
}
