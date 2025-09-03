package main

import (
	"log"

	"github.com/Abuelodelanada/suitecrm-exporter/internal/collector"
	"github.com/Abuelodelanada/suitecrm-exporter/internal/config"
)

func main() {
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	client := collector.NewClient(
		cfg.SuiteCRMURL,
		cfg.Username,
		cfg.Password,
		cfg.ClientID,
		cfg.ClientSecret,
	)

	if err := client.Login(); err != nil {
		log.Fatalf("login failed: %v", err)
	}
	log.Println("Access token acquired ✅")

	accounts, err := client.FetchAccounts()
	if err != nil {
		log.Fatalf("fetch accounts failed: %v", err)
	}
	log.Printf("Accounts response: %s", string(accounts))
}
