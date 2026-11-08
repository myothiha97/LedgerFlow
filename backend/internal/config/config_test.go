package config

import (
	"strings"
	"testing"
)

func TestLoad_Valid(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want Config
	}{
		{
			name: "defaults when only DATABASE_URL is set",
			env:  map[string]string{"DATABASE_URL": "postgres://localhost/ledgerflow"},
			want: Config{Port: 4000, Env: EnvProd, DatabaseURL: "postgres://localhost/ledgerflow"},
		},
		{
			name: "every value set",
			env: map[string]string{
				"PORT":         "8080",
				"APP_ENV":      "dev",
				"DATABASE_URL": "postgres://db/ledgerflow",
			},
			want: Config{Port: 8080, Env: EnvDev, DatabaseURL: "postgres://db/ledgerflow"},
		},
		{
			name: "lowest port",
			env:  map[string]string{"PORT": "1", "DATABASE_URL": "postgres://x"},
			want: Config{Port: 1, Env: EnvProd, DatabaseURL: "postgres://x"},
		},
		{
			name: "highest port",
			env:  map[string]string{"PORT": "65535", "DATABASE_URL": "postgres://x"},
			want: Config{Port: 65535, Env: EnvProd, DatabaseURL: "postgres://x"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setEnv(t, tt.env)

			got, err := Load()

			if err != nil {
				t.Fatalf("Load() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("Load() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestLoad_Invalid(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		wantErr []string
	}{
		{
			name:    "missing DATABASE_URL",
			env:     map[string]string{},
			wantErr: []string{"DATABASE_URL is required"},
		},
		{
			name:    "port is not a number",
			env:     map[string]string{"PORT": "abc", "DATABASE_URL": "postgres://x"},
			wantErr: []string{"PORT must be"},
		},
		{
			name:    "port one above the highest",
			env:     map[string]string{"PORT": "65536", "DATABASE_URL": "postgres://x"},
			wantErr: []string{"PORT must be"},
		},
		{
			name:    "negative port",
			env:     map[string]string{"PORT": "-1", "DATABASE_URL": "postgres://x"},
			wantErr: []string{"PORT must be"},
		},
		{
			name:    "unknown APP_ENV",
			env:     map[string]string{"APP_ENV": "staging", "DATABASE_URL": "postgres://x"},
			wantErr: []string{"APP_ENV must be"},
		},
		{
			name:    "reports every problem at once",
			env:     map[string]string{"PORT": "0", "APP_ENV": "staging"},
			wantErr: []string{"PORT must be", "APP_ENV must be", "DATABASE_URL is required"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setEnv(t, tt.env)

			_, err := Load()

			if err == nil {
				t.Fatal("Load() error = nil, want an error")
			}
			for _, want := range tt.wantErr {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("Load() error = %q, want it to contain %q", err, want)
				}
			}
		})
	}
}

func TestLoad_ErrorNeverShowsDatabaseURL(t *testing.T) {
	//nolint:gosec // A fake password, so the test can check it never appears in the error.
	setEnv(t, map[string]string{"PORT": "abc", "DATABASE_URL": "postgres://app:s3cret-pw@db/ledgerflow"})

	cfg, err := Load()

	if err == nil {
		t.Fatal("Load() error = nil, want an error")
	}
	if strings.Contains(err.Error(), "s3cret-pw") {
		t.Errorf("Load() error leaks the database password: %q", err)
	}
	if cfg != (Config{}) {
		t.Errorf("Load() config = %+v, want the zero Config on error", cfg)
	}
}

// setEnv clears every variable Load reads, then sets the given ones.
// t.Setenv restores the old values when the test ends.
func setEnv(t *testing.T, env map[string]string) {
	t.Helper()
	for _, key := range []string{"PORT", "APP_ENV", "DATABASE_URL"} {
		t.Setenv(key, env[key])
	}
}
