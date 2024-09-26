package gootel

import (
	"context"

	"github.com/erajayatech/go-opentelemetry/v3/internal/caller"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// RecordSpan to record span.
func RecordSpan(ctx context.Context) (context.Context, trace.Span) {
	if c, ok := ctx.(*gin.Context); ok {
		return otel.Tracer("").Start(c.Request.Context(), caller.FuncName(caller.WithSkip(1)))
	}
	return otel.Tracer("").Start(ctx, caller.FuncName(caller.WithSkip(1)))
}
