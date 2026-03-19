package tracer

import (
	"context"
	"manga-go/internal/pkg/config"
	"manga-go/internal/pkg/logger"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/encoding/gzip"
)

var dsn = os.Getenv("UPTRACE_DSN")

type Tracer struct {
	*otlptrace.Exporter
}

func InitTracer(cfg *config.Config, logger *logger.Logger) func(context.Context) error {
	secureOption := otlptracegrpc.WithInsecure()

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(cfg.Otlp.Endpoint),
			otlptracegrpc.WithHeaders(map[string]string{
				"uptrace-dsn": dsn,
			}),
			otlptracegrpc.WithCompressor(gzip.Name),
		),
	)

	if err != nil {
		logger.Error("Init tracer", err)
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", cfg.Service.Name),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		logger.Errorf("Could not create resources: %v", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)
	return exporter.Shutdown
}
