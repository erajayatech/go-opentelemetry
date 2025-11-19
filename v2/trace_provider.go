package gootel

import (
	"context"
	"fmt"

	"github.com/erajayatech/go-opentelemetry/v2/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
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

	// Create the resource with common attributes
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

	// Setup trace exporters based on configuration
	var exporters []sdktrace.SpanExporter

	// Setup New Relic exporter if enabled
	if config.IsNewRelicEnabled() {
		opt, err := getNROption()
		if err != nil {
			return fail(err, "error get new relic option")
		}

		nrExporter, err := otlptracegrpc.New(ctx, opt...)
		if err != nil {
			return fail(err, "error create new relic otlp trace grpc exporter")
		}

		exporters = append(exporters, nrExporter)
	}

	// Setup Jaeger exporter if enabled
	if config.IsJaegerEnabled() {
		jaegerOpt, err := getJaegerOption()
		if err != nil {
			return fail(err, "error get jaeger option")
		}

		jaegerExporter, err := otlptracegrpc.New(ctx, jaegerOpt...)
		if err != nil {
			return fail(err, "error create jaeger otlp trace grpc exporter")
		}

		exporters = append(exporters, jaegerExporter)
	}

	// Setup Datadog exporter if enabled
	if config.IsDatadogEnabled() {
		datadogExporter, err := createDatadogExporter(ctx)
		if err != nil {
			return fail(err, "error create datadog exporter")
		}

		exporters = append(exporters, datadogExporter)
	}

	// Ensure we have at least one exporter
	if len(exporters) == 0 {
		return fail(fmt.Errorf("no exporters configured"), "error setting up trace provider")
	}

	// Create a BatchSpanProcessor for each exporter
	opts := []sdktrace.TracerProviderOption{
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(_resource),
	}

	for _, exporter := range exporters {
		opts = append(opts, sdktrace.WithBatcher(exporter))
	}

	// Create trace provider with all the options
	tp := sdktrace.NewTracerProvider(opts...)

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

func getJaegerOption() ([]otlptracegrpc.Option, error) {
	fail := func(err error, msg string) ([]otlptracegrpc.Option, error) {
		return nil, fmt.Errorf("%s:: %w", msg, err)
	}

	jaegerEndpoint, err := config.GetJaegerEndpoint()
	if err != nil {
		return fail(err, "error get jaeger endpoint")
	}

	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(jaegerEndpoint),
		otlptracegrpc.WithInsecure(), // Jaeger typically doesn't require authentication
	}

	return opts, nil
}

func createDatadogExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	fail := func(err error, msg string) (sdktrace.SpanExporter, error) {
		return nil, fmt.Errorf("%s:: %w", msg, err)
	}

	datadogEndpoint, err := config.GetDatadogEndpoint()
	if err != nil {
		return fail(err, "error get datadog endpoint")
	}

	datadogAPIKey, err := config.GetDatadogAPIKey()
	if err != nil {
		return fail(err, "error get datadog api key")
	}

	protocol := config.GetDatadogProtocol()

	if protocol == "http" {
		// HTTP exporter - construct proper HTTP endpoint
		httpEndpoint := datadogEndpoint
		opts := []otlptracehttp.Option{
			otlptracehttp.WithEndpoint(httpEndpoint),
			otlptracehttp.WithHeaders(map[string]string{"DD-API-KEY": datadogAPIKey}),
			otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
			otlptracehttp.WithInsecure(),
		}
		return otlptracehttp.New(ctx, opts...)
	} else {
		// gRPC exporter (default)
		opts := []otlptracegrpc.Option{
			otlptracegrpc.WithEndpoint(datadogEndpoint),
			otlptracegrpc.WithHeaders(map[string]string{"DD-API-KEY": datadogAPIKey}),
			otlptracegrpc.WithCompressor("gzip"),
			otlptracegrpc.WithInsecure(),
		}
		return otlptracegrpc.New(ctx, opts...)
	}
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
