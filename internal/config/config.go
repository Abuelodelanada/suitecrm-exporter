package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	SuiteCRMURL string
	Username    string
	Password    string
	APIToken    string
	HTTPTimeout time.Duration
	ListenAddr  string
}

func Load() *Config {
	return LoadFromArgs(os.Args[1:])
}

func LoadFromArgs(args []string) *Config {
	defURL := os.Getenv("SUITECRM_URL")
	defUser := os.Getenv("SUITECRM_USER")
	defPass := os.Getenv("SUITECRM_PASSWORD")
	defToken := os.Getenv("SUITECRM_TOKEN")
	defListen := os.Getenv("LISTEN_ADDR")
	if defListen == "" {
		defListen = ":8080"
	}
	defTimeout := 5
	if s := os.Getenv("HTTP_TIMEOUT_SECONDS"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			defTimeout = v
		}
	}

	cfg := &Config{}

	// Use a local FlagSet so tests no rompen el flag global
	fs := flag.NewFlagSet("suitecrm-exporter", flag.ContinueOnError)
	fs.StringVar(&cfg.SuiteCRMURL, "suitecrm-url", defURL, "SuiteCRM base URL (or env SUITECRM_URL)")
	fs.StringVar(&cfg.Username, "suitecrm-user", defUser, "SuiteCRM username (or env SUITECRM_USER)")
	fs.StringVar(&cfg.Password, "suitecrm-password", defPass, "SuiteCRM password (or env SUITECRM_PASSWORD)")
	fs.StringVar(&cfg.APIToken, "suitecrm-token", defToken, "SuiteCRM API token (or env SUITECRM_TOKEN)")
	fs.StringVar(&cfg.ListenAddr, "listen-addr", defListen, "Address to listen on (or env LISTEN_ADDR)")
	timeout := fs.Int("http-timeout", defTimeout, "HTTP timeout in seconds (or env HTTP_TIMEOUT_SECONDS)")

	// Parse (ignoramos el error por simplicidad; en producción podrías devolverlo)
	_ = fs.Parse(args)

	cfg.HTTPTimeout = time.Duration(*timeout) * time.Second
	return cfg
}

func (c *Config) Validate() error {
	if c.SuiteCRMURL == "" {
		return fmt.Errorf("suitecrm URL empty: configure SUITECRM_URL or --suitecrm-url")
	}
	return nil
}
