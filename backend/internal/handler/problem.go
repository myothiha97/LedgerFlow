// Package handler turns HTTP requests into service calls and service results into
// HTTP responses. It holds no business rules.
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Problem is an RFC 9457 Problem Details body. Every API error uses this shape.
// Code is a stable value the frontend can switch on; Title and Detail are for people.
type Problem struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail,omitempty"`
	Code   string `json:"code"`
}

// Stable error codes shared by every handler.
const (
	CodeNotFound       = "not_found"
	CodeInternal       = "internal_error"
	CodeUnavailable    = "unavailable"
	CodeNotImplemented = "not_implemented"
)

// writeProblem sends a Problem Details response and stops the handler chain.
// Never pass internal error text as detail: log it and send a generic message.
func writeProblem(c *gin.Context, status int, code, detail string) {
	c.Header("Content-Type", "application/problem+json")
	c.AbortWithStatusJSON(status, Problem{
		// "about:blank" is the RFC 9457 default when no extra docs page exists.
		Type:   "about:blank",
		Title:  http.StatusText(status),
		Status: status,
		Detail: detail,
		Code:   code,
	})
}
