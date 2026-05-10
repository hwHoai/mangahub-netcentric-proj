package benchmarks_prometheus

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	port	string
	ActiveConnections prometheus.Gauge
	TotalRequests prometheus.Counter
	ResponsesSent prometheus.Counter
}

func InitMetrics(serverName string, port string) *Metrics {
	if port == "" {
		port = "2112"
	}
	return &Metrics{
		port: port,
		ActiveConnections: promauto.NewGauge(prometheus.GaugeOpts{
			Name: serverName+"_active_connections",
			Help: "Current number of active connections",
		}),
		TotalRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: serverName+"_requests_total",
			Help: "Total number of processed requests",
		}),
		ResponsesSent: promauto.NewCounter(prometheus.CounterOpts{
			Name: serverName+"_responses_sent_total",
			Help: "Total number of responses sent",
		}),
	}
}

func (m *Metrics) ExportMetrics() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		fmt.Println("Prometheus metrics available at http://localhost:" + m.port + "/metrics")
		http.ListenAndServe(":"+m.port, nil)
	}()
}