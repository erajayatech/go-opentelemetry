package gootel

import (
	"context"
	"fmt"

	"github.com/erajayatech/go-opentelemetry/v3/internal/config"
	sentryotel "github.com/getsentry/sentry-go/otel"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
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
		sdktrace.WithSpanProcessor(sentryotel.NewSentrySpanProcessor()),
	)

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(sentryotel.NewSentryPropagator())

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
