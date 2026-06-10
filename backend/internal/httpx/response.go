// Package httpx holds the HTTP response helpers shared by every handler: one consistent
// error envelope so the frontend handles a single error shape (Architecture Guidelines §4.3).
package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorBody is the inner payload of an error response.
type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Error ErrorBody `json:"error"`
}

// Error writes the consistent error envelope: { "error": { "code", "message" } }.
func Error(c *gin.Context, status int, code, message string) {
	c.JSON(status, errorResponse{Error: ErrorBody{Code: code, Message: message}})
}

// JSON writes a success payload with the given status code.
func JSON(c *gin.Context, status int, payload any) {
	c.JSON(status, payload)
}

// NotImplemented is the placeholder response for the learning-half routes you haven't
// built yet — they return 501 through the same envelope until you fill them in.
func NotImplemented(c *gin.Context) {
	Error(c, http.StatusNotImplemented, "not_implemented", "this endpoint is not implemented yet")
}
