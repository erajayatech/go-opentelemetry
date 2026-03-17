package gootel

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func setupTestTracerForSpan() (*trace.TracerProvider, *tracetest.InMemoryExporter) {
	exporter := tracetest.NewInMemoryExporter()
	tp := trace.NewTracerProvider(trace.WithSyncer(exporter))
	return tp, exporter
}

func TestRecordErrorToSpan(t *testing.T) {
	tp, _ := setupTestTracerForSpan()
	defer tp.Shutdown(context.Background())

	tracer := tp.Tracer("test")
	ctx, _ := tracer.Start(context.Background(), "test-span")

	testErr := errors.New("test error")
	RecordErrorToSpan(ctx, testErr)
}

func TestAddEventToSpan(t *testing.T) {
	tp, _ := setupTestTracerForSpan()
	defer tp.Shutdown(context.Background())

	tracer := tp.Tracer("test")
	ctx, _ := tracer.Start(context.Background(), "test-span")

	eventData := map[string]interface{}{
		"event.type":     "user.login",
		"user.id":        "123",
		"login.success":  true,
		"attempt.count":  3,
		"login.duration": 150.5,
	}

	AddEventToSpan(ctx, "login.attempt", eventData)
}

func TestGetTraceID(t *testing.T) {
	tp, _ := setupTestTracerForSpan()
	defer tp.Shutdown(context.Background())

	tracer := tp.Tracer("test")
	ctx, _ := tracer.Start(context.Background(), "test-span")

	traceID := GetTraceID(ctx)
	assert.NotEmpty(t, traceID)
}

func TestGetSpanID(t *testing.T) {
	tp, _ := setupTestTracerForSpan()
	defer tp.Shutdown(context.Background())

	tracer := tp.Tracer("test")
	ctx, _ := tracer.Start(context.Background(), "test-span")

	spanID := GetSpanID(ctx)
	assert.NotEmpty(t, spanID)
}

func TestAddExceptionToSpan(t *testing.T) {
	tp, _ := setupTestTracerForSpan()
	defer tp.Shutdown(context.Background())

	tracer := tp.Tracer("test")
	ctx, _ := tracer.Start(context.Background(), "test-span")

	stackTrace := "goroutine 1 [running]:\nmain.main()\n\t/path/to/file.go:10"
	AddExceptionToSpan(ctx, "NullPointerException", "something went wrong", stackTrace)
}

func TestRecordErrorToSpanWithNilContext(t *testing.T) {
	tp, _ := setupTestTracerForSpan()
	defer tp.Shutdown(context.Background())

	testErr := errors.New("test error")
	RecordErrorToSpan(context.Background(), testErr)
}

func TestAddEventToSpanWithNilContext(t *testing.T) {
	tp, _ := setupTestTracerForSpan()
	defer tp.Shutdown(context.Background())

	eventData := map[string]interface{}{
		"test.key": "test.value",
	}

	AddEventToSpan(context.Background(), "test.event", eventData)
}

func TestGetTraceIDWithNilContext(t *testing.T) {
	tp, _ := setupTestTracerForSpan()
	defer tp.Shutdown(context.Background())

	traceID := GetTraceID(context.Background())
	assert.Empty(t, traceID)
}

func TestGetSpanIDWithNilContext(t *testing.T) {
	tp, _ := setupTestTracerForSpan()
	defer tp.Shutdown(context.Background())

	spanID := GetSpanID(context.Background())
	assert.Empty(t, spanID)
}

func TestAddEventToSpanWithEmptyData(t *testing.T) {
	tp, _ := setupTestTracerForSpan()
	defer tp.Shutdown(context.Background())

	tracer := tp.Tracer("test")
	ctx, _ := tracer.Start(context.Background(), "test-span")

	AddEventToSpan(ctx, "test.event", map[string]interface{}{})
}

func TestAddExceptionToSpanWithEmptyStackTrace(t *testing.T) {
	tp, _ := setupTestTracerForSpan()
	defer tp.Shutdown(context.Background())

	tracer := tp.Tracer("test")
	ctx, _ := tracer.Start(context.Background(), "test-span")

	AddExceptionToSpan(ctx, "Error", "message", "")
}
