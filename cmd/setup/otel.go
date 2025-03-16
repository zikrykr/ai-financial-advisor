package setup

import (
	"context"
	"log"

	"github.com/ai-financial-advisor/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func InitTracer() (*trace.TracerProvider, error) {
	ctx := context.Background()

	conf := config.GetConfig()

	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpoint(conf.Grafana.Host), // Grafana Tempo's OTLP gRPC endpoint
		otlptracegrpc.WithInsecure(),                  // Use TLS in production!
	)

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		log.Fatalf("Failed to create OTLP exporter: %v", err)
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(conf.App.Name),
		)),
	)

	otel.SetTracerProvider(tp)
	return tp, nil
}
