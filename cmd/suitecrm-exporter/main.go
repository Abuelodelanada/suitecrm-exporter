package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Abuelodelanada/suitecrm-exporter/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	suitecrmUp = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "suitecrm_up",
			Help: "Shows whether the SuiteCRM exporter is working or not.",
		},
	)
)

func init() {
	prometheus.MustRegister(suitecrmUp)
}

func initialiseMetrics() {
	suitecrmUp.Set(1)
}

func main() {
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		// Lo dejamos no-fatal por ahora; más adelante podrías decidir fallar si falta URL.
		log.Printf("WARNING: %v. The exporter will start, but suitecrm_up will not reflect the API until configuration is complete.", err)
	}
	initialiseMetrics()
	http.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{
		Addr:              cfg.ListenAddr,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Initialising suitecrm-exporter in %s (timeout=%s)", cfg.ListenAddr, cfg.HTTPTimeout)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
