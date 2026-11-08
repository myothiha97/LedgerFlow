package handler

import (
	"crypto/rand"
	"errors"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	requestIDHeader = "X-Request-ID"
	requestIDKey    = "request_id"
)

// RequestID gives every request a random ID, stores it on the context, and returns
// it in the X-Request-ID header so a user report can be matched to a log line.
// It always makes a new ID: trusting a client value would let callers forge log entries.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := rand.Text()
		c.Set(requestIDKey, id)
		c.Header(requestIDHeader, id)
		c.Next()
	}
}

// Logger writes one structured log line per request after it finishes.
// It logs the path without the query string, so query values never reach the logs.
func Logger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status := c.Writer.Status()
		level := slog.LevelInfo
		if status >= http.StatusInternalServerError {
			level = slog.LevelError
		}
		logger.LogAttrs(c.Request.Context(), level, "request",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", status),
			slog.Duration("duration", time.Since(start)),
			slog.String(requestIDKey, c.GetString(requestIDKey)),
		)
	}
}

// Recovery turns a panic in a handler into a logged 500 Problem response,
// so one bad request cannot crash the server or leak a stack trace.
func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok && errors.Is(err, http.ErrAbortHandler) {
					// net/http uses this panic to abort a response on purpose.
					panic(r)
				}
				logger.ErrorContext(c.Request.Context(), "panic recovered",
					slog.Any("panic", r),
					slog.String("stack", string(debug.Stack())),
					slog.String(requestIDKey, c.GetString(requestIDKey)),
				)
				if c.Writer.Written() {
					c.Abort()
					return
				}
				writeProblem(c, http.StatusInternalServerError, CodeInternal, "")
			}
		}()
		c.Next()
	}
}
