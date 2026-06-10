package store

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/myothiha97/ledgerflow/backend/internal/store/gen"
)

// Store is the concrete data-access implementation. Its methods (across users_store.go
// and sessions_store.go) satisfy the consumer-defined interfaces declared in the service
// package ("accept interfaces, return structs" — Architecture Guidelines §3.3).
type Store struct {
	pool    *pgxpool.Pool
	queries *gen.Queries
}

// New builds a Store over an existing pool. gen.Queries runs its statements through the
// pool (which satisfies sqlc's DBTX interface).
func New(pool *pgxpool.Pool) *Store {
	return &Store{
		pool:    pool,
		queries: gen.New(pool),
	}
}
