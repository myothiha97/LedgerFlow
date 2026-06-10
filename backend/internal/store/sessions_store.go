package store

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/myothiha97/ledgerflow/backend/internal/domain"
)

// errNotImplemented marks the learning-half stubs. Returning it (instead of leaving a
// panic) keeps the build green and lets the server boot while you fill these in.
var errNotImplemented = errors.New("not implemented: TODO(you)")

// --- Sessions store: YOUR work (the learning half) -------------------------------------
//
// TODO(you): 1) write db/queries/sessions.sql with CreateSession, GetSession,
// DeleteSession (the sessions schema is in db/migrations/0001_init.up.sql);
// 2) run `make generate`; 3) implement the three methods below against s.queries,
// mirroring users_store.go — wrap errors with fmt.Errorf("...: %w", err) and map
// pgx.ErrNoRows to domain.ErrNotFound. See Architecture Guidelines §5 (DB) and §3.2.

// CreateSession persists a new session and returns it.
func (s *Store) CreateSession(ctx context.Context, token string, userID uuid.UUID, expiresAt time.Time) (domain.Session, error) {
	return domain.Session{}, errNotImplemented
}

// GetSession loads a session by its token (return domain.ErrNotFound when absent).
func (s *Store) GetSession(ctx context.Context, token string) (domain.Session, error) {
	return domain.Session{}, errNotImplemented
}

// DeleteSession removes a session by token (used on logout).
func (s *Store) DeleteSession(ctx context.Context, token string) error {
	return errNotImplemented
}
