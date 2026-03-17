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

func setupTestTracerForRedis() (*trace.TracerProvider, *tracetest.InMemoryExporter) {
	exporter := tracetest.NewInMemoryExporter()
	tp := trace.NewTracerProvider(trace.WithSyncer(exporter))
	return tp, exporter
}

func TestTraceRedisOperation(t *testing.T) {
	tp, _ := setupTestTracerForRedis()
	defer tp.Shutdown(context.Background())

	tracer := tp.Tracer("test")
	ctx := context.Background()
	ctx, span := tracer.Start(ctx, "test-span")

	ctx, redisSpan := TraceRedisOperation(ctx, "GET", "user:123", 0)

	assert.NotNil(t, redisSpan)
	assert.True(t, redisSpan.SpanContext().IsValid())

	redisSpan.End()
	span.End()
}

func TestRecordRedisSuccess(t *testing.T) {
	tp, _ := setupTestTracerForRedis()
	defer tp.Shutdown(context.Background())

	tracer := tp.Tracer("test")
	ctx := context.Background()
	ctx, span := tracer.Start(ctx, "test-span")

	ctx, redisSpan := TraceRedisOperation(ctx, "SET", "user:123", 0)

	RecordRedisSuccess(redisSpan, 5*time.Millisecond, true)

	redisSpan.End()
	span.End()
	assert.True(t, redisSpan.SpanContext().IsValid())
}

func TestRecordRedisError(t *testing.T) {
	tp, _ := setupTestTracerForRedis()
	defer tp.Shutdown(context.Background())

	tracer := tp.Tracer("test")
	ctx := context.Background()
	ctx, span := tracer.Start(ctx, "test-span")

	ctx, redisSpan := TraceRedisOperation(ctx, "GET", "user:123", 0)

	testErr := errors.New("connection timeout")
	RecordRedisError(redisSpan, testErr)

	redisSpan.End()
	span.End()
	assert.True(t, redisSpan.SpanContext().IsValid())
}

func TestRedisKeys(t *testing.T) {
	assert.NotNil(t, RedisOperationKey)
	assert.NotNil(t, RedisKeyKey)
	assert.NotNil(t, RedisDBKey)
	assert.NotNil(t, RedisDurationKey)
	assert.NotNil(t, RedisSuccessKey)
	assert.NotNil(t, RedisErrorKey)
	assert.NotNil(t, RedisFoundKey)
}
