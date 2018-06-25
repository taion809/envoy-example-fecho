package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Subsystem: "echoer",
		Name:      "http_requests_total",
		Help:      "The total number of http requests",
	})
)

func init() {
	prometheus.MustRegister(requestCounter)
}

func main() {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		requestCounter.Inc()

		w.WriteHeader(http.StatusOK)
	})

	r.Get("/junk", func(w http.ResponseWriter, r *http.Request) {
		requestCounter.Inc()

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(junk))
	})

	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	http.ListenAndServe(":5555", r)
}
