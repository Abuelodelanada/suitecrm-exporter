package collector

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Abuelodelanada/suitecrm-exporter/internal/config"
)

func TestPing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/Api/V8/me" {
			t.Errorf("expected URL /Api/V8/me, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := &config.Config{
		SuiteCRMURL: server.URL,
		HTTPTimeout: 2 * time.Second,
	}

	client := NewClient(cfg)

	up, err := client.Ping()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !up {
		t.Fatalf("expected up=true, got false")
	}
}

func TestPing_Non200Status(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	cfg := &config.Config{
		SuiteCRMURL: server.URL,
		HTTPTimeout: 2 * time.Second,
	}

	client := NewClient(cfg)
	up, err := client.Ping()

	if err == nil {
		t.Fatalf("expected error for non-200 status, got nil")
	}
	if up {
		t.Fatalf("expected up=false for non-200 status, got true")
	}
}


func TestPing_WithToken(t *testing.T) {
	expectedToken := "my-secret-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer "+expectedToken {
			t.Errorf("expected Authorization header 'Bearer %s', got '%s'", expectedToken, auth)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := &config.Config{
		SuiteCRMURL: server.URL,
		APIToken:    expectedToken,
		HTTPTimeout: 2 * time.Second,
	}

	client := NewClient(cfg)
	up, err := client.Ping()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !up {
		t.Fatalf("expected up=true, got false")
	}
}


func TestPing_WithBasicAuth(t *testing.T) {
	expectedUser := "user"
	expectedPass := "pass"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			t.Errorf("expected basic auth, none found")
		}
		if user != expectedUser || pass != expectedPass {
			t.Errorf("expected user/pass '%s/%s', got '%s/%s'", expectedUser, expectedPass, user, pass)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := &config.Config{
		SuiteCRMURL: server.URL,
		Username:    expectedUser,
		Password:    expectedPass,
		HTTPTimeout: 2 * time.Second,
	}

	client := NewClient(cfg)
	up, err := client.Ping()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !up {
		t.Fatalf("expected up=true, got false")
	}
}
