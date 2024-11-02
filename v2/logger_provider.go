package gootel

import (
	"context"
	"fmt"
	"time"

	"github.com/erajayatech/go-opentelemetry/v2/internal/config"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
)

// source: https://github.com/newrelic/newrelic-opentelemetry-examples/blob/main/getting-started-guides/go/otel.go#L114
// source: https://opentelemetry.io/docs/languages/go/instrumentation/#logs

// NewLoggerProvider return opentelemetry logger provider.
// Beaware that current version of go implementation for opentelemetry logger is experimental.
func NewLoggerProvider(ctx context.Context) (*log.LoggerProvider, error) {
	fail := func(err error, msg string) (*log.LoggerProvider, error) {
		return nil, fmt.Errorf("%s:: %w", msg, err)
	}

	// otelNRHost, err := config.GetOtelOTLPNewrelicHost()
	// if err != nil {
	// 	return fail(err, "error get otel otlp new relic host")
	// }

	otelNRHeaderAPIKey, err := config.GetOtelOTLPNewrelicHeaderAPIKey()
	if err != nil {
		return fail(err, "error get otel otlp new relic header api key")
	}

	opst := []otlploghttp.Option{
		otlploghttp.WithEndpointURL("https://otlp.nr-data.net:4317/v1/logs"),
		otlploghttp.WithHeaders(map[string]string{"api-key": otelNRHeaderAPIKey}),
		// otlploghttp.WithCompression(otlploghttp.GzipCompression),
	}

	exporter, err := otlploghttp.New(ctx, opst...)
	if err != nil {
		return fail(err, "error create otel otlp logger http exporter")
	}

	resource, err := getResource()
	if err != nil {
		return fail(err, "error get resource")
	}

	lp := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(exporter, log.WithExportInterval(time.Second))),
		log.WithResource(resource),
	)

	global.SetLoggerProvider(lp)

	return lp, nil
}
