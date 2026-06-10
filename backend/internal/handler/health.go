// Package handler holds the thin Gin HTTP layer: parse/validate request shape, call a
// service, map the result or error to an HTTP response. No business logic lives here
// (Architecture Guidelines §2.2).
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// HealthHandler reports service liveness, including database connectivity.
type HealthHandler struct {
	pool *pgxpool.Pool
}

// NewHealthHandler wires the DB pool used for the readiness ping.
func NewHealthHandler(pool *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{pool: pool}
}

// Check is WORKED REFERENCE #1 — it proves the wiring (config → pool → router) is correct
// by pinging Postgres. 200 when the DB answers, 503 when it doesn't.
func (h *HealthHandler) Check(c *gin.Context) {
	if err := h.pool.Ping(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable", "db": "down"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "db": "up"})
}
