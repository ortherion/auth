package app

import (
	"auth/internal/adapters/grpc"
	"auth/internal/adapters/repository/mongo"
	"auth/internal/adapters/rest"
	"auth/internal/config"
	"auth/internal/domain/service"
	"context"
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"golang.org/x/sync/errgroup"
	"net/http"
	"time"
)

var httpServer *rest.Server

func Start(ctx context.Context) {

	/* LOGGER INIT */
	log := logrus.New()
	log.SetFormatter(new(logrus.JSONFormatter))

	/* CONFIG INIT */
	cfg, err := config.NewConfigs()
	if err != nil {
		log.Fatal("can't load config")
	}

	/* DATABASE INIT */
	collection, err := mongo.NewClientMongoDB(ctx, cfg)
	if err != nil {
		log.Error(err)
	}

	//fileRepo := repository.NewFileRepo()
	mongoRepo := mongo.NewMongoRepo(collection)

	/* SERVICES INIT */
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:   "https://40376d00d7b8408f8bc64950ce173be9@sentry.k8s.golang-mts-teta.ru/49",
		Debug: true,
	}); err != nil {
		log.Error("can't init sentry")
	}
	defer sentry.Flush(2 * time.Second)

	/* Jaeger init */
	exporter, err := jaeger.New(jaeger.WithAgentEndpoint(jaeger.WithAgentHost(cfg.Jaeger.Host), jaeger.WithAgentPort(cfg.Jaeger.Port)))
	if err != nil {
		log.Error("can't init jaeger collector")
		sentry.CaptureException(err)
	}

	/* Tracer provider init */
	traceProvider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String(cfg.Jaeger.Service))))

	otel.SetTracerProvider(traceProvider)
	authService := service.NewAuthService(mongoRepo, cfg)

	/* REST INIT */
	helpers := rest.NewHelpers(log)
	router := rest.NewRouter(log, authService, helpers)

	/* SERVERS INIT */
	var g errgroup.Group

	g.Go(func() error {
		err := grpc.NewGrpcServer(authService).Start(cfg)
		return fmt.Errorf("failed creating grpc server. %w", err)
	})

	g.Go(func() error {
		httpServer = rest.NewServer(log, router, cfg)
		log.Printf("Start Server on %s", cfg.GetHTTPAddr())

		err = httpServer.Start()

		return fmt.Errorf("rest server was terminated with an error. %w", err)
	})

	g.Go(func() error {
		docServer := rest.NewServer(log, rest.SwaggerRouter(fmt.Sprintf("%s:%s", cfg.HTTP, cfg.HTTP.SwaggerPort)), cfg)
		err := docServer.Start()
		if err != nil {
			return err
		}
		return fmt.Errorf("swagger server was terminated with an error. %w", err)
	})

	err = g.Wait()

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err, "server start fail")
	}
}

func Stop() {
	logrus.Warn("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(2)*time.Second)
	defer cancel()

	err := httpServer.Stop(ctx)
	if err != nil {
		logrus.Error(err, "Error while stopping")
	}

	logrus.Warn("app has stopped")
}
