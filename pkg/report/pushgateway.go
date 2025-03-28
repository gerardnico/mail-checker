// see example at https://pkg.go.dev/github.com/prometheus/client_golang
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

	// The registry will all metrics
	var reg = prometheus.NewRegistry()

	// Collectors (ie CounterVec, GaugeVec, ...) can only be registered/created once
	// Collectors are vector that gather Metrics data
	var createdVector = make(map[string]prometheus.Collector)

	// Process and add metrics to the vectors
	for _, metricDef := range opts.Metrics {

		switch metricDef.Type {
		case Counter:
			counterVec := getOrCreateVector(metricDef, createdVector, reg).(*prometheus.CounterVec)
			// Add a metric to the vector
			counterVec.
				With(metricDef.Labels).
				Add(metricDef.Value)
		case Gauge:
			gaugeVec := getOrCreateVector(metricDef, createdVector, reg).(*prometheus.GaugeVec)
			// Add a metric to the vector
			gaugeVec.
				With(metricDef.Labels).
				Set(metricDef.Value)
		default:
			return fmt.Errorf("unsupported metric type: %s", metricDef.Type)
		}

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

func getOrCreateVector(metricDef MetricDefinition, createdCollector map[string]prometheus.Collector, reg *prometheus.Registry) prometheus.Collector {
	key := metricDef.Name
	// Check if the key exists
	// Could do with reg.Collectors()
	metric, exists := createdCollector[key]
	if exists {
		return metric
	}
	var collectorVec prometheus.Collector
	switch metricDef.Type {
	case Counter:
		counterOpts := prometheus.CounterOpts{
			Name: metricDef.Name,
			Help: fmt.Sprintf("Counter metric for %s", metricDef.Name),
		}
		collectorVec = prometheus.NewCounterVec(
			counterOpts,
			extractLabelKeys(metricDef.Labels),
		)
	case Gauge:
		gaugeOpts := prometheus.GaugeOpts{
			Name: metricDef.Name,
			Help: fmt.Sprintf("Gauge metric for %s", metricDef.Name),
		}
		collectorVec = prometheus.NewGaugeVec(
			gaugeOpts,
			extractLabelKeys(metricDef.Labels),
		)
	default:
		log.Fatalf("unsupported metric type: %s", metricDef.Type)
		return nil
	}

	reg.MustRegister(collectorVec)
	createdCollector[key] = collectorVec
	return collectorVec

}
