/*
CreateTime: 2020/6/3 14:55
Author:     ALittleCat
*/

package monitor_go

import (
	"github.com/prometheus/client_golang/prometheus"
)

func Register(collector prometheus.Collector) error {
	return prometheus.Register(collector)
}

func MustRegister(collectors ...prometheus.Collector) {
	prometheus.MustRegister(collectors...)
}
