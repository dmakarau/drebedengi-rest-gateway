package config

import (
	"os"
	"testing"
)

func setEnv(t *testing.T, pairs ...string) {
	t.Helper()
	for i := 0; i < len(pairs); i += 2 {
		key, val := pairs[i], pairs[i+1]
		orig, exists := os.LookupEnv(key)
		os.Setenv(key, val)
		t.Cleanup(func() {
			if exists {
				os.Setenv(key, orig)
			} else {
				os.Unsetenv(key)
			}
		})
	}
}

func unsetEnv(t *testing.T, keys ...string) {
	t.Helper()
	for _, key := range keys {
		orig, exists := os.LookupEnv(key)
		os.Unsetenv(key)
		t.Cleanup(func() {
			if exists {
				os.Setenv(key, orig)
			}
		})
	}
}

func TestLoad_AllVarsPresent(t *testing.T) {
	setEnv(t,
		"DD_API_ID", "myapi",
		"DD_LOGIN", "user@example.com",
		"DD_PASSWORD", "secret",
		"DD_SOAP_URL", "http://custom.soap/",
		"PORT", "9090",
	)

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.APIId != "myapi" {
		t.Errorf("APIId = %q, want %q", cfg.APIId, "myapi")
	}
	if cfg.Login != "user@example.com" {
		t.Errorf("Login = %q, want %q", cfg.Login, "user@example.com")
	}
	if cfg.Password != "secret" {
		t.Errorf("Password = %q, want %q", cfg.Password, "secret")
	}
	if cfg.SoapURL != "http://custom.soap/" {
		t.Errorf("SoapURL = %q, want %q", cfg.SoapURL, "http://custom.soap/")
	}
	if cfg.Port != "9090" {
		t.Errorf("Port = %q, want %q", cfg.Port, "9090")
	}
}

func TestLoad_Defaults(t *testing.T) {
	setEnv(t,
		"DD_API_ID", "myapi",
		"DD_LOGIN", "user@example.com",
		"DD_PASSWORD", "secret",
	)
	unsetEnv(t, "DD_SOAP_URL", "PORT")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.SoapURL != "http://www.drebedengi.ru/soap/" {
		t.Errorf("SoapURL = %q, want default", cfg.SoapURL)
	}
	if cfg.Port != "8080" {
		t.Errorf("Port = %q, want default 8080", cfg.Port)
	}
}

func TestLoad_MissingAPIId(t *testing.T) {
	unsetEnv(t, "DD_API_ID")
	setEnv(t, "DD_LOGIN", "u", "DD_PASSWORD", "p")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing DD_API_ID")
	}
}

func TestLoad_MissingLogin(t *testing.T) {
	setEnv(t, "DD_API_ID", "api")
	unsetEnv(t, "DD_LOGIN")
	setEnv(t, "DD_PASSWORD", "p")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing DD_LOGIN")
	}
}

func TestLoad_MissingPassword(t *testing.T) {
	setEnv(t, "DD_API_ID", "api", "DD_LOGIN", "u")
	unsetEnv(t, "DD_PASSWORD")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing DD_PASSWORD")
	}
}
