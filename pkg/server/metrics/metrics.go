package metrics

import (
	"net"
	"net/http"
	"strings"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"k8s.io/component-base/metrics/legacyregistry"
	"k8s.io/klog/v2"
)

var (
	// metricsList is a list of all raw metrics that should be registered always, regardless of any feature gate's value.
	metricsList       []prometheus.Collector
	grpcServerMetrics *grpcprom.ServerMetrics
)

// Register registers a list of metrics.
func Register() {
	grpcServerMetrics = grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)

	metricsList = []prometheus.Collector{
		grpcServerMetrics,
	}

	for _, metric := range metricsList {
		legacyregistry.RawMustRegister(metric)
	}
}

func ExportMetrics(metricsAddress string) {
	if metricsAddress == "" {
		return
	}
	l, err := net.Listen("tcp", metricsAddress)
	if err != nil {
		klog.Warningf("failed to get listener for metrics endpoint: %v", err)
		return
	}
	serve(l, serveMetrics)
}

func serve(l net.Listener, serveFunc func(net.Listener) error) {
	path := l.Addr().String()
	klog.V(2).Infof("set up prometheus server on %v", path)
	go func() {
		defer l.Close()
		if err := serveFunc(l); err != nil {
			klog.Fatalf("serve failure(%v), address(%v)", err, path)
		}
	}()
}

func serveMetrics(l net.Listener) error {
	m := http.NewServeMux()
	m.Handle("/metrics", legacyregistry.Handler())
	return trapClosedConnErr(http.Serve(l, m))
}

func trapClosedConnErr(err error) error {
	if err == nil {
		return nil
	}
	if strings.Contains(err.Error(), "use of closed network connection") {
		return nil
	}
	return err
}

// GRPCServerMetricsOptions returns the ServerOption applying on gRPC server
// to collect server metrics
func GRPCServerMetricsOptions() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			grpcServerMetrics.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			grpcServerMetrics.StreamServerInterceptor(),
		),
	}
}
