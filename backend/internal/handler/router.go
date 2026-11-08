package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewRouter builds the HTTP routes and the middleware every request passes through.
// Order matters: the request ID must exist before the logger and recovery read it.
func NewRouter(logger *slog.Logger, db Pinger) *gin.Engine {
	r := gin.New()
	// Trust no proxy headers, so X-Forwarded-For cannot fake c.ClientIP().
	// Add the real proxy here when the app is deployed behind one. A nil list never errors.
	_ = r.SetTrustedProxies(nil)
	r.Use(RequestID(), Logger(logger), Recovery(logger))

	r.GET("/healthz", Health)
	r.GET("/readyz", Ready(db))

	r.NoRoute(func(c *gin.Context) {
		writeProblem(c, http.StatusNotFound, CodeNotFound, "no route matches this path")
	})
	return r
}
