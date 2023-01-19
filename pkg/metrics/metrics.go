package metrics

import (
	"net/http"

	"github.com/gin-gonic/gin"
	adapter "github.com/gwatts/gin-adapter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpCnt = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "wscp_restful_http_requests_total",
			Help: "Total HTTP requests processed by the rest, excluding scrapes.",
		},
		[]string{"handler", "code", "method"},
	)
	httpPushSize = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "http_push_size_bytes",
			Help:       "HTTP request size ",
			Objectives: map[float64]float64{0.1: 0.01, 0.5: 0.05, 0.9: 0.01},
		},
		[]string{"method"},
	)
	httpPushDuration = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "wscp_restful_http_duration_seconds",
			Help:       "HTTP request duration ",
			Objectives: map[float64]float64{0.1: 0.01, 0.5: 0.05, 0.9: 0.01},
		},
		[]string{"method"},
	)
)

func InstrumentWithCounter(handlerName string, handler http.Handler) http.HandlerFunc {
	return promhttp.InstrumentHandlerCounter(
		httpCnt.MustCurryWith(prometheus.Labels{"handler": handlerName}),
		handler,
	)
}

func metrics(handler http.Handler) http.Handler {
	return promhttp.InstrumentHandlerDuration(httpPushDuration, handler)
}

func GinMetricsMiddleware() gin.HandlerFunc {

	return adapter.Wrap(metrics)
}
