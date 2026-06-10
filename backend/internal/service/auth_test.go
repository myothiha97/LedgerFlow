package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/myothiha97/ledgerflow/backend/internal/domain"
)

// mockStore is a hand-written test double for AuthStore. Each method delegates to a
// function field, so a test sets only the funcs it exercises; calling an unset method
// panics, which surfaces an unexpected store call. This is the §2.4 payoff: services are
// tested with no database.
type mockStore struct {
	createUserFn     func(ctx context.Context, name, email, passwordHash string) (domain.User, error)
	getUserByEmailFn func(ctx context.Context, email string) (domain.User, error)
	getUserByIDFn    func(ctx context.Context, id uuid.UUID) (domain.User, error)
	createSessionFn  func(ctx context.Context, token string, userID uuid.UUID, expiresAt time.Time) (domain.Session, error)
	getSessionFn     func(ctx context.Context, token string) (domain.Session, error)
	deleteSessionFn  func(ctx context.Context, token string) error
}

func (m *mockStore) CreateUser(ctx context.Context, name, email, passwordHash string) (domain.User, error) {
	return m.createUserFn(ctx, name, email, passwordHash)
}
func (m *mockStore) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	return m.getUserByEmailFn(ctx, email)
}
func (m *mockStore) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return m.getUserByIDFn(ctx, id)
}
func (m *mockStore) CreateSession(ctx context.Context, token string, userID uuid.UUID, expiresAt time.Time) (domain.Session, error) {
	return m.createSessionFn(ctx, token, userID, expiresAt)
}
func (m *mockStore) GetSession(ctx context.Context, token string) (domain.Session, error) {
	return m.getSessionFn(ctx, token)
}
func (m *mockStore) DeleteSession(ctx context.Context, token string) error {
	return m.deleteSessionFn(ctx, token)
}

// TestAuthService_Register is the worked table-driven test you mirror (Arrange-Act-Assert,
// Architecture Guidelines §9). It checks the two things Register is responsible for:
// the email is normalized, and the stored value is a bcrypt hash that verifies.
func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name          string
		inputName     string
		inputEmail    string
		inputPassword string
		createUserErr error
		wantErr       bool
	}{
		{
			name:          "hashes password and normalizes email",
			inputName:     "Ada",
			inputEmail:    "  ADA@example.com ",
			inputPassword: "supersecret",
			wantErr:       false,
		},
		{
			name:          "propagates a store conflict",
			inputName:     "Ada",
			inputEmail:    "ada@example.com",
			inputPassword: "supersecret",
			createUserErr: domain.ErrEmailTaken,
			wantErr:       true,
		},
		// TODO(you): add cases — e.g. a password longer than bcrypt's 72-byte limit.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			var gotEmail, gotHash string
			store := &mockStore{
				createUserFn: func(_ context.Context, name, email, passwordHash string) (domain.User, error) {
					gotEmail, gotHash = email, passwordHash
					if tt.createUserErr != nil {
						return domain.User{}, tt.createUserErr
					}
					return domain.User{ID: uuid.New(), Name: name, Email: email, PasswordHash: passwordHash}, nil
				},
			}
			svc := NewAuthService(store)

			// Act
			user, err := svc.Register(context.Background(), tt.inputName, tt.inputEmail, tt.inputPassword)

			// Assert
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotEmail != "ada@example.com" {
				t.Errorf("email not normalized: got %q, want %q", gotEmail, "ada@example.com")
			}
			if cmpErr := bcrypt.CompareHashAndPassword([]byte(gotHash), []byte(tt.inputPassword)); cmpErr != nil {
				t.Errorf("stored hash does not verify against the password: %v", cmpErr)
			}
			if user.Email != "ada@example.com" {
				t.Errorf("returned user email = %q, want %q", user.Email, "ada@example.com")
			}
		})
	}
}

// TODO(you): implement Login, then test: valid creds → non-empty token + CreateSession
// called; wrong password → domain.ErrInvalidCredentials; unknown email → also
// domain.ErrInvalidCredentials (no user-enumeration leak).
func TestAuthService_Login(t *testing.T) {
	t.Skip("TODO(you): implement AuthService.Login, then write this table")
}

// TODO(you): implement ValidateSession, then test: valid unexpired token → user;
// expired session → domain.ErrNotFound; missing token → domain.ErrNotFound.
func TestAuthService_ValidateSession(t *testing.T) {
	t.Skip("TODO(you): implement AuthService.ValidateSession, then write this table")
}
