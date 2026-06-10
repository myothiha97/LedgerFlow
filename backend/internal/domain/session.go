package domain

import (
	"time"

	"github.com/google/uuid"
)

// Session is a server-side login session. Token is the opaque value stored in the
// HttpOnly cookie; the session is valid until ExpiresAt.
type Session struct {
	Token     string
	UserID    uuid.UUID
	ExpiresAt time.Time
	CreatedAt time.Time
}
