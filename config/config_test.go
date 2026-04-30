package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		f, err := os.CreateTemp(t.TempDir(), "config-*.yaml")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer f.Close()

		_, err = f.WriteString(`
redis:
  host: "localhost:6379"
  password: "secret"
  db: 1
firestore:
  projectID: "my-project"
  credentialsFile: "/path/to/creds.json"
mappings:
  - source: "events_queue"
    target: "events"
  - source: "users_queue"
    target: "users"
`)
		if err != nil {
			t.Fatalf("failed to write config: %v", err)
		}

		cfg, err := Load(f.Name())
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		if cfg.Redis.Host != "localhost:6379" {
			t.Errorf("Redis.Host = %q, want %q", cfg.Redis.Host, "localhost:6379")
		}
		if cfg.Redis.Password != "secret" {
			t.Errorf("Redis.Password = %q, want %q", cfg.Redis.Password, "secret")
		}
		if cfg.Redis.DB != 1 {
			t.Errorf("Redis.DB = %d, want %d", cfg.Redis.DB, 1)
		}
		if cfg.Firestore.ProjectID != "my-project" {
			t.Errorf("Firestore.ProjectID = %q, want %q", cfg.Firestore.ProjectID, "my-project")
		}
		if cfg.Firestore.CredentialsFile != "/path/to/creds.json" {
			t.Errorf("Firestore.CredentialsFile = %q, want %q", cfg.Firestore.CredentialsFile, "/path/to/creds.json")
		}
		if len(cfg.Mappings) != 2 {
			t.Fatalf("len(Mappings) = %d, want 2", len(cfg.Mappings))
		}
		if cfg.Mappings[0].Source != "events_queue" || cfg.Mappings[0].Target != "events" {
			t.Errorf("Mappings[0] = %+v, want {Source:events_queue Target:events}", cfg.Mappings[0])
		}
		if cfg.Mappings[1].Source != "users_queue" || cfg.Mappings[1].Target != "users" {
			t.Errorf("Mappings[1] = %+v, want {Source:users_queue Target:users}", cfg.Mappings[1])
		}
	})

	t.Run("file not found", func(t *testing.T) {
		_, err := Load("/nonexistent/path/config.yaml")
		if err == nil {
			t.Error("Load() expected error for missing file, got nil")
		}
	})

	t.Run("invalid yaml", func(t *testing.T) {
		f, err := os.CreateTemp(t.TempDir(), "config-*.yaml")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer f.Close()

		if _, err = f.WriteString(":: invalid: yaml: ["); err != nil {
			t.Fatalf("failed to write config: %v", err)
		}

		_, err = Load(f.Name())
		if err == nil {
			t.Error("Load() expected error for invalid YAML, got nil")
		}
	})

	t.Run("redis password from env", func(t *testing.T) {
		f, err := os.CreateTemp(t.TempDir(), "config-*.yaml")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer f.Close()

		if _, err = f.WriteString("redis:\n  host: \"localhost:6379\"\n"); err != nil {
			t.Fatalf("failed to write config: %v", err)
		}

		t.Setenv("REDIS_PASSWORD", "env-password")

		cfg, err := Load(f.Name())
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
		if cfg.Redis.Password != "env-password" {
			t.Errorf("Redis.Password = %q, want %q", cfg.Redis.Password, "env-password")
		}
	})

	t.Run("firestore credentials from env", func(t *testing.T) {
		f, err := os.CreateTemp(t.TempDir(), "config-*.yaml")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer f.Close()

		if _, err = f.WriteString("firestore:\n  projectID: \"my-project\"\n"); err != nil {
			t.Fatalf("failed to write config: %v", err)
		}

		t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/env/creds.json")

		cfg, err := Load(f.Name())
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
		if cfg.Firestore.CredentialsFile != "/env/creds.json" {
			t.Errorf("Firestore.CredentialsFile = %q, want %q", cfg.Firestore.CredentialsFile, "/env/creds.json")
		}
	})

	t.Run("firestore credentials file not overridden by env when set in config", func(t *testing.T) {
		f, err := os.CreateTemp(t.TempDir(), "config-*.yaml")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer f.Close()

		if _, err = f.WriteString("firestore:\n  projectID: \"my-project\"\n  credentialsFile: \"/config/creds.json\"\n"); err != nil {
			t.Fatalf("failed to write config: %v", err)
		}

		t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/env/creds.json")

		cfg, err := Load(f.Name())
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
		if cfg.Firestore.CredentialsFile != "/config/creds.json" {
			t.Errorf("Firestore.CredentialsFile = %q, want %q (env should not override config)", cfg.Firestore.CredentialsFile, "/config/creds.json")
		}
	})
}
