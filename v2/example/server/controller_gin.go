package main

import (
	"net/http"

	gootel "github.com/erajayatech/go-opentelemetry/v2"
	"github.com/gin-gonic/gin"
)

func controllerGinFoo(c *gin.Context) {
	ctx, span := gootel.RecordSpan(c)
	defer span.End()

	serviceFoo(ctx)

	c.JSON(http.StatusOK, gin.H{"trace_id": span.SpanContext().TraceID().String()})
}
