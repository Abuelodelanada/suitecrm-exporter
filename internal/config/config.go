package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	SuiteCRMURL  string
	Username     string
	Password     string
	ClientID     string
	ClientSecret string
	HTTPTimeout  time.Duration
	ListenAddr   string
}

func Load() *Config {
	return LoadFromArgs(os.Args[1:])
}

func LoadFromArgs(args []string) *Config {
	defURL := os.Getenv("SUITECRM_URL")
	defUser := os.Getenv("SUITECRM_USER")
	defPass := os.Getenv("SUITECRM_PASSWORD")
	defClientID := os.Getenv("SUITECRM_CLIENT_ID")
	defClientSecret := os.Getenv("SUITECRM_CLIENT_SECRET")
	defListen := os.Getenv("LISTEN_ADDR")
	if defListen == "" {
		defListen = ":8080"
	}
	defTimeout := 10
	if s := os.Getenv("HTTP_TIMEOUT_SECONDS"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			defTimeout = v
		}
	}

	cfg := &Config{}

	fs := flag.NewFlagSet("suitecrm-exporter", flag.ContinueOnError)
	fs.StringVar(&cfg.SuiteCRMURL, "suitecrm-url", defURL, "SuiteCRM base URL")
	fs.StringVar(&cfg.Username, "suitecrm-user", defUser, "SuiteCRM username")
	fs.StringVar(&cfg.Password, "suitecrm-password", defPass, "SuiteCRM password")
	fs.StringVar(&cfg.ClientID, "client-id", defClientID, "OAuth2 client ID")
	fs.StringVar(&cfg.ClientSecret, "client-secret", defClientSecret, "OAuth2 client secret")
	fs.StringVar(&cfg.ListenAddr, "listen-addr", defListen, "Listen address")
	timeout := fs.Int("http-timeout", defTimeout, "HTTP timeout in seconds")
	_ = fs.Parse(args)

	cfg.HTTPTimeout = time.Duration(*timeout) * time.Second
	return cfg
}

func (c *Config) Validate() error {
	if c.SuiteCRMURL == "" {
		return fmt.Errorf("suitecrm URL empty")
	}
	if c.Username == "" || c.Password == "" {
		return fmt.Errorf("username or password empty")
	}
	if c.ClientID == "" || c.ClientSecret == "" {
		return fmt.Errorf("client_id or client_secret empty")
	}
	return nil
}
