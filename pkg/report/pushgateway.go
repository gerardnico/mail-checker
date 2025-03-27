package report

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

type PushGateway struct {
	Url string `yaml:"url" validate:"omitempty,url"`
}

// PushMetrics pushes metrics to Prometheus Pushgateway
func ToPushgateway(pushgateway PushGateway, opts MetaCheck) error {

	// Create a new registry
	reg := prometheus.NewRegistry()

	// Slice to store created metrics
	var createdMetrics []prometheus.Collector

	// Process and create metrics
	for _, metricDef := range opts.Metrics {
		var metric prometheus.Collector

		switch metricDef.Type {
		case Counter:
			metric = prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: metricDef.Name,
					Help: fmt.Sprintf("Counter metric for %s", metricDef.Name),
				},
				extractLabelKeys(metricDef.Labels),
			)
			counterVec := metric.(*prometheus.CounterVec)
			counterVec.WithLabelValues(extractLabelValues(metricDef.Labels)...).Add(metricDef.Value)
		case Gauge:
			metric = prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: metricDef.Name,
					Help: fmt.Sprintf("Gauge metric for %s", metricDef.Name),
				},
				extractLabelKeys(metricDef.Labels),
			)
			gaugeVec := metric.(*prometheus.GaugeVec)
			gaugeVec.WithLabelValues(extractLabelValues(metricDef.Labels)...).Set(metricDef.Value)
		default:
			return fmt.Errorf("unsupported metric type: %s", metricDef.Type)
		}

		reg.MustRegister(metric)
		createdMetrics = append(createdMetrics, metric)
	}

	// Create pusher
	pusher := push.New(pushgateway.Url, opts.Job)
	pusher.Gatherer(reg)

	// Add job and instance labels if provided
	if opts.Job != "" {
		pusher.Grouping("job", opts.Job)
	}
	if opts.Instance != "" {
		pusher.Grouping("instance", opts.Instance)
	}

	// Set timeout
	if opts.Timeout == 0 {
		opts.Timeout = 5 * time.Second
	}
	pusher.Client(&http.Client{
		Timeout: opts.Timeout,
	})

	// Push metrics with retry
	return pushWithRetry(pusher, 3)
}

// pushWithRetry attempts to push metrics multiple times
func pushWithRetry(pusher *push.Pusher, maxRetries int) error {
	for i := 0; i < maxRetries; i++ {
		err := pusher.Push()
		if err == nil {
			return nil
		}
		log.Printf("Push attempt %d failed: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("failed to push metrics after %d attempts", maxRetries)
}

// Helper functions to extract label keys and values
func extractLabelKeys(labels map[string]string) []string {
	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	return keys
}

func extractLabelValues(labels map[string]string) []string {
	values := make([]string, 0, len(labels))
	for _, v := range labels {
		values = append(values, v)
	}
	return values
}
