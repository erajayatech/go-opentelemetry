package gootel

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func setupTestTracerForHTTP() (*trace.TracerProvider, *tracetest.InMemoryExporter) {
	exporter := tracetest.NewInMemoryExporter()
	tp := trace.NewTracerProvider(trace.WithSyncer(exporter))
	return tp, exporter
}

func TestTraceHTTPRequest(t *testing.T) {
	tp, _ := setupTestTracerForHTTP()
	defer tp.Shutdown(context.Background())

	tracer := tp.Tracer("test")
	ctx := context.Background()
	ctx, span := tracer.Start(ctx, "test-span")

	ctx, httpSpan := TraceHTTPRequest(ctx, "test-service", "GET", "http://example.com/api")

	assert.NotNil(t, httpSpan)
	assert.True(t, httpSpan.SpanContext().IsValid())

	httpSpan.End()
	span.End()
}

func TestRecordHTTPSuccess(t *testing.T) {
	tp, _ := setupTestTracerForHTTP()
	defer tp.Shutdown(context.Background())

	tracer := tp.Tracer("test")
	ctx := context.Background()
	ctx, span := tracer.Start(ctx, "test-span")

	ctx, httpSpan := TraceHTTPRequest(ctx, "test-service", "GET", "http://example.com/api")

	RecordHTTPSuccess(httpSpan, 200, 100*time.Millisecond)

	httpSpan.End()
	span.End()
	assert.True(t, httpSpan.SpanContext().IsValid())
}

func TestRecordHTTPError(t *testing.T) {
	tp, _ := setupTestTracerForHTTP()
	defer tp.Shutdown(context.Background())

	tracer := tp.Tracer("test")
	ctx := context.Background()
	ctx, span := tracer.Start(ctx, "test-span")

	ctx, httpSpan := TraceHTTPRequest(ctx, "test-service", "POST", "http://example.com/api")

	testErr := errors.New("test error")
	RecordHTTPError(httpSpan, testErr)

	httpSpan.End()
	span.End()
	assert.True(t, httpSpan.SpanContext().IsValid())
}

func TestHTTPKeys(t *testing.T) {
	assert.NotNil(t, HTTPMethodKey)
	assert.NotNil(t, HTTPURLKey)
	assert.NotNil(t, HTTPStatusCodeKey)
	assert.NotNil(t, HTTPServiceKey)
	assert.NotNil(t, HTTPDurationKey)
	assert.NotNil(t, HTTPErrorKey)
	assert.NotNil(t, HTTPTargetKey)
	assert.NotNil(t, SpanKindKey)
}
