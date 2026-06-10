// Package service holds all business logic (Architecture Guidelines §2.2). HTTP handlers
// and the future AI module are both just callers of these functions — the rules live here,
// once, so every caller gets them.
package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/myothiha97/ledgerflow/backend/internal/domain"
)

// errNotImplemented marks the learning-half stubs so the build stays green.
var errNotImplemented = errors.New("not implemented: TODO(you)")

// AuthStore is the consumer-defined interface (Architecture Guidelines §2.4 / §3.3):
// the auth service declares exactly the storage methods it needs. The concrete
// *store.Store satisfies it structurally, which is what lets these services be
// unit-tested against a mock with no database (see auth_test.go).
type AuthStore interface {
	CreateUser(ctx context.Context, name, email, passwordHash string) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	CreateSession(ctx context.Context, token string, userID uuid.UUID, expiresAt time.Time) (domain.Session, error)
	GetSession(ctx context.Context, token string) (domain.Session, error)
	DeleteSession(ctx context.Context, token string) error
}

// AuthService implements the auth business operations.
type AuthService struct {
	store      AuthStore
	sessionTTL time.Duration
}

// NewAuthService returns an AuthService with a default 7-day session lifetime.
func NewAuthService(store AuthStore) *AuthService {
	return &AuthService{
		store:      store,
		sessionTTL: 7 * 24 * time.Hour,
	}
}

// Register creates a user with a bcrypt-hashed password.
//
// WORKED REFERENCE: this is the full service-layer template you mirror for every
// create-style operation — normalize input, do the business work (hash), call the store,
// wrap errors with context. Note the password plaintext never leaves this function.
func (s *AuthService) Register(ctx context.Context, name, email, password string) (domain.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.store.CreateUser(ctx, name, email, string(hash))
	if err != nil {
		return domain.User{}, fmt.Errorf("register: %w", err)
	}
	return user, nil
}

// --- YOUR work (the learning half) -----------------------------------------------------

// Login verifies credentials and creates a session, returning the session token.
//
// TODO(you):
//  1. normalize email; load the user via s.store.GetUserByEmail.
//  2. compare with bcrypt.CompareHashAndPassword(user.PasswordHash, password).
//     Return domain.ErrInvalidCredentials on BOTH a wrong password AND a not-found user,
//     so the API never reveals which emails exist (Architecture Guidelines §6).
//  3. generate a random token (crypto/rand → base64/hex).
//  4. persist it: s.store.CreateSession(ctx, token, user.ID, time.Now().Add(s.sessionTTL)).
//  5. return the token and user (the handler sets the HttpOnly cookie).
func (s *AuthService) Login(ctx context.Context, email, password string) (token string, user domain.User, err error) {
	return "", domain.User{}, errNotImplemented
}

// Logout revokes the session for the given token.
//
// TODO(you): call s.store.DeleteSession. Idempotent logout (a missing session is a no-op,
// not an error) is usually the friendlier behavior.
func (s *AuthService) Logout(ctx context.Context, token string) error {
	return errNotImplemented
}

// ValidateSession resolves a session token to its user, rejecting expired sessions.
// The auth middleware calls this on every protected request.
//
// TODO(you): load via s.store.GetSession; if missing or session.ExpiresAt is in the past,
// return domain.ErrNotFound (treated as unauthenticated). Otherwise return the user via
// s.store.GetUserByID(session.UserID).
func (s *AuthService) ValidateSession(ctx context.Context, token string) (domain.User, error) {
	return domain.User{}, errNotImplemented
}
