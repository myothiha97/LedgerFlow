package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

// fakePinger stands in for the database, which is a system boundary.
// Like a real database call, it must be bounded: a ping with no deadline fails,
// so a /readyz that could hang forever on a stuck database fails the tests.
type fakePinger struct{ err error }

func (f fakePinger) Ping(ctx context.Context) error {
	if _, ok := ctx.Deadline(); !ok {
		return errors.New("ping called without a deadline")
	}
	return f.err
}

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func serve(t *testing.T, h http.Handler, method, target string) *httptest.ResponseRecorder {
	t.Helper()
	return serveRequest(h, httptest.NewRequestWithContext(t.Context(), method, target, nil))
}

func serveRequest(h http.Handler, req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func decodeProblem(t *testing.T, rec *httptest.ResponseRecorder) Problem {
	t.Helper()
	if got := rec.Header().Get("Content-Type"); got != "application/problem+json" {
		t.Fatalf("Content-Type = %q, want application/problem+json", got)
	}
	var p Problem
	if err := json.Unmarshal(rec.Body.Bytes(), &p); err != nil {
		t.Fatalf("decode problem body %q: %v", rec.Body.String(), err)
	}
	return p
}

func TestHealth_ReturnsOK(t *testing.T) {
	router := NewRouter(discardLogger(), fakePinger{err: errors.New("database down")})

	rec := serve(t, router, http.MethodGet, "/healthz")

	// Liveness must stay 200 even when the database is down.
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got, want := rec.Body.String(), `{"status":"ok"}`; got != want {
		t.Errorf("body = %s, want %s", got, want)
	}
}

func TestReady_DependsOnDatabase(t *testing.T) {
	t.Skip("TODO(you): Day 1 exercise. Remove this line, run `make test`, then implement Ready in health.go")

	tests := []struct {
		name       string
		pingErr    error
		wantStatus int
		wantBody   string // checked when the database is up
		wantCode   string // checked when the database is down
	}{
		{name: "database up", pingErr: nil, wantStatus: http.StatusOK, wantBody: `{"status":"ready"}`},
		{name: "database down", pingErr: errors.New("connection refused"), wantStatus: http.StatusServiceUnavailable, wantCode: CodeUnavailable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := NewRouter(discardLogger(), fakePinger{err: tt.pingErr})

			rec := serve(t, router, http.MethodGet, "/readyz")

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantBody != "" && rec.Body.String() != tt.wantBody {
				t.Errorf("body = %s, want %s", rec.Body.String(), tt.wantBody)
			}
			if tt.wantCode != "" {
				p := decodeProblem(t, rec)
				if p.Code != tt.wantCode || p.Status != tt.wantStatus {
					t.Errorf("problem = %+v, want code %q and status %d", p, tt.wantCode, tt.wantStatus)
				}
				// The ping error is internal and must not reach the client.
				if bytes.Contains(rec.Body.Bytes(), []byte("connection refused")) {
					t.Errorf("body leaks the internal error: %s", rec.Body.String())
				}
			}
		})
	}
}

func TestRouter_UnknownPathReturnsProblem404(t *testing.T) {
	router := NewRouter(discardLogger(), fakePinger{})

	rec := serve(t, router, http.MethodGet, "/no-such-path")

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
	want := Problem{Type: "about:blank", Title: "Not Found", Status: 404, Detail: "no route matches this path", Code: CodeNotFound}
	if got := decodeProblem(t, rec); got != want {
		t.Errorf("problem = %+v, want %+v", got, want)
	}
}

func TestRequestID_UniquePerRequest(t *testing.T) {
	router := NewRouter(discardLogger(), fakePinger{})

	first := serve(t, router, http.MethodGet, "/healthz").Header().Get(requestIDHeader)
	second := serve(t, router, http.MethodGet, "/healthz").Header().Get(requestIDHeader)

	if first == "" || second == "" {
		t.Fatalf("missing %s header: first %q, second %q", requestIDHeader, first, second)
	}
	if first == second {
		t.Errorf("two requests got the same ID %q", first)
	}
}

