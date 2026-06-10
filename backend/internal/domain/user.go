// Package domain holds the core entities and rules. It is pure Go — no Gin, no SQL,
// no framework imports (Architecture Guidelines §2.1: dependencies point inward).
package domain

import (
	"time"

	"github.com/google/uuid"
)

// User is a registered account holder. PasswordHash is the bcrypt hash; the plaintext
// password never lives on this struct and is never logged.
type User struct {
	ID           uuid.UUID
	Name         string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}
