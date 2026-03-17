package gootel

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	RedisOperationKey = attribute.Key("redis.operation")
	RedisKeyKey       = attribute.Key("redis.key")
	RedisDBKey        = attribute.Key("redis.db_index")
	RedisDurationKey  = attribute.Key("redis.duration_ms")
	RedisSuccessKey   = attribute.Key("redis.success")
	RedisErrorKey     = attribute.Key("redis.error")
	RedisFoundKey     = attribute.Key("redis.found")
)

// TraceRedisOperation creates a span for Redis operations
func TraceRedisOperation(ctx context.Context, operation, key string, dbIndex int) (context.Context, trace.Span) {
	tracer := otel.Tracer("redis")

	attrs := []attribute.KeyValue{
		RedisOperationKey.String(operation),
		RedisKeyKey.String(key),
		RedisDBKey.Int(dbIndex),
		SpanKindKey.String("client"),
	}

	return tracer.Start(ctx, "Redis "+operation, trace.WithAttributes(attrs...))
}

// RecordRedisSuccess records successful Redis operation with duration
func RecordRedisSuccess(span trace.Span, duration time.Duration, found bool) {
	span.SetAttributes(
		RedisSuccessKey.Bool(true),
		RedisDurationKey.Int64(duration.Milliseconds()),
		RedisFoundKey.Bool(found),
	)
	span.SetStatus(codes.Ok, "Redis operation succeeded")
}

// RecordRedisError records failed Redis operation with error details
func RecordRedisError(span trace.Span, err error) {
	span.SetAttributes(
		RedisSuccessKey.Bool(false),
		RedisErrorKey.String(err.Error()),
	)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}
