package main

import (
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/willglynn/purpleair_exporter/purpleair"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func handler(w http.ResponseWriter, r *http.Request, logger log.Logger) {
	query := r.URL.Query()

	target := query.Get("target")
	if len(query["target"]) != 1 || target == "" {
		http.Error(w, "'target' parameter must be specified once", http.StatusBadRequest)
		return
	}
	targetIP := net.ParseIP(target)
	if targetIP == nil {
		http.Error(w, "'target' parameter must be an IP address", http.StatusBadRequest)
		return
	}

	period := query.Get("period")
	var oneSecond, twoMinute bool
	switch period {
	case "1", "1s", "live":
		oneSecond = true
	case "2", "2m", "avg", "average":
		twoMinute = true
	default:
		oneSecond = true
		twoMinute = true
	}

	logger = log.With(logger, "target", target, "oneSecond", oneSecond, "twoMinute", twoMinute)
	level.Debug(logger).Log("msg", "Starting scrape")

	start := time.Now()
	registry := prometheus.NewRegistry()
	registry.MustRegister(purpleair.NewScraper(targetIP, oneSecond, twoMinute))
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
	duration := time.Since(start).Seconds()
	level.Debug(logger).Log("msg", "Finished scrape", "duration_seconds", duration)
}

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)

	mux := http.NewServeMux()
	mux.HandleFunc("/purpleair", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, logger)
	})

	var addr string
	addr = os.Getenv("LISTEN")
	if port := os.Getenv("PORT"); addr == "" && port != "" {
		addr = ":" + port
	}
	if addr == "" {
		addr = ":2020"
	}

	level.Info(logger).Log("msg", "Starting HTTP server", "addr", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		_ = level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}
}
