package goopentelemetry

import (
	"context"
	"runtime"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type (
	OtelTracer interface {
		SetTraceProviderJaeger() error
		SetTraceProviderNewRelic(ctx context.Context) error
	}

	otelTracer struct {
		env              string
		version          string
		service          string
		sampled          bool
		jaegerUrl        string
		exporterEndpoint *string
		apiKey           *string
	}
)

func ConstructOtelTracer(options ...OtelTracerOptionFunc) OtelTracer {
	otelTracerImplement := &otelTracer{
		env:     EnvironmentMode(),
		version: AppVersion(),
		service: AppName(),
		sampled: OtelSampled(),
	}

	// Run the options on it
	for _, option := range options {
		option(otelTracerImplement)
	}

	return otelTracerImplement
}

func (ot *otelTracer) SetTraceProviderJaeger() error {
	if ot.jaegerUrl == "" {
		ot.jaegerUrl = OtelJaegerURL()
	}

	tp, err := ot.tracerProviderJaeger()

	if err != nil {
		return err
	}

	otel.SetTracerProvider(tp)

	return nil
}

func (ot *otelTracer) SetTraceProviderNewRelic(context context.Context) error {
	tracerProvider, err := ot.tracerProviderNewRelic(context)
	if err != nil {
		return err
	}

	otel.SetTracerProvider(tracerProvider)

	return nil
}

// tracerProviderJaeger returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func (ot *otelTracer) tracerProviderJaeger() (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(ot.jaegerUrl)))
	if err != nil {
		return nil, err
	}

	var sampler = tracesdk.NeverSample()
	if ot.sampled {
		sampler = tracesdk.AlwaysSample()
	}

	traceProvider := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(sampler),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exporter),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(ot.service),
			semconv.ServiceVersionKey.String(ot.version),
			attribute.String("environment", ot.env),
		)),
	)

	return traceProvider, nil
}

// tracerProviderNewRelic returns an OpenTelemetry TracerProvider configured to use
// the NewRelic exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func (ot *otelTracer) tracerProviderNewRelic(ctx context.Context) (*tracesdk.TracerProvider, error) {
	// Create the NewRelic exporter
	var options []otlptracegrpc.Option

	if ot.exporterEndpoint != nil && ot.apiKey != nil {
		options = append(options, otlptracegrpc.WithEndpoint(*ot.exporterEndpoint))
		options = append(options, otlptracegrpc.WithHeaders(map[string]string{
			"api-key": *ot.apiKey,
		}))
		options = append(options, otlptracegrpc.WithCompressor("gzip"))
	}

	exporter, err := otlptrace.New(ctx, otlptracegrpc.NewClient(options...))
	if err != nil {
		return nil, err
	}

	var sampler = tracesdk.NeverSample()
	if ot.sampled {
		sampler = tracesdk.AlwaysSample()
	}

	traceProvider := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(sampler),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exporter),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(ot.service),
			semconv.ServiceVersionKey.String(ot.version),
			attribute.String("environment", ot.env),
		)),
	)

	return traceProvider, nil
}

func Start(ctx context.Context) (context.Context, trace.Span) {
	c, ok := ctx.(*gin.Context)
	if ok {
		requestMethod := c.Request.Method
		urlPath := c.FullPath()
		operation := WriteStringTemplate("[%s] %s %s", EnvironmentMode(), requestMethod, urlPath)
		return otel.Tracer("").Start(ctx, operation)
	}

	method, ok := grpc.Method(ctx)
	if ok {
		operation := WriteStringTemplate("[%s] %s", EnvironmentMode(), method)
		return otel.Tracer("").Start(ctx, operation)
	}

	return otel.Tracer("").Start(ctx, MyCaller(2))
}

func StartWorker(ctx context.Context) (context.Context, trace.Span) {
	actionName := GetActionName()
	operation := WriteStringTemplate("[%s] WORKER %s", EnvironmentMode(), GetFunctionName(2))

	return NewSpan(ctx, actionName, operation)
}

func MyCaller(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	details := runtime.FuncForPC(pc)

	if ok && details != nil {
		return details.Name()
	}
	return "failed to identify method caller"
}
