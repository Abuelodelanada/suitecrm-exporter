package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadFromArgs_Defaults(t *testing.T) {
	os.Clearenv()

	cfg := LoadFromArgs([]string{})

	if cfg.ListenAddr != ":8080" {
		t.Errorf("expected default ListenAddr :8080, got %s", cfg.ListenAddr)
	}

	if cfg.HTTPTimeout != 5*time.Second {
		t.Errorf("expected default timeout 5s, got %s", cfg.HTTPTimeout)
	}
}

func TestLoadFromArgs_WithEnvVars(t *testing.T) {
	os.Setenv("SUITECRM_URL", "http://example.com")
	os.Setenv("LISTEN_ADDR", ":9999")
	os.Setenv("HTTP_TIMEOUT_SECONDS", "12")
	defer os.Clearenv()

	cfg := LoadFromArgs([]string{})

	if cfg.SuiteCRMURL != "http://example.com" {
		t.Errorf("expected SuiteCRMURL=http://example.com, got %s", cfg.SuiteCRMURL)
	}

	if cfg.ListenAddr != ":9999" {
		t.Errorf("expected ListenAddr=:9999, got %s", cfg.ListenAddr)
	}

	if cfg.HTTPTimeout != 12*time.Second {
		t.Errorf("expected timeout=12s, got %s", cfg.HTTPTimeout)
	}
}

func TestLoadFromArgs_FlagsOverrideEnv(t *testing.T) {
	os.Setenv("SUITECRM_URL", "http://env-url")
	os.Setenv("LISTEN_ADDR", ":9999")
	os.Setenv("HTTP_TIMEOUT_SECONDS", "15")
	defer os.Clearenv()

	args := []string{
		"--suitecrm-url", "http://flag-url",
		"--listen-addr", ":1234",
		"--http-timeout", "30",
	}

	cfg := LoadFromArgs(args)

	if cfg.SuiteCRMURL != "http://flag-url" {
		t.Errorf("expected SuiteCRMURL from flag, got %s", cfg.SuiteCRMURL)
	}
	if cfg.ListenAddr != ":1234" {
		t.Errorf("expected ListenAddr from flag, got %s", cfg.ListenAddr)
	}
	if cfg.HTTPTimeout != 30*time.Second {
		t.Errorf("expected HTTPTimeout=30s, got %s", cfg.HTTPTimeout)
	}
}

func TestValidate_MissingSuiteCRMURL(t *testing.T) {
	cfg := &Config{}

	err := cfg.Validate()
	if err == nil {
		t.Errorf("expected error for empty SuiteCRMURL, got nil")
	}
}

func TestValidate_WithSuiteCRMURL(t *testing.T) {
	cfg := &Config{SuiteCRMURL: "http://example.com"}

	err := cfg.Validate()
	if err != nil {
		t.Errorf("expected no error for valid SuiteCRMURL, got %v", err)
	}
}
