package gootel

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	BusinessErrorKey = attribute.Key("business.error")
	BusinessEventKey = attribute.Key("business.event")
)

// AddBusinessAttribute adds a custom business attribute to the current span
// This function is generic and can be used for any business domain (e-commerce, finance, logistics, etc.)
// Example: AddBusinessAttribute(span, "order.id", "12345")
// Example: AddBusinessAttribute(span, "customer.segment", "premium")
// Example: AddBusinessAttribute(span, "payment.method", "credit_card")
func AddBusinessAttribute(span trace.Span, key, value string) {
	span.SetAttributes(attribute.String(key, value))
}

// AddBusinessAttributes adds multiple custom business attributes to the current span at once
// This is more efficient than calling AddBusinessAttribute multiple times
// Example: AddBusinessAttributes(span, map[string]string{"order.id": "12345", "customer.id": "67890"})
func AddBusinessAttributes(span trace.Span, attrs map[string]string) {
	kv := make([]attribute.KeyValue, 0, len(attrs))
	for k, v := range attrs {
		kv = append(kv, attribute.String(k, v))
	}
	span.SetAttributes(kv...)
}

// AddBusinessEvent adds a business event to the current span
// This function supports various data types: string, int, float64, bool, int64
// Example: AddBusinessEvent(span, "order_placed", map[string]interface{}{"order.id": "12345", "total": 99.99})
// Example: AddBusinessEvent(span, "user_registered", map[string]interface{}{"user.id": "abc123", "source": "web"})
// Example: AddBusinessEvent(span, "payment_processed", map[string]interface{}{"payment.id": "pay_123", "success": true})
func AddBusinessEvent(span trace.Span, name string, data map[string]interface{}) {
	attrs := make([]attribute.KeyValue, 0, len(data))

	for k, v := range data {
		switch val := v.(type) {
		case string:
			attrs = append(attrs, attribute.String(k, val))
		case int:
			attrs = append(attrs, attribute.Int64(k, int64(val)))
		case int64:
			attrs = append(attrs, attribute.Int64(k, val))
		case float64:
			attrs = append(attrs, attribute.Float64(k, val))
		case bool:
			attrs = append(attrs, attribute.Bool(k, val))
		}
	}

	span.AddEvent(name, trace.WithAttributes(attrs...))
}

// RecordBusinessError records a business logic error with custom error type
// This is generic and can be used for any business domain
// Example: RecordBusinessError(span, "PAYMENT_FAILED", err)
// Example: RecordBusinessError(span, "INVENTORY_OUT_OF_STOCK", err)
// Example: RecordBusinessError(span, "VALIDATION_ERROR", err)
func RecordBusinessError(span trace.Span, errorType string, err error) {
	span.SetAttributes(
		BusinessErrorKey.String(errorType),
		attribute.String("error.message", err.Error()),
	)
	span.RecordError(err)
	span.SetStatus(codes.Error, "Business error: "+errorType)
}

// RecordBusinessMetric records a business metric/value to the current span
// This is useful for tracking KPIs, business metrics, and performance indicators
// Example: RecordBusinessMetric(span, "cart.value", 150.50)
// Example: RecordBusinessMetric(span, "conversion.rate", 0.75)
// Example: RecordBusinessMetric(span, "session.duration", 300)
func RecordBusinessMetric(span trace.Span, metricName string, value interface{}) {
	switch v := value.(type) {
	case string:
		span.SetAttributes(attribute.String(metricName, v))
	case int:
		span.SetAttributes(attribute.Int64(metricName, int64(v)))
	case int64:
		span.SetAttributes(attribute.Int64(metricName, v))
	case float64:
		span.SetAttributes(attribute.Float64(metricName, v))
	case bool:
		span.SetAttributes(attribute.Bool(metricName, v))
	}
}

// AddBusinessContext adds business context information to the current span
// This is useful for adding context that spans across multiple operations
// Example: AddBusinessContext(span, "order", map[string]string{"id": "12345", "status": "processing"})
// Example: AddBusinessContext(span, "user", map[string]string{"id": "abc123", "role": "admin"})
func AddBusinessContext(span trace.Span, contextType string, contextData map[string]string) {
	prefix := contextType + "."
	attrs := make([]attribute.KeyValue, 0, len(contextData))
	
	for k, v := range contextData {
		attrs = append(attrs, attribute.String(prefix+k, v))
	}
	
	span.SetAttributes(attrs...)
}
