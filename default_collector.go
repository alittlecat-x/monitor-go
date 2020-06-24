/*
CreateTime: 2020/6/5 14:27
Author:     ALittleCat
*/

package monitor_go

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"reflect"
)

var DefaultCollector = NewDefaultCollector()

type collectors struct {
	HttpServerRequestsDurationSummary *prometheus.SummaryVec
	HttpServerRequestsDurationHistogram *prometheus.HistogramVec
	HttpServerRequestsRequestSizeHistogram *prometheus.HistogramVec
	HttpServerRequestsResponseSizeHistogram *prometheus.HistogramVec
}

func NewDefaultCollector() *collectors {
	collectors := collectors{
		HttpServerRequestsDurationSummary: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:    "http_server_requests_seconds_summary",
				Help:    "A summary of request latencies for requests.",
				Objectives: map[float64]float64{0.5:0.01, 0.75: 0.01, 0.9: 0.01, 0.95: 0.001, 0.99: 0.001},
			},
			[]string{"method", "path", "status"},
		),
		HttpServerRequestsDurationHistogram: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_server_requests_seconds_histogram",
				Help:    "A histogram of request size for requests.",
				Buckets: []float64{10, 50, 100, 200, 500, 1000, 2000, 3000, 5000, 10000},
			},
			[]string{"method", "path", "status"},
		),
		HttpServerRequestsRequestSizeHistogram: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_server_requests_request_size_bytes",
				Help:    "A histogram of request size for requests.",
				Buckets: []float64{200, 500, 900, 1500, 3000, 10000},
			},
			[]string{"method", "path", "status"},
		),
		HttpServerRequestsResponseSizeHistogram: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_server_requests_response_size_bytes",
				Help:    "A histogram of response size for requests.",
				Buckets: []float64{200, 500, 900, 1500, 3000, 10000},
			},
			[]string{"method", "path", "status"},
		),
	}
	v := reflect.ValueOf(collectors)
	collector := make([]prometheus.Collector, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		c, ok := v.Field(i).Interface().(prometheus.Collector)
		if !ok {
			panic(errors.New("collector must be type prometheus.Collector"))
		}
		collector[i] = c
	}
	prometheus.MustRegister(collector...)
	return &collectors
}
