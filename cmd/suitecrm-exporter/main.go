package main

import (
	"net/http"

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
	initialiseMetrics()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
