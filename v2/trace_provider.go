package gootel

import (
	"context"
	"fmt"

	"github.com/erajayatech/go-opentelemetry/v2/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

// source: https://opentelemetry.io/docs/languages/go/instrumentation/#traces
// source: https://opentelemetry.io/docs/languages/go/exporters/#otlp-traces-over-grpc

// NewTraceProvider return opentelemetry trace provider.
func NewTraceProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	fail := func(err error, msg string) (*sdktrace.TracerProvider, error) {
		return nil, fmt.Errorf("%s:: %w", msg, err)
	}

	opt, err := getNROption()
	if err != nil {
		return fail(err, "error get new relic option")
	}

	exporter, err := otlptracegrpc.New(ctx, opt...)
	if err != nil {
		return fail(err, "error create otel otlp trace grpc exporter")
	}

	serviceName, err := config.GetServiceName()
	if err != nil {
		return fail(err, "error get service name")
	}

	appVersion, err := config.GetAppVersion()
	if err != nil {
		return fail(err, "error get app version")
	}

	appEnv, err := config.GetAppEnvironment()
	if err != nil {
		return fail(err, "error get app environment")
	}

	_resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(appVersion),
		attribute.String("environment", appEnv),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(_resource),
	)

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp, nil
}

func getNROption() ([]otlptracegrpc.Option, error) {
	fail := func(err error, msg string) ([]otlptracegrpc.Option, error) {
		return nil, fmt.Errorf("%s:: %w", msg, err)
	}

	otelNRHost, err := config.GetOtelOTLPNewrelicHost()
	if err != nil {
		return fail(err, "error get otel otlp new relic host")
	}

	otelNRHeaderAPIKey, err := config.GetOtelOTLPNewrelicHeaderAPIKey()
	if err != nil {
		return fail(err, "error get otel otlp new relic header api key")
	}

	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(otelNRHost),
		otlptracegrpc.WithHeaders(map[string]string{"api-key": otelNRHeaderAPIKey}),
		otlptracegrpc.WithCompressor("gzip"),
	}

	return opts, nil
}

// -----------------------------------------------------------
// -------- retained for compatibility with version 1 --------

type otelTracer struct {
	tp *sdktrace.TracerProvider
}

// ConstructOtelTracer is retained for compatibility with version 1.
//
// Deprecated: Use NewTraceProvider instead. See v2/example/server/main.go
func ConstructOtelTracer() *otelTracer {
	return &otelTracer{}
}

// SetTraceProviderNewRelic is retained for compatibility with version 1.
//
// Deprecated: Use NewTraceProvider instead. See v2/example/server/main.go
func (o *otelTracer) SetTraceProviderNewRelic(ctx context.Context) error {
	fail := func(err error, msg string) error {
		return fmt.Errorf("%s:: %w", msg, err)
	}

	tp, err := NewTraceProvider(ctx)
	if err != nil {
		return fail(err, "error create new trace provider")
	}

	o.tp = tp

	return nil
}

// -------- retained for compatibility with version 1 --------
// -----------------------------------------------------------

func (o *otelTracer) Shutdown(ctx context.Context) error {
	return o.tp.Shutdown(ctx)
}
