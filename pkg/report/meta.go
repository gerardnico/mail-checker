package report

import "time"

// MetricType represents the allowed types of metrics
type MetricType string

// Constant values for the allowed metric types
const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

// MetricDefinition represents a single metric to be pushed
type MetricDefinition struct {
	Name   string
	Type   MetricType
	Value  float64
	Labels map[string]string
}

// MetaCheck contains information about the checks for pushing metrics
type MetaCheck struct {
	Job      string
	Instance string
	Errors   []string
	Metrics  []MetricDefinition
	Timeout  time.Duration
}
