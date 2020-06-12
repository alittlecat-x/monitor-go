/*
CreateTime: 2020/6/5 11:03
Author:     ALittleCat
*/

package main

import (
	monitor_go "github.com/alittlecat-x/monitor-go"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	myCounter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "my_counter",
			Help: "This is my counter",
		},
	)
	if err := monitor_go.Register(myCounter); err != nil {
		panic(err)
	}
	p := pushgateway.New("127.0.0.1:10086", "test").Collector(myCounter)
	if err := p.Push(); err != nil {
		panic(err)
	}
}
