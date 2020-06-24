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
	length int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

func HttpMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		sw := statusWriter{ResponseWriter: w}
		next.ServeHTTP(&sw, r)
		duration := time.Since(startTime).Seconds()
		DefaultCollector.HttpServerRequestsDurationSummary.With(prometheus.Labels{
			"method": r.Method,
			"path":   r.URL.Path,
			"status": strconv.Itoa(sw.status),
		}).Observe(duration)
		DefaultCollector.HttpServerRequestsDurationHistogram.With(prometheus.Labels{
			"method": r.Method,
			"path":   r.URL.Path,
			"status": strconv.Itoa(sw.status),
		}).Observe(duration)
		DefaultCollector.HttpServerRequestsRequestSizeHistogram.With(prometheus.Labels{
			"method": r.Method,
			"path":   r.URL.Path,
			"status": strconv.Itoa(sw.status),
		}).Observe(float64(r.ContentLength))
		DefaultCollector.HttpServerRequestsResponseSizeHistogram.With(prometheus.Labels{
			"method": r.Method,
			"path":   r.URL.Path,
			"status": strconv.Itoa(sw.status),
		}).Observe(float64(sw.length))
	})
}

func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		duration := time.Since(startTime).Seconds()
		DefaultCollector.HttpServerRequestsDurationSummary.With(prometheus.Labels{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"status": strconv.Itoa(c.Writer.Status()),
		}).Observe(duration)
		DefaultCollector.HttpServerRequestsDurationHistogram.With(prometheus.Labels{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"status": strconv.Itoa(c.Writer.Status()),
		}).Observe(duration)
		DefaultCollector.HttpServerRequestsRequestSizeHistogram.With(prometheus.Labels{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"status": strconv.Itoa(c.Writer.Status()),
		}).Observe(float64(c.Request.ContentLength))
		DefaultCollector.HttpServerRequestsResponseSizeHistogram.With(prometheus.Labels{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"status": strconv.Itoa(c.Writer.Status()),
		}).Observe(float64(c.Writer.Size()))
	}
}

func Handler() http.Handler {
	return promhttp.Handler()
}
