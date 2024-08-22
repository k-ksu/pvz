package metric

import (
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	OrdersGiven = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "number_of_given_orders",
		Help: "Number of given orders.",
	})
)

func RegisterMetrics() *prometheus.Registry {
	promMetrics := prometheus.NewRegistry()

	grpcMetrics := grpc_prometheus.NewServerMetrics()

	promMetrics.MustRegister(grpcMetrics, OrdersGiven)

	return promMetrics
}
