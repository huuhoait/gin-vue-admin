package core

import (
	"context"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.uber.org/zap"

	"github.com/huuhoait/gin-vue-admin/server/global"
)

// InitTracer configures the global OpenTelemetry TracerProvider.
// It reads OTEL_EXPORTER_OTLP_ENDPOINT (default: localhost:4317) and
// OTEL_SERVICE_NAME (default: gin-vue-admin) from the environment so that
// the same binary works in dev (Jaeger all-in-one) and prod (OTEL Collector).
//
// Returns a shutdown function that must be deferred in main().
func InitTracer() func(context.Context) error {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:4317"
	}
	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "gin-vue-admin"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(), // TLS controlled by OTEL_EXPORTER_OTLP_INSECURE env
	)
	if err != nil {
		if global.GVA_LOG != nil {
			global.GVA_LOG.Warn("otel exporter init failed — tracing disabled", zap.Error(err))
		}
		return func(context.Context) error { return nil }
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(sampleRate()))),
	)
	otel.SetTracerProvider(tp)

	if global.GVA_LOG != nil {
		global.GVA_LOG.Info("otel tracer initialized", zap.String("endpoint", endpoint), zap.String("service", serviceName))
	}
	return tp.Shutdown
}

// sampleRate returns 1.0 in dev/test and 0.1 in release to keep span volume reasonable.
func sampleRate() float64 {
	if os.Getenv("GIN_MODE") == "release" {
		return 0.1
	}
	return 1.0
}
