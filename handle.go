/*
CreateTime: 2020/6/5 15:32
Author:     ALittleCat
*/

package monitor_go

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strconv"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func HttpMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		sw := statusWriter{w, 200}
		next.ServeHTTP(&sw, r)
		duration := time.Since(startTime).Seconds()
		DefaultCollector.HttpRequestDuration.With(prometheus.Labels{
			"method": r.Method,
			"path":   r.URL.Path,
			"status": strconv.Itoa(sw.status),
		}).Observe(duration)
	})
}

func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		duration := time.Since(startTime).Seconds()
		DefaultCollector.HttpRequestDuration.With(prometheus.Labels{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"status": strconv.Itoa(c.Writer.Status()),
		}).Observe(duration)
	}
}

func Handler() http.Handler {
	return promhttp.Handler()
}
