package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/myothiha97/ledgerflow/backend/internal/service"
)

// NewRouter builds the Gin engine and registers the full API surface. Implemented routes
// work today; the learning-half routes are registered too (they return 501) so the whole
// surface is visible from day one.
func NewRouter(pool *pgxpool.Pool, auth *service.AuthService, cookieSecure bool) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	health := NewHealthHandler(pool)
	r.GET("/health", health.Check)

	authHandler := NewAuthHandler(auth, cookieSecure)
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/register", authHandler.Register)       // worked reference #2
		authGroup.POST("/login", authHandler.Login)             // TODO(you)
		authGroup.POST("/logout", authHandler.Logout)           // TODO(you)
		authGroup.GET("/me", RequireAuth(auth), authHandler.Me) // TODO(you)
	}

	return r
}