func TestRequestID_IgnoresClientValue(t *testing.T) {
	router := NewRouter(discardLogger(), fakePinger{})
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/healthz", nil)
	req.Header.Set(requestIDHeader, "forged")

	rec := serveRequest(router, req)

	if got := rec.Header().Get(requestIDHeader); got == "forged" || got == "" {
		t.Errorf("%s = %q, want a new server-made ID", requestIDHeader, got)
	}
}

func TestLogger_OneLinePerRequestWithoutQuery(t *testing.T) {
	var buf bytes.Buffer
	router := NewRouter(slog.New(slog.NewJSONHandler(&buf, nil)), fakePinger{})

	rec := serve(t, router, http.MethodGet, "/healthz?token=secret")

	lines := bytes.Split(bytes.TrimSpace(buf.Bytes()), []byte("\n"))
	if len(lines) != 1 {
		t.Fatalf("got %d log lines, want 1:\n%s", len(lines), buf.String())
	}
	var entry struct {
		Msg       string `json:"msg"`
		Method    string `json:"method"`
		Path      string `json:"path"`
		Status    int    `json:"status"`
		RequestID string `json:"request_id"`
	}
	if err := json.Unmarshal(lines[0], &entry); err != nil {
		t.Fatalf("decode log line %s: %v", lines[0], err)
	}
	if entry.Msg != "request" || entry.Method != "GET" || entry.Path != "/healthz" || entry.Status != 200 {
		t.Errorf("log entry = %+v, want GET /healthz 200", entry)
	}
	if entry.RequestID != rec.Header().Get(requestIDHeader) {
		t.Errorf("log request_id = %q, want the response header %q", entry.RequestID, rec.Header().Get(requestIDHeader))
	}
	if bytes.Contains(buf.Bytes(), []byte("secret")) {
		t.Errorf("log leaks the query string: %s", buf.String())
	}
}

func TestRecovery_PanicReturnsProblem500(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	r := gin.New()
	r.Use(RequestID(), Recovery(logger))
	r.GET("/boom", func(*gin.Context) { panic("secret internal detail") })

	rec := serve(t, r, http.MethodGet, "/boom")

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	if p := decodeProblem(t, rec); p.Code != CodeInternal {
		t.Errorf("problem code = %q, want %q", p.Code, CodeInternal)
	}
	if bytes.Contains(rec.Body.Bytes(), []byte("secret internal detail")) {
		t.Errorf("response leaks the panic value: %s", rec.Body.String())
	}
	if !bytes.Contains(buf.Bytes(), []byte("secret internal detail")) {
		t.Errorf("panic value not logged: %s", buf.String())
	}
}

func TestLogger_ServerErrorLogsAtErrorLevel(t *testing.T) {
	var buf bytes.Buffer
	r := gin.New()
	r.Use(RequestID(), Logger(slog.New(slog.NewJSONHandler(&buf, nil))), Recovery(discardLogger()))
	r.GET("/boom", func(*gin.Context) { panic("boom") })

	serve(t, r, http.MethodGet, "/boom")

	var entry struct {
		Level  string `json:"level"`
		Status int    `json:"status"`
	}
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &entry); err != nil {
		t.Fatalf("decode log line %s: %v", buf.String(), err)
	}
	if entry.Level != "ERROR" || entry.Status != http.StatusInternalServerError {
		t.Errorf("log entry = %+v, want level ERROR and status 500", entry)
	}
}

func TestRecovery_RepanicsOnAbortHandler(t *testing.T) {
	r := gin.New()
	r.Use(Recovery(discardLogger()))
	r.GET("/abort", func(*gin.Context) { panic(http.ErrAbortHandler) })

	defer func() {
		// net/http needs this panic to reach it, so it can drop the connection quietly.
		got := recover()
		if err, _ := got.(error); !errors.Is(err, http.ErrAbortHandler) {
			t.Errorf("recovered %v, want http.ErrAbortHandler to pass through", got)
		}
	}()
	serve(t, r, http.MethodGet, "/abort")
	t.Error("ServeHTTP returned normally, want it to panic")
}
