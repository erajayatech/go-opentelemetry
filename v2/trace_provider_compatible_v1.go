package gootel

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/sdk/trace"
)

// -----------------------------------------------------------
// -------- retained for compatibility with version 1 --------

type otelTracer struct {
	tp *trace.TracerProvider
}

// ConstructOtelTracer is retained for compatibility with version 1.
//
// Deprecated: Use NewTraceProvider instead. See v2/example/server/main.go
func ConstructOtelTracer() *otelTracer {
	return &otelTracer{}
}

// SetTraceProviderNewRelic is retained for compatibility with version 1.
//
// Deprecated: Use NewTraceProvider instead. See v2/example/server/main.go
func (o *otelTracer) SetTraceProviderNewRelic(ctx context.Context) error {
	fail := func(err error, msg string) error {
		return fmt.Errorf("%s:: %w", msg, err)
	}

	tp, err := NewTraceProvider(ctx)
	if err != nil {
		return fail(err, "error create new trace provider")
	}

	o.tp = tp

	return nil
}

func (o *otelTracer) Shutdown(ctx context.Context) error {
	return o.tp.Shutdown(ctx)
}

// -------- retained for compatibility with version 1 --------
// -----------------------------------------------------------
