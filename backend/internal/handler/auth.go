package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/myothiha97/ledgerflow/backend/internal/domain"
	"github.com/myothiha97/ledgerflow/backend/internal/httpx"
	"github.com/myothiha97/ledgerflow/backend/internal/service"
)

// sessionCookieName is the cookie that carries the session token.
const sessionCookieName = "session"

// AuthHandler adapts HTTP requests to the auth service. cookieSecure comes from config
// (off for local http, on behind HTTPS) — you'll use it when you set the login cookie.
type AuthHandler struct {
	auth         *service.AuthService
	cookieSecure bool
}

// NewAuthHandler wires the auth service and cookie policy.
func NewAuthHandler(auth *service.AuthService, cookieSecure bool) *AuthHandler {
	return &AuthHandler{auth: auth, cookieSecure: cookieSecure}
}

// registerRequest is the request *shape*. Gin's binding tags validate shape only
// (required, format); business validity belongs in the service (Architecture Guidelines §4.4).
type registerRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// userResponse is the public view of a user (never includes the password hash).
type userResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func toUserResponse(u domain.User) userResponse {
	return userResponse{ID: u.ID.String(), Name: u.Name, Email: u.Email}
}

// Register is WORKED REFERENCE #2 — the full create-style vertical slice you mirror:
// bind+validate shape → call the service → map result/error to the HTTP envelope.
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	user, err := h.auth.Register(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrEmailTaken) {
			httpx.Error(c, http.StatusConflict, "email_taken", "that email is already registered")
			return
		}
		httpx.Error(c, http.StatusInternalServerError, "internal_error", "could not register user")
		return
	}

	httpx.JSON(c, http.StatusCreated, toUserResponse(user))
}

// --- YOUR work (the learning half) -----------------------------------------------------

// Login — TODO(you): bind a {email, password} request (mirror registerRequest), call
// h.auth.Login, and on success set the HttpOnly session cookie, then respond 200 with the
// user. Use c.SetCookie(sessionCookieName, token, maxAge, "/", "", h.cookieSecure, true)
// — the final `true` is HttpOnly. Map domain.ErrInvalidCredentials → 401.
func (h *AuthHandler) Login(c *gin.Context) {
	httpx.NotImplemented(c)
}

// Logout — TODO(you): read the cookie via c.Cookie(sessionCookieName), call h.auth.Logout,
// clear the cookie (same SetCookie call with maxAge < 0), respond 204 (no body).
func (h *AuthHandler) Logout(c *gin.Context) {
	httpx.NotImplemented(c)
}

// Me — TODO(you): RequireAuth will have stored the authenticated user under contextUserKey;
// read it with c.Get(contextUserKey), type-assert to domain.User, respond 200 with
// toUserResponse(user). (RequireAuth returns 401 before this runs when there's no session.)
func (h *AuthHandler) Me(c *gin.Context) {
	httpx.NotImplemented(c)
}
