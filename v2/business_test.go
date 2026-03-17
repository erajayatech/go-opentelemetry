package gootel

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func setupTestTracer() (*trace.TracerProvider, *tracetest.InMemoryExporter) {
	exporter := tracetest.NewInMemoryExporter()
	tp := trace.NewTracerProvider(trace.WithSyncer(exporter))
	return tp, exporter
}

func TestAddBusinessAttribute(t *testing.T) {
	tp, _ := setupTestTracer()
	defer tp.Shutdown(context.Background())

	tracer := tp.Tracer("test")
	_, span := tracer.Start(context.Background(), "test-span")

	AddBusinessAttribute(span, "business.process.id", "ORDER-12345")

	span.End()
	assert.True(t, span.SpanContext().IsValid())
}

func TestAddBusinessEvent(t *testing.T) {
	tp, _ := setupTestTracer()
	defer tp.Shutdown(context.Background())

	tracer := tp.Tracer("test")
	_, span := tracer.Start(context.Background(), "test-span")

	eventData := map[string]interface{}{
		"event.type":     "order.created",
		"customer.id":    "CUST-001",
		"order.amount":   100.50,
		"order.currency": "USD",
		"is_verified":    true,
	}

	AddBusinessEvent(span, "order.created", eventData)

	span.End()
	assert.True(t, span.SpanContext().IsValid())
}

func TestRecordBusinessError(t *testing.T) {
	tp, _ := setupTestTracer()
	defer tp.Shutdown(context.Background())

	tracer := tp.Tracer("test")
	_, span := tracer.Start(context.Background(), "test-span")

	testErr := errors.New("payment failed")
	RecordBusinessError(span, "PAYMENT_ERROR", testErr)

	span.End()
	assert.True(t, span.SpanContext().IsValid())
}

func TestBusinessKeys(t *testing.T) {
	assert.NotNil(t, BusinessErrorKey)
	assert.NotNil(t, BusinessEventKey)
}

func TestAddBusinessEventWithEmptyData(t *testing.T) {
	tp, _ := setupTestTracer()
	defer tp.Shutdown(context.Background())

	tracer := tp.Tracer("test")
	_, span := tracer.Start(context.Background(), "test-span")

	AddBusinessEvent(span, "test.event", map[string]interface{}{})

	span.End()
	assert.True(t, span.SpanContext().IsValid())
}

func TestAddBusinessEventWithNilData(t *testing.T) {
	tp, _ := setupTestTracer()
	defer tp.Shutdown(context.Background())

	tracer := tp.Tracer("test")
	_, span := tracer.Start(context.Background(), "test-span")

	AddBusinessEvent(span, "test.event", nil)

	span.End()
	assert.True(t, span.SpanContext().IsValid())
}

func TestAddBusinessAttributes(t *testing.T) {
	tp, _ := setupTestTracer()
	defer tp.Shutdown(context.Background())

	tracer := tp.Tracer("test")
	_, span := tracer.Start(context.Background(), "test-span")

	attrs := map[string]string{
		"order.id":     "ORDER-12345",
		"customer.id":  "CUST-001",
		"payment.method": "credit_card",
	}

	AddBusinessAttributes(span, attrs)

	span.End()
	assert.True(t, span.SpanContext().IsValid())
}

func TestRecordBusinessMetric(t *testing.T) {
	tp, _ := setupTestTracer()
	defer tp.Shutdown(context.Background())

	tracer := tp.Tracer("test")
	_, span := tracer.Start(context.Background(), "test-span")

	RecordBusinessMetric(span, "cart.value", 150.50)
	RecordBusinessMetric(span, "conversion.rate", 0.75)
	RecordBusinessMetric(span, "session.duration", 300)
	RecordBusinessMetric(span, "order.count", 42)
	RecordBusinessMetric(span, "is_premium_user", true)
	RecordBusinessMetric(span, "user.tier", "gold")

	span.End()
	assert.True(t, span.SpanContext().IsValid())
}

func TestAddBusinessContext(t *testing.T) {
	tp, _ := setupTestTracer()
	defer tp.Shutdown(context.Background())

	tracer := tp.Tracer("test")
	_, span := tracer.Start(context.Background(), "test-span")

	orderContext := map[string]string{
		"id":     "ORDER-12345",
		"status": "processing",
	}

	AddBusinessContext(span, "order", orderContext)

	userContext := map[string]string{
		"id":   "CUST-001",
		"role": "admin",
	}

	AddBusinessContext(span, "user", userContext)

	span.End()
	assert.True(t, span.SpanContext().IsValid())
}

func TestGenericBusinessAttributes(t *testing.T) {
	tp, _ := setupTestTracer()
	defer tp.Shutdown(context.Background())

	tracer := tp.Tracer("test")
	_, span := tracer.Start(context.Background(), "test-span")

	AddBusinessAttribute(span, "workflow.id", "WF-001")
	AddBusinessAttribute(span, "shipment.tracking", "TRK-12345")
	AddBusinessAttribute(span, "inventory.sku", "SKU-67890")
	AddBusinessAttribute(span, "transaction.ref", "TXN-ABC123")

	span.End()
	assert.True(t, span.SpanContext().IsValid())
}
