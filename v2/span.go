package gootel

import (
	"context"

	"github.com/erajayatech/go-opentelemetry/v2/internal/caller"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.opentelemetry.io/otel/trace"
)

// RecordSpan to record span.
func RecordSpan(ctx context.Context) (context.Context, trace.Span) {
	if c, ok := ctx.(*gin.Context); ok {
		return otel.Tracer("").Start(c.Request.Context(), caller.FuncName(caller.WithSkip(1)))
	}
	return otel.Tracer("").Start(ctx, caller.FuncName(caller.WithSkip(1)))
}

// -----------------------------------------------------------
// -------- retained for compatibility with version 1 --------

// Start is retained for compatibility with version 1.
//
// Deprecated: Use RecordSpan instead.
func Start(ctx context.Context) (context.Context, trace.Span) {
	if c, ok := ctx.(*gin.Context); ok {
		return otel.Tracer("").Start(c.Request.Context(), caller.FuncName(caller.WithSkip(1)))
	}
	return otel.Tracer("").Start(ctx, caller.FuncName(caller.WithSkip(1)))
}

// NewSpan is retained for compatibility with version 1.
//
// Deprecated: Use RecordSpan instead.
func NewSpan(ctx context.Context, args ...string) (context.Context, trace.Span) {
	name := caller.FuncName(caller.WithSkip(1))
	if len(args) > 1 {
		name = args[0]
	}
	return otel.Tracer("").Start(ctx, name)
}

// StartWorker is retained for compatibility with version 1.
//
// Deprecated: Use RecordSpan instead.
func StartWorker(ctx context.Context) (context.Context, trace.Span) {
	return otel.Tracer("").Start(ctx, caller.FuncName(caller.WithSkip(1)))
}

// AddSpanError is retained for compatibility with version 1.
//
// Deprecated: Use span.RecordError(err) instead.
func AddSpanError(span trace.Span, err error) {
	span.RecordError(err)
}

// SetStatus is retained for compatibility with version 1.
//
// Deprecated: Use span.SetStatus(codes.Error, msg) instead.
func FailSpan(span trace.Span, msg string) {
	span.SetStatus(codes.Error, msg)
}

// AddSpanTags is retained for compatibility with version 1.
//
// Deprecated: Use span.SetAttributes instead.
func AddSpanTags(span trace.Span, tags map[string]string) {
	list := make([]attribute.KeyValue, len(tags))

	var i int
	for k, v := range tags {
		list[i] = attribute.Key(k).String(v)
		i++
	}

	span.SetAttributes(list...)
}

// AddSpanEvents is retained for compatibility with version 1.
//
// Deprecated: Use span.AddEvent instead.
func AddSpanEvents(span trace.Span, name string, events map[string]string) {
	list := make([]trace.EventOption, len(events))

	var i int
	for k, v := range events {
		list[i] = trace.WithAttributes(attribute.Key(k).String(v))
		i++
	}

	span.AddEvent(name, list...)
}

// SpanFromContext is retained for compatibility with version 1.
//
// Deprecated: Use trace.SpanFromContext(ctx) instead.
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

type HttpSpanAttribute struct {
	Method     string
	Version    string
	Url        string
	IP         string
	StatusCode int
}

var httpFlavorKey = "1.0"

// AddSpanTags is retained for compatibility with version 1.
//
// Deprecated: Use RecordSpan instead.
func NewHttpSpan(ctx context.Context, name string, operation string, httpSpanAttribute HttpSpanAttribute) (context.Context, trace.Span) {
	return otel.Tracer(name).Start(
		ctx,
		operation,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			semconv.HTTPMethodKey.String(httpSpanAttribute.Method),
			semconv.HTTPFlavorKey.String(httpFlavorKey),
			semconv.HTTPURLKey.String(httpSpanAttribute.Url),
			attribute.Key("net.peer.ip").String(httpSpanAttribute.IP),
		),
	)
}

// -------- retained for compatibility with version 1 --------
// -----------------------------------------------------------
