/*
CreateTime: 2020/6/3 17:48
Author:     ALittleCat
*/

package main

import (
	monitor_go "github.com/alittlecat-x/monitor-go"
	"github.com/gin-gonic/gin"
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
	router := gin.New()
	router.Use(monitor_go.GinMiddleware())
	router.GET("/hello", func(c *gin.Context) {
		myMonitor.MyCounter.Inc()
		c.AbortWithStatus(http.StatusOK)
	})
	router.GET("/metrics", gin.WrapH(monitor_go.Handler()))
	if err := router.Run("127.0.0.1:10086"); err != nil {
		panic(err)
	}
}
