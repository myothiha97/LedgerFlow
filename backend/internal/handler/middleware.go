package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/myothiha97/ledgerflow/backend/internal/service"
)

// contextUserKey is the gin.Context key under which RequireAuth stores the authenticated
// user for downstream handlers (e.g. AuthHandler.Me) to read.
const contextUserKey = "currentUser"

// RequireAuth is the session-cookie auth middleware (Architecture Guidelines §6.3:
// authenticate every /api/* route except register/login, and scope work by user_id).
//
// SKELETON — the signature and wiring are provided; you implement the body.
//
// TODO(you):
//  1. token, err := c.Cookie(sessionCookieName); if err != nil → 401 + c.Abort().
//  2. user, err := auth.ValidateSession(c.Request.Context(), token); if err != nil → 401 + c.Abort().
//  3. c.Set(contextUserKey, user); c.Next().
//
// For the 401 paths use httpx.Error(c, http.StatusUnauthorized, "unauthorized", "...")
// then c.Abort() so the protected handler never runs.
func RequireAuth(auth *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO(you): implement per the contract above. Until then this is a pass-through
		// and the protected handlers themselves return 501.
		c.Next()
	}
}
