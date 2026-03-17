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
	HTTPMethodKey      = attribute.Key("http.method")
	HTTPURLKey         = attribute.Key("http.url")
	HTTPStatusCodeKey  = attribute.Key("http.status_code")
	HTTPServiceKey     = attribute.Key("http.service_name")
	HTTPDurationKey    = attribute.Key("http.duration_ms")
	HTTPErrorKey       = attribute.Key("http.error")
	HTTPTargetKey      = attribute.Key("http.target")
	SpanKindKey        = attribute.Key("span.kind")
)

// TraceHTTPRequest creates a span for HTTP requests with detailed attributes
func TraceHTTPRequest(ctx context.Context, serviceName, method, url string) (context.Context, trace.Span) {
	tracer := otel.Tracer("http-client")

	attrs := []attribute.KeyValue{
		HTTPMethodKey.String(method),
		HTTPURLKey.String(url),
		HTTPServiceKey.String(serviceName),
		SpanKindKey.String("client"),
	}

	return tracer.Start(ctx, "HTTP "+method, trace.WithAttributes(attrs...))
}

// RecordHTTPSuccess records successful HTTP request with status code and duration
func RecordHTTPSuccess(span trace.Span, statusCode int, duration time.Duration) {
	span.SetAttributes(
		HTTPStatusCodeKey.Int(statusCode),
		HTTPDurationKey.Int64(duration.Milliseconds()),
	)
	span.SetStatus(codes.Ok, "HTTP request succeeded")
}

// RecordHTTPError records failed HTTP request with error details
func RecordHTTPError(span trace.Span, err error) {
	span.SetAttributes(
		HTTPErrorKey.String(err.Error()),
	)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}
