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
	HttpRequestDuration *prometheus.HistogramVec
}

func NewDefaultCollector() *collectors {
	collectors := collectors{
		HttpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration",
				Help:    "The HTTP request latencies in seconds",
				Buckets: prometheus.DefBuckets,
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
