/*
CreateTime: 2020/6/5 15:33
Author:     ALittleCat
*/

package pushgateway

import (
	"errors"
	monitor_go "github.com/alittlecat-x/monitor-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"reflect"
)

func New(url string, job string) *push.Pusher {
	p := push.New(url, job).Collector(prometheus.NewGoCollector())
	v := reflect.ValueOf(*monitor_go.DefaultCollector)
	for i := 0; i < v.NumField(); i++ {
		c, ok := v.Field(i).Interface().(prometheus.Collector)
		if !ok {
			panic(errors.New("collector must be type prometheus.Collector"))
		}
		p = p.Collector(c)
	}
	return p
}
