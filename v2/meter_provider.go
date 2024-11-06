package gootel

import (
	"context"
	"fmt"

	"github.com/erajayatech/go-opentelemetry/v2/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
)

// source: https://opentelemetry.io/docs/languages/go/instrumentation/#metrics
// source; https://opentelemetry.io/docs/languages/go/exporters/#otlp-metrics-over-grpc

// NewMeterProvider return opentelemetry meter provider.
func NewMeterProvider(ctx context.Context) (*metric.MeterProvider, error) {
	fail := func(err error, msg string) (*metric.MeterProvider, error) {
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

	opts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(otelNRHost),
		otlpmetricgrpc.WithHeaders(map[string]string{"api-key": otelNRHeaderAPIKey}),
		otlpmetricgrpc.WithCompressor("gzip"),
	}

	exporter, err := otlpmetricgrpc.New(ctx, opts...)
	if err != nil {
		return fail(err, "error create otel otlp metric grpc exporter")
	}

	resource, err := getResource()
	if err != nil {
		return fail(err, "error get resource")
	}

	mp := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter)),
		metric.WithResource(resource),
	)

	otel.SetMeterProvider(mp)

	return mp, nil
}
