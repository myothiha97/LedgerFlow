package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Pinger checks that a dependency answers. *pgxpool.Pool satisfies it.
type Pinger interface {
	Ping(ctx context.Context) error
}

// Health handles GET /healthz (liveness). It answers 200 while the process runs
// and never checks dependencies, so a database outage does not restart the app.
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Ready handles GET /readyz (readiness). It tells a load balancer whether this
// instance can serve traffic right now, which needs a working database.
func Ready(db Pinger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO(you): Day 1 exercise. Ping db with a short timeout (2 seconds) and
		// answer 200 {"status":"ready"} when it works. When it fails, answer with
		// writeProblem: status 503 and code CodeUnavailable. Then remove the
		// t.Skip line in TestReady_DependsOnDatabase (handler_test.go) and make it pass.
		_ = db
		writeProblem(c, http.StatusNotImplemented, CodeNotImplemented, "readiness check is not implemented yet")
	}
}
