package gootel

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestNewTraceProviderWithStdout(t *testing.T) {
	config := ExporterConfig{
		ServiceName:    "test-service",
		ServiceVersion: "1.0.0",
		Environment:    "test",
	}

	tp, err := NewTraceProviderWithStdout(config, false)
	assert.NoError(t, err)
	assert.NotNil(t, tp)

	ctx := context.Background()
	err = tp.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestNewTraceProviderWithStdoutPrettyPrint(t *testing.T) {
	config := ExporterConfig{
		ServiceName:    "test-service",
		ServiceVersion: "1.0.0",
		Environment:    "test",
	}

	tp, err := NewTraceProviderWithStdout(config, true)
	assert.NoError(t, err)
	assert.NotNil(t, tp)

	ctx := context.Background()
	err = tp.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestNewTraceProviderWithJaeger(t *testing.T) {
	config := ExporterConfig{
		ServiceName:    "test-service",
		ServiceVersion: "1.0.0",
		Environment:    "test",
	}

	tp, err := NewTraceProviderWithJaeger(context.Background(), config, "localhost:4317")
	assert.NoError(t, err)
	assert.NotNil(t, tp)

	ctx := context.Background()
	err = tp.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestNewTraceProviderWithOTLP(t *testing.T) {
	config := ExporterConfig{
		ServiceName:    "test-service",
		ServiceVersion: "1.0.0",
		Environment:    "test",
	}

	t.Run("gRPC", func(t *testing.T) {
		tp, err := NewTraceProviderWithOTLP(context.Background(), config, "localhost:4317", false, nil)
		assert.NoError(t, err)
		assert.NotNil(t, tp)

		ctx := context.Background()
		err = tp.Shutdown(ctx)
		assert.NoError(t, err)
	})

	t.Run("HTTP", func(t *testing.T) {
		tp, err := NewTraceProviderWithOTLP(context.Background(), config, "localhost:4318", true, nil)
		assert.NoError(t, err)
		assert.NotNil(t, tp)

		ctx := context.Background()
		err = tp.Shutdown(ctx)
		assert.NoError(t, err)
	})

	t.Run("With Headers", func(t *testing.T) {
		headers := map[string]string{
			"Authorization": "Bearer token",
			"X-Custom":      "value",
		}
		tp, err := NewTraceProviderWithOTLP(context.Background(), config, "localhost:4317", false, headers)
		assert.NoError(t, err)
		assert.NotNil(t, tp)

		ctx := context.Background()
		err = tp.Shutdown(ctx)
		assert.NoError(t, err)
	})
}

func TestNewTraceProviderWithMultipleExporters(t *testing.T) {
	config := ExporterConfig{
		ServiceName:    "test-service",
		ServiceVersion: "1.0.0",
		Environment:    "test",
	}

	exporters := []sdktrace.SpanExporter{
		tracetest.NewInMemoryExporter(),
		tracetest.NewInMemoryExporter(),
	}

	tp, err := NewTraceProviderWithMultipleExporters(context.Background(), config, exporters)
	assert.NoError(t, err)
	assert.NotNil(t, tp)

	ctx := context.Background()
	err = tp.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestNewTraceProviderWithMultipleExportersEmpty(t *testing.T) {
	config := ExporterConfig{
		ServiceName:    "test-service",
		ServiceVersion: "1.0.0",
		Environment:    "test",
	}

	exporters := []sdktrace.SpanExporter{}

	tp, err := NewTraceProviderWithMultipleExporters(context.Background(), config, exporters)
	assert.Error(t, err)
	assert.Nil(t, tp)
}

func TestExportProviderConfig(t *testing.T) {
	config := ExporterConfig{
		ServiceName:    "my-service",
		ServiceVersion: "2.0.0",
		Environment:    "production",
	}

	assert.Equal(t, "my-service", config.ServiceName)
	assert.Equal(t, "2.0.0", config.ServiceVersion)
	assert.Equal(t, "production", config.Environment)
}

func TestNewTraceProviderWithNewRelic(t *testing.T) {
	config := ExporterConfig{
		ServiceName:    "test-service",
		ServiceVersion: "1.0.0",
		Environment:    "test",
	}

	tp, err := NewTraceProviderWithNewRelic(context.Background(), config, "test-api-key", "otlp.nr-data.net:4317")
	assert.NoError(t, err)
	assert.NotNil(t, tp)

	ctx := context.Background()
	err = tp.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestNewTraceProviderWithDatadog(t *testing.T) {
	config := ExporterConfig{
		ServiceName:    "test-service",
		ServiceVersion: "1.0.0",
		Environment:    "test",
	}

	t.Run("gRPC", func(t *testing.T) {
		tp, err := NewTraceProviderWithDatadog(context.Background(), config, "test-api-key", "trace-agent.datadoghq.com:4317", false)
		assert.NoError(t, err)
		assert.NotNil(t, tp)

		ctx := context.Background()
		err = tp.Shutdown(ctx)
		assert.NoError(t, err)
	})

	t.Run("HTTP", func(t *testing.T) {
		tp, err := NewTraceProviderWithDatadog(context.Background(), config, "test-api-key", "trace-agent.datadoghq.com:4318", true)
		assert.NoError(t, err)
		assert.NotNil(t, tp)

		ctx := context.Background()
		err = tp.Shutdown(ctx)
		assert.NoError(t, err)
	})
}

func TestExportProviderDefaults(t *testing.T) {
	config := ExporterConfig{
		ServiceName:    "test-service",
		ServiceVersion: "1.0.0",
		Environment:    "test",
	}

	t.Run("NewRelic default endpoint", func(t *testing.T) {
		tp, err := NewTraceProviderWithNewRelic(context.Background(), config, "test-api-key", "")
		assert.NoError(t, err)
		assert.NotNil(t, tp)

		ctx := context.Background()
		err = tp.Shutdown(ctx)
		assert.NoError(t, err)
	})

	t.Run("Jaeger default endpoint", func(t *testing.T) {
		tp, err := NewTraceProviderWithJaeger(context.Background(), config, "")
		assert.NoError(t, err)
		assert.NotNil(t, tp)

		ctx := context.Background()
		err = tp.Shutdown(ctx)
		assert.NoError(t, err)
	})

	t.Run("OTLP default endpoint gRPC", func(t *testing.T) {
		tp, err := NewTraceProviderWithOTLP(context.Background(), config, "", false, nil)
		assert.NoError(t, err)
		assert.NotNil(t, tp)

		ctx := context.Background()
		err = tp.Shutdown(ctx)
		assert.NoError(t, err)
	})

	t.Run("OTLP default endpoint HTTP", func(t *testing.T) {
		tp, err := NewTraceProviderWithOTLP(context.Background(), config, "", true, nil)
		assert.NoError(t, err)
		assert.NotNil(t, tp)

		ctx := context.Background()
		err = tp.Shutdown(ctx)
		assert.NoError(t, err)
	})

	t.Run("Datadog default endpoint gRPC", func(t *testing.T) {
		tp, err := NewTraceProviderWithDatadog(context.Background(), config, "test-api-key", "", false)
		assert.NoError(t, err)
		assert.NotNil(t, tp)

		ctx := context.Background()
		err = tp.Shutdown(ctx)
		assert.NoError(t, err)
	})

	t.Run("Datadog default endpoint HTTP", func(t *testing.T) {
		tp, err := NewTraceProviderWithDatadog(context.Background(), config, "test-api-key", "", true)
		assert.NoError(t, err)
		assert.NotNil(t, tp)

		ctx := context.Background()
		err = tp.Shutdown(ctx)
		assert.NoError(t, err)
	})
}
