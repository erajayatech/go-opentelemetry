package gootel

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func setupGlobalTracer() (*trace.TracerProvider, *tracetest.InMemoryExporter) {
	exporter := tracetest.NewInMemoryExporter()
	tp := trace.NewTracerProvider(trace.WithSyncer(exporter))
	otel.SetTracerProvider(tp)
	return tp, exporter
}

func TestRecordSpanBackwardCompatibility(t *testing.T) {
	tp, _ := setupGlobalTracer()
	defer tp.Shutdown(context.Background())

	ctx := context.Background()
	ctx, span := RecordSpan(ctx)

	assert.NotNil(t, span)
	assert.NotNil(t, ctx)
	assert.True(t, span.SpanContext().IsValid())

	span.End()
}

func TestRecordSpanWithTracer(t *testing.T) {
	tp, _ := setupGlobalTracer()
	defer tp.Shutdown(context.Background())

	ctx := context.Background()
	ctx, span := RecordSpan(ctx)

	assert.NotNil(t, span)
	assert.NotNil(t, ctx)

	span.SetAttributes()
	span.End()
}

func TestEnhancedHelpersBackwardCompatibility(t *testing.T) {
	tp, _ := setupGlobalTracer()
	defer tp.Shutdown(context.Background())

	ctx := context.Background()

	ctx, span := RecordSpan(ctx)
	assert.NotNil(t, span)
	assert.True(t, span.SpanContext().IsValid())

	ctx, httpSpan := TraceHTTPRequest(ctx, "test-service", "GET", "http://example.com")
	assert.NotNil(t, httpSpan)
	assert.True(t, httpSpan.SpanContext().IsValid())

	httpSpan.End()
	span.End()
}

func TestBusinessHelpersBackwardCompatibility(t *testing.T) {
	tp, _ := setupGlobalTracer()
	defer tp.Shutdown(context.Background())

	ctx := context.Background()
	ctx, span := RecordSpan(ctx)
	assert.NotNil(t, span)

	AddBusinessAttribute(span, "test.key", "test.value")
	RecordBusinessError(span, "TEST_ERROR", assert.AnError)

	span.End()
}

func TestEnhancedSpanHelpersBackwardCompatibility(t *testing.T) {
	tp, _ := setupGlobalTracer()
	defer tp.Shutdown(context.Background())

	ctx := context.Background()
	ctx, span := RecordSpan(ctx)
	assert.NotNil(t, span)

	RecordErrorToSpan(ctx, assert.AnError)
	AddEventToSpan(ctx, "test.event", map[string]interface{}{"key": "value"})

	traceID := GetTraceID(ctx)
	assert.NotEmpty(t, traceID)

	spanID := GetSpanID(ctx)
	assert.NotEmpty(t, spanID)

	span.End()
}
