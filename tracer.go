package goopentelemetry

import (
	"context"
	"fmt"
	"runtime"
	"strings"

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
)

type OtelTracer struct {
	env     string
	version string
	service string
	sampled bool
}

func ConstructOtelTracer() OtelTracer {
	return OtelTracer{
		env:     GetEnv("MODE"),
		version: GetEnv("APP_VERSION"),
		service: GetEnv("APP_NAME"),
		sampled: StringToBool(GetEnv("OTEL_SAMPLED")),
	}
}

// tracerProviderJaeger returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func tracerProviderJaeger(url string, service string, version string, env string, sampled bool) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}

	var sampler = tracesdk.NeverSample()
	if sampled {
		sampler = tracesdk.AlwaysSample()
	}

	traceProvider := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(sampler),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exporter),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			semconv.ServiceVersionKey.String(version),
			attribute.String("environment", env),
		)),
	)

	return traceProvider, nil
}

func (otelTracer *OtelTracer) SetTraceProviderJaeger() error {
	env := otelTracer.env
	version := otelTracer.version
	sampled := otelTracer.sampled
	service := GetEnv("APP_NAME")
	jaegerUrl := GetEnv("OTEL_JAEGER_URL")

	tp, err := tracerProviderJaeger(jaegerUrl, service, version, env, sampled)

	if err != nil {
		return err
	}

	otel.SetTracerProvider(tp)

	return nil
}

func (otelTracer *OtelTracer) SetTraceProviderNewRelic(context context.Context) error {
	env := otelTracer.env
	version := otelTracer.version
	sampled := otelTracer.sampled
	service := GetEnv("APP_NAME")

	tracerProvider, err := tracerProviderNewRelic(context, service, version, env, sampled)
	if err != nil {
		return err
	}

	otel.SetTracerProvider(tracerProvider)

	return nil
}

// tracerProviderNewRelic returns an OpenTelemetry TracerProvider configured to use
// the NewRelic exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func tracerProviderNewRelic(ctx context.Context, service string, version string, env string, sampled bool) (*tracesdk.TracerProvider, error) {
	// Create the NewRelic exporter
	exporter, err := otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(),
	)
	if err != nil {
		return nil, err
	}

	var sampler = tracesdk.NeverSample()
	if sampled {
		sampler = tracesdk.AlwaysSample()
	}

	traceProvider := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(sampler),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exporter),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			semconv.ServiceVersionKey.String(version),
			attribute.String("environment", env),
		)),
	)

	return traceProvider, nil
}

func Start(ctx *gin.Context) (context.Context, trace.Span) {
	actionName := getActionName()
	request := ctx.Request
	requestMethod := request.Method
	urlPath := ctx.Request.URL.Path
	env := GetEnv("MODE")

	operation := fmt.Sprintf("[%s] %s %s", env, requestMethod, urlPath)

	return NewSpan(ctx, actionName, operation)
}

func getActionName() string {
	c, _, _, _ := runtime.Caller(1)
	f := runtime.FuncForPC(c).Name()
	fs := strings.SplitN(f, ".", 2)
	replacer := strings.NewReplacer("(", "", ")", "", "*", "")
	actionName := replacer.Replace(fs[1])

	return actionName
}
