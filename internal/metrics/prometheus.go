package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests processed, by status code.",
		},
		[]string{"method", "status"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of response durations for HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	pvzCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "pvz_created_total",
			Help: "Total number of PVZ created.",
		},
	)

	receptionCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "reception_created_total",
			Help: "Total number of receptions created.",
		},
	)

	productAdded = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "product_added_total",
			Help: "Total number of products added.",
		},
	)
)

func init() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(pvzCreated)
	prometheus.MustRegister(receptionCreated)
	prometheus.MustRegister(productAdded)
}

func LogRequestDuration(method string, duration time.Duration) {
	requestDuration.WithLabelValues(method).Observe(duration.Seconds())
}

func LogRequestCount(method, status string) {
	requestCount.WithLabelValues(method, status).Inc()
}

func IncrementPVZCreated() {
	pvzCreated.Inc()
}

func IncrementReceptionCreated() {
	receptionCreated.Inc()
}

func IncrementProductAdded() {
	productAdded.Inc()
}
