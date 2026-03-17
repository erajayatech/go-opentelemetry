package gootel

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// RecordErrorToSpan records an error to the current span from context
func RecordErrorToSpan(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
}

// AddEventToSpan adds an event to the current span from context
func AddEventToSpan(ctx context.Context, name string, attributes map[string]interface{}) {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		return
	}

	attrs := make([]attribute.KeyValue, 0, len(attributes))

	for k, v := range attributes {
		switch val := v.(type) {
		case string:
			attrs = append(attrs, attribute.String(k, val))
		case int:
			attrs = append(attrs, attribute.Int64(k, int64(val)))
		case float64:
			attrs = append(attrs, attribute.Float64(k, val))
		case bool:
			attrs = append(attrs, attribute.Bool(k, val))
		}
	}

	span.AddEvent(name, trace.WithAttributes(attrs...))
}

// GetTraceID extracts the trace ID from context
func GetTraceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		return ""
	}

	spanContext := span.SpanContext()
	if !spanContext.IsValid() {
		return ""
	}

	return spanContext.TraceID().String()
}

// GetSpanID extracts the span ID from context
func GetSpanID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		return ""
	}

	spanContext := span.SpanContext()
	if !spanContext.IsValid() {
		return ""
	}

	return spanContext.SpanID().String()
}

// AddExceptionToSpan adds an exception event to the current span
func AddExceptionToSpan(ctx context.Context, exceptionType, message string, stackTrace string) {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		return
	}

	attrs := []attribute.KeyValue{
		attribute.String("exception.type", exceptionType),
		attribute.String("exception.message", message),
		attribute.String("exception.stacktrace", stackTrace),
	}

	span.AddEvent("exception", trace.WithAttributes(attrs...))
	span.SetStatus(codes.Error, message)
}
