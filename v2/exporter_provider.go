package gootel

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

// ExporterConfig holds configuration for trace providers
type ExporterConfig struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
}

// NewTraceProviderWithNewRelic creates a trace provider with New Relic exporter
func NewTraceProviderWithNewRelic(ctx context.Context, config ExporterConfig, apiKey, endpoint string) (*sdktrace.TracerProvider, error) {
	fail := func(err error, msg string) (*sdktrace.TracerProvider, error) {
		return nil, fmt.Errorf("%s:: %w", msg, err)
	}

	if endpoint == "" {
		endpoint = "otlp.nr-data.net:4317"
	}

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithHeaders(map[string]string{"api-key": apiKey}),
		otlptracegrpc.WithCompressor("gzip"),
	)
	if err != nil {
		return fail(err, "error create new relic exporter")
	}

	tp := createTraceProvider(exporter, config)
	return tp, nil
}

// NewTraceProviderWithDatadog creates a trace provider with Datadog exporter
func NewTraceProviderWithDatadog(ctx context.Context, config ExporterConfig, apiKey, endpoint string, useHTTP bool) (*sdktrace.TracerProvider, error) {
	fail := func(err error, msg string) (*sdktrace.TracerProvider, error) {
		return nil, fmt.Errorf("%s:: %w", msg, err)
	}

	if endpoint == "" {
		if useHTTP {
			endpoint = "trace-agent.datadoghq.com:4318"
		} else {
			endpoint = "trace-agent.datadoghq.com:4317"
		}
	}

	var exporter sdktrace.SpanExporter
	var err error

	if useHTTP {
		exporter, err = otlptracehttp.New(ctx,
			otlptracehttp.WithEndpoint(endpoint),
			otlptracehttp.WithHeaders(map[string]string{"DD-API-KEY": apiKey}),
			otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
			otlptracehttp.WithInsecure(),
		)
	} else {
		exporter, err = otlptracegrpc.New(ctx,
			otlptracegrpc.WithEndpoint(endpoint),
			otlptracegrpc.WithHeaders(map[string]string{"DD-API-KEY": apiKey}),
			otlptracegrpc.WithCompressor("gzip"),
			otlptracegrpc.WithInsecure(),
		)
	}

	if err != nil {
		return fail(err, "error create datadog exporter")
	}

	tp := createTraceProvider(exporter, config)
	return tp, nil
}

// NewTraceProviderWithJaeger creates a trace provider with Jaeger exporter
func NewTraceProviderWithJaeger(ctx context.Context, config ExporterConfig, endpoint string) (*sdktrace.TracerProvider, error) {
	fail := func(err error, msg string) (*sdktrace.TracerProvider, error) {
		return nil, fmt.Errorf("%s:: %w", msg, err)
	}

	if endpoint == "" {
		endpoint = "localhost:4317"
	}

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return fail(err, "error create jaeger exporter")
	}

	tp := createTraceProvider(exporter, config)
	return tp, nil
}

// NewTraceProviderWithOTLP creates a trace provider with generic OTLP exporter
func NewTraceProviderWithOTLP(ctx context.Context, config ExporterConfig, endpoint string, useHTTP bool, headers map[string]string) (*sdktrace.TracerProvider, error) {
	fail := func(err error, msg string) (*sdktrace.TracerProvider, error) {
		return nil, fmt.Errorf("%s:: %w", msg, err)
	}

	if endpoint == "" {
		if useHTTP {
			endpoint = "localhost:4318"
		} else {
			endpoint = "localhost:4317"
		}
	}

	var exporter sdktrace.SpanExporter
	var err error

	if useHTTP {
		opts := []otlptracehttp.Option{
			otlptracehttp.WithEndpoint(endpoint),
			otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
			otlptracehttp.WithInsecure(),
		}
		if len(headers) > 0 {
			opts = append(opts, otlptracehttp.WithHeaders(headers))
		}
		exporter, err = otlptracehttp.New(ctx, opts...)
	} else {
		opts := []otlptracegrpc.Option{
			otlptracegrpc.WithEndpoint(endpoint),
			otlptracegrpc.WithCompressor("gzip"),
			otlptracegrpc.WithInsecure(),
		}
		if len(headers) > 0 {
			opts = append(opts, otlptracegrpc.WithHeaders(headers))
		}
		exporter, err = otlptracegrpc.New(ctx, opts...)
	}

	if err != nil {
		return fail(err, "error create otlp exporter")
	}

	tp := createTraceProvider(exporter, config)
	return tp, nil
}

// NewTraceProviderWithStdout creates a trace provider with stdout exporter
func NewTraceProviderWithStdout(config ExporterConfig, prettyPrint bool) (*sdktrace.TracerProvider, error) {
	fail := func(err error, msg string) (*sdktrace.TracerProvider, error) {
		return nil, fmt.Errorf("%s:: %w", msg, err)
	}

	opts := []stdouttrace.Option{}
	if prettyPrint {
		opts = append(opts, stdouttrace.WithPrettyPrint())
	}

	exporter, err := stdouttrace.New(opts...)
	if err != nil {
		return fail(err, "error create stdout exporter")
	}

	tp := createTraceProvider(exporter, config)
	return tp, nil
}

// NewTraceProviderWithMultipleExporters creates a trace provider with multiple exporters
func NewTraceProviderWithMultipleExporters(ctx context.Context, config ExporterConfig, exporters []sdktrace.SpanExporter) (*sdktrace.TracerProvider, error) {
	fail := func(err error, msg string) (*sdktrace.TracerProvider, error) {
		return nil, fmt.Errorf("%s:: %w", msg, err)
	}

	if len(exporters) == 0 {
		return fail(fmt.Errorf("no exporters provided"), "error create trace provider")
	}

	opts := []sdktrace.TracerProviderOption{
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(createResource(config)),
	}

	for _, exporter := range exporters {
		opts = append(opts, sdktrace.WithBatcher(exporter))
	}

	tp := sdktrace.NewTracerProvider(opts...)
	return tp, nil
}

func createTraceProvider(exporter sdktrace.SpanExporter, config ExporterConfig) *sdktrace.TracerProvider {
	return sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(createResource(config)),
		sdktrace.WithBatcher(exporter),
	)
}

func createResource(config ExporterConfig) *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(config.ServiceName),
		semconv.ServiceVersionKey.String(config.ServiceVersion),
		semconv.DeploymentEnvironmentKey.String(config.Environment),
	)
}
