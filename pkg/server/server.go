package server

import (
	"fmt"
	"net"
	"sync"

	"github.com/Microsoft/go-winio"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/kubernetes-csi/csi-proxy/client"
	srvtypes "github.com/kubernetes-csi/csi-proxy/pkg/server/types"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

// Server aggregates a number of API groups and versions,
// and serves requests for all of them.
type Server struct {
	versionedAPIs      []*srvtypes.VersionedAPI
	started            bool
	mutex              *sync.Mutex
	grpcServers        []*grpc.Server
	prometheusRegistry *prometheus.Registry
	prometheusMetrics  *grpcprom.ServerMetrics
}

// NewServer creates a new Server for the given API groups.
func NewServer(reg *prometheus.Registry, apiGroups ...srvtypes.APIGroup) *Server {
	versionedAPIs := make([]*srvtypes.VersionedAPI, 0, len(apiGroups))
	for _, apiGroup := range apiGroups {
		versionedAPIs = append(versionedAPIs, apiGroup.VersionedAPIs()...)
	}

	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)
	reg.MustRegister(srvMetrics)

	return &Server{
		versionedAPIs:      versionedAPIs,
		mutex:              &sync.Mutex{},
		prometheusRegistry: reg,
		prometheusMetrics:  srvMetrics,
	}
}

// Start starts one GRPC server per API version; it is a blocking call, that returns
// as soon as any of those servers shuts down (at which point it also shuts down all the
// others).
// If passed a listeningChan, it will close it when it's started listening.
func (s *Server) Start(listeningChan chan interface{}) []error {
	doneChan, ListenErr := s.startListening()
	if len(ListenErr) != 0 {
		return ListenErr
	}
	defer close(doneChan)

	if listeningChan != nil {
		close(listeningChan)
	}

	return s.waitForGRPCServersToStop(doneChan)
}

// startListening creates the named pipes, and starts GRPC servers listening on them.
func (s *Server) startListening() (chan *versionedAPIDone, []error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.started {
		return nil, []error{fmt.Errorf("server already started")}
	}
	s.started = true

	listeners, ListenErr := s.createListeners()
	if len(ListenErr) != 0 {
		return nil, ListenErr
	}

	return s.createAndStartGRPCServers(listeners), nil
}

//
//func (s *Server) createOtelExporter(ctx context.Context) {
//	exporter, err := otlptracegrpc.New(ctx,
//		otlptracegrpc.WithInsecure(),
//	)
//	if err != nil {
//		log.Fatalf("failed to create exporter: %v", err)
//	}
//
//	tp := sdktrace.NewTracerProvider(
//		sdktrace.WithSampler(sdktrace.AlwaysSample()),
//		sdktrace.WithBatcher(exporter),
//	)
//	otel.SetTracerProvider(tp)
//	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
//	//defer func() { _ = exporter.Shutdown(context.Background()) }()
//}

// createListeners creates the named pipes.
func (s *Server) createListeners() (listeners []net.Listener, errors []error) {
	listeners = make([]net.Listener, len(s.versionedAPIs))

	for i, versionedAPI := range s.versionedAPIs {
		pipePath := client.PipePath(versionedAPI.Group, versionedAPI.Version)

		listener, err := winio.ListenPipe(pipePath, nil)
		if err == nil {
			listeners[i] = listener
		} else {
			errors = append(errors, err)
		}
	}

	if len(errors) != 0 {
		// let's do a best effort to close all the listeners that we did manage to create
		for _, listener := range listeners {
			if listener != nil {
				listener.Close()
			}
		}
	}

	return
}

type versionedAPIDone struct {
	index int
	err   error
}

// createAndStartGRPCServers creates the GRPC servers, but doesn't start them just yet.
func (s *Server) createAndStartGRPCServers(listeners []net.Listener) chan *versionedAPIDone {
	doneChan := make(chan *versionedAPIDone, len(s.versionedAPIs))
	s.grpcServers = make([]*grpc.Server, len(s.versionedAPIs))

	//s.createOtelExporter(context.Background())

	for i, versionedAPI := range s.versionedAPIs {
		opts := []grpc.ServerOption{
			grpc.StatsHandler(otelgrpc.NewServerHandler()),
			grpc.ChainUnaryInterceptor(
				s.prometheusMetrics.UnaryServerInterceptor(), //grpcprom.WithExemplarFromContext(exemplarFromContext)),
				//	logging.UnaryServerInterceptor(interceptorLogger(rpcLogger), logging.WithFieldsFromContext(logTraceID)),
				//	selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(authFn), selector.MatchFunc(allButHealthZ)),
				//	recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
			),
			grpc.ChainStreamInterceptor(
				s.prometheusMetrics.StreamServerInterceptor(), //grpcprom.WithExemplarFromContext(exemplarFromContext)),
				//	logging.StreamServerInterceptor(interceptorLogger(rpcLogger), logging.WithFieldsFromContext(logTraceID)),
				//	selector.StreamServerInterceptor(auth.StreamServerInterceptor(authFn), selector.MatchFunc(allButHealthZ)),
				//	recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
			),
		}
		grpcServer := grpc.NewServer(opts...)
		s.grpcServers[i] = grpcServer

		versionedAPI.Registrant(grpcServer)
		// this next line is not a tautology, because of how go treats closures...
		index := i

		go func() {
			err := grpcServer.Serve(listeners[index])

			doneChan <- &versionedAPIDone{
				index: index,
				err:   err,
			}
		}()
	}

	return doneChan
}

func (s *Server) waitForGRPCServersToStop(doneChan chan *versionedAPIDone) (errs []error) {
	processServerDoneEvent := func(event *versionedAPIDone) {
		if event.err != nil {
			versionedAPI := s.versionedAPIs[event.index]
			err := errors.Wrapf(event.err, "GRPC server for API group %s version %s failed", versionedAPI.Group, versionedAPI.Version)
			errs = append(errs, err)
		}
	}

	// and now let's wait for at least one server to be done
	processServerDoneEvent(<-doneChan)

	// let's stop all other servers
	if err := s.Stop(); err != nil {
		// cannot happen, as the only error Stop can return is if the server hasn't been started yet
		panic(err)
	}

	// and wait for them to stop
	// TODO: do we want a timeout here?
	for doneCount := 1; doneCount < len(s.versionedAPIs); doneCount++ {
		processServerDoneEvent(<-doneChan)
	}

	return
}

// Stop stops all GRPC servers.
func (s *Server) Stop() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.started {
		return fmt.Errorf("server not started yet")
	}

	for _, grpcServer := range s.grpcServers {
		if grpcServer != nil {
			grpcServer.Stop()
		}
	}

	return nil
}
