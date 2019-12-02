package server

import (
	"crypto/tls"
	"net"
	"net/http"
	"strings"
	"time"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	ilogger "github.com/meateam/elasticsearch-logger"
	pb "github.com/meateam/search-service/proto"
	"github.com/meateam/search-service/service"
	"github.com/meateam/search-service/service/elasticsearch"
	es "github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const (
	envPrefix                   = "SS"
	configPort                  = "port"
	configElasticsearchURL      = "elasticsearch_url"
	configElasticsearchUser     = "elasticsearch_user"
	configElasticsearchPassword = "elasticsearch_password"
	configElasticsearchIndex    = "elasticsearch_index"
	configTLSSkipVerify         = "tls_skip_verify"
	configHealthCheckInterval   = "health_check_interval"
	configElasticAPMIgnoreURLS  = "elastic_apm_ignore_urls"
	configElasticsearchSniff    = "elasticsearch_sniff"
)

func init() {
	viper.SetDefault(configPort, "8080")
	viper.SetDefault(configElasticsearchURL, "http://localhost:9200")
	viper.SetDefault(configElasticsearchUser, "")
	viper.SetDefault(configElasticsearchPassword, "")
	viper.SetDefault(configElasticsearchIndex, "files")
	viper.SetDefault(configTLSSkipVerify, true)
	viper.SetDefault(configHealthCheckInterval, 3)
	viper.SetDefault(configElasticAPMIgnoreURLS, "/grpc.health.v1.Health/Check")
	viper.SetDefault(configElasticsearchSniff, false)
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
}

// SearchServer is a structure that holds the search grpc server
// and its services and configuration.
type SearchServer struct {
	*grpc.Server
	logger              *logrus.Logger
	port                string
	healthCheckInterval int
	SearchService       service.Service
}

// Serve accepts incoming connections on the listener `lis`, creating a new
// ServerTransport and service goroutine for each. The service goroutines
// read gRPC requests and then call the registered handlers to reply to them.
// Serve returns when `lis.Accept` fails with fatal errors. `lis` will be closed when
// this method returns.
// If `lis` is nil then Serve creates a `net.Listener` with "tcp" network listening
// on the configured `TCP_PORT`, which defaults to "8080".
// Serve will return a non-nil error unless Stop or GracefulStop is called.
func (s SearchServer) Serve(lis net.Listener) {
	listener := lis
	if lis == nil {
		l, err := net.Listen("tcp", ":"+s.port)
		if err != nil {
			s.logger.Fatalf("failed to listen: %v", err)
		}

		listener = l
	}

	s.logger.Infof("listening and serving grpc server on port %s", s.port)
	if err := s.Server.Serve(listener); err != nil {
		s.logger.Fatalf(err.Error())
	}
}

// NewServer configures and creates a grpc.Server instance with the download service
// health check service.
// Configure using environment variables.
// `HEALTH_CHECK_INTERVAL`: Interval to update serving state of the health check server.
// `PORT`: TCP port on which the grpc server would serve on.
func NewServer(logger *logrus.Logger) *SearchServer {
	// If no logger is given, create a new default logger for the server.
	if logger == nil {
		logger = ilogger.NewLogger()
	}

	// Set up grpc server opts with logger interceptor.
	serverOpts := append(
		serverLoggerInterceptor(logger),
		grpc.MaxRecvMsgSize(16<<20),
	)

	// Create a new grpc server.
	grpcServer := grpc.NewServer(
		serverOpts...,
	)

	controller, err := initController()
	if err != nil {
		logger.Fatalf("%v", err)
	}

	// Create a search service and register it on the grpc server.
	searchService := service.NewService(controller, logger)
	pb.RegisterSearchServer(grpcServer, searchService)

	// Create a health server and register it on the grpc server.
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	searchServer := &SearchServer{
		Server:              grpcServer,
		logger:              logger,
		port:                viper.GetString(configPort),
		healthCheckInterval: viper.GetInt(configHealthCheckInterval),
		SearchService:       searchService,
	}

	// Health check validation goroutine worker.
	go searchServer.healthCheckWorker(healthServer)

	return searchServer
}

func initController() (service.Controller, error) {
	controller, err := elasticsearch.NewController(initESConfig())
	if err != nil {
		return nil, err
	}
	return controller, nil
}

func initESConfig() ([]es.ClientOptionFunc, string) {
	elasticURL := viper.GetString(configElasticsearchURL)
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: viper.GetBool(configTLSSkipVerify), // ignore expired SSL certificates
		},
	}
	httpClient := &http.Client{Transport: transCfg}

	elasticOpts := []es.ClientOptionFunc{
		es.SetURL(strings.Split(elasticURL, ",")...),
		es.SetSniff(viper.GetBool(configElasticsearchSniff)),
		es.SetHttpClient(httpClient),
	}

	elasticUser := viper.GetString(configElasticsearchUser)
	elasticPassword := viper.GetString(configElasticsearchPassword)
	if elasticUser != "" && elasticPassword != "" {
		elasticOpts = append(elasticOpts, es.SetBasicAuth(elasticUser, elasticPassword))
	}

	return elasticOpts, viper.GetString(configElasticsearchIndex)
}

// serverLoggerInterceptor configures the logger interceptor for the search server.
func serverLoggerInterceptor(logger *logrus.Logger) []grpc.ServerOption {
	// Create new logrus entry for logger interceptor.
	logrusEntry := logrus.NewEntry(logger)

	ignorePayload := ilogger.IgnoreServerMethodsDecider(
		strings.Split(viper.GetString(configElasticAPMIgnoreURLS), ",")...,
	)

	ignoreInitialRequest := ilogger.IgnoreServerMethodsDecider(
		strings.Split(viper.GetString(configElasticAPMIgnoreURLS), ",")...,
	)

	// Shared options for the logger, with a custom gRPC code to log level function.
	loggerOpts := []grpc_logrus.Option{
		grpc_logrus.WithDecider(func(fullMethodName string, err error) bool {
			return ignorePayload(fullMethodName)
		}),
		grpc_logrus.WithLevels(grpc_logrus.DefaultCodeToLevel),
	}

	return ilogger.ElasticsearchLoggerServerInterceptor(
		logrusEntry,
		ignorePayload,
		ignoreInitialRequest,
		loggerOpts...,
	)
}

// healthCheckWorker is running an infinite loop that sets the serving status once
// in s.healthCheckInterval seconds.
func (s SearchServer) healthCheckWorker(healthServer *health.Server) {
	for {
		if s.SearchService.HealthCheck() {
			healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
		} else {
			healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
		}

		time.Sleep(time.Second * time.Duration(s.healthCheckInterval))
	}
}
