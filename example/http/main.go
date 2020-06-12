/*
CreateTime: 2020/6/4 11:57
Author:     ALittleCat
*/

package main

import (
	monitor_go "github.com/alittlecat-x/monitor-go"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
)

type myMonitor struct {
	MyCounter prometheus.Counter
}

func newMonitor() (*myMonitor, error) {
	myCounter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "my_counter",
			Help: "This is my counter",
		},
	)
	err := monitor_go.Register(myCounter)
	if err != nil {
		return nil, err
	}
	return &myMonitor{
		MyCounter: myCounter,
	}, nil
}

func main() {
	myMonitor, err := newMonitor()
	if err != nil {
		panic(err)
	}
	http.Handle("/hello", monitor_go.HttpMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		myMonitor.MyCounter.Inc()
	})))
	http.Handle("/metrics", monitor_go.HttpMiddleware(monitor_go.Handler()))
	server := &http.Server{
		Addr: "127.0.0.1:10086",
	}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
