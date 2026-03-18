package main

import (
	"context"
	"log"
	"os"
	"time"

	gootel "github.com/erajayatech/go-opentelemetry/v2"
	"go.opentelemetry.io/otel"
)

func main() {
	ctx := context.Background()

	config := gootel.ExporterConfig{
		ServiceName:    "datadog-example-service",
		ServiceVersion: "1.0.0",
		Environment:    "production",
	}

	datadogAPIKey := os.Getenv("DATADOG_API_KEY")
	if datadogAPIKey == "" {
		log.Fatal("DATADOG_API_KEY environment variable is required")
	}

	datadogEndpoint := os.Getenv("DATADOG_ENDPOINT")
	if datadogEndpoint == "" {
		datadogEndpoint = "trace-agent.datadoghq.com:4317"
	}

	useHTTP := os.Getenv("DATADOG_USE_HTTP") == "true"

	tp, err := gootel.NewTraceProviderWithDatadog(
		ctx,
		config,
		datadogAPIKey,
		datadogEndpoint,
		useHTTP,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down trace provider: %v", err)
		}
	}()

	otel.SetTracerProvider(tp)

	runExampleOperations(ctx)
}

func runExampleOperations(ctx context.Context) {
	ctx, span := gootel.RecordSpan(ctx)
	defer span.End()

	gootel.AddBusinessAttribute(span, "operation.type", "demo")
	gootel.AddBusinessAttribute(span, "environment", "production")

	processHTTPRequest(ctx)
	processDatabaseOperation(ctx)
	processCacheOperation(ctx)
	processBusinessLogic(ctx)
}

func processHTTPRequest(ctx context.Context) {
	ctx, span := gootel.TraceHTTPRequest(ctx, "datadog-example-service", "POST", "https://api.example.com/orders")
	defer span.End()

	gootel.AddBusinessAttribute(span, "http.category", "external-api")
	gootel.AddBusinessAttribute(span, "api.resource", "orders")

	startTime := time.Now()
	time.Sleep(75 * time.Millisecond)
	duration := time.Since(startTime)

	gootel.RecordHTTPSuccess(span, 201, duration)
	gootel.AddEventToSpan(ctx, "http.request.completed", map[string]interface{}{
		"response_size": 2048,
		"cache_hit":     false,
		"retry_count":   0,
	})
}

func processDatabaseOperation(ctx context.Context) {
	ctx, span := gootel.TraceDBOperation(ctx, "postgresql", "appdb", "INSERT", "INSERT INTO orders (user_id, amount, status) VALUES ($1, $2, $3)")
	defer span.End()

	gootel.AddBusinessAttribute(span, "db.operation", "insert")
	gootel.AddBusinessAttribute(span, "db.table", "orders")

	startTime := time.Now()
	time.Sleep(45 * time.Millisecond)
	duration := time.Since(startTime)

	gootel.RecordDBSuccess(span, 1, duration)
	gootel.RecordBusinessMetric(span, "db.rows.affected", 1)
	gootel.AddEventToSpan(ctx, "db.transaction.committed", map[string]interface{}{
		"transaction_id": "txn_123456",
		"isolation_level": "read_committed",
	})
}

func processCacheOperation(ctx context.Context) {
	ctx, span := gootel.TraceRedisOperation(ctx, "SET", "order:12345", 0)
	defer span.End()

	gootel.AddBusinessAttribute(span, "cache.key", "order:12345")
	gootel.AddBusinessAttribute(span, "cache.ttl", "3600")

	startTime := time.Now()
	time.Sleep(15 * time.Millisecond)
	duration := time.Since(startTime)

	gootel.RecordRedisSuccess(span, duration, true)
	gootel.AddEventToSpan(ctx, "cache.set.completed", map[string]interface{}{
		"key":      "order:12345",
		"ttl":      3600,
		"size":     1024,
		"location": "redis-cluster-01",
	})
}

func processBusinessLogic(ctx context.Context) {
	ctx, span := gootel.RecordSpan(ctx)
	defer span.End()

	gootel.AddBusinessAttribute(span, "business.process", "order.fulfillment")
	gootel.AddBusinessAttribute(span, "order.id", "ORD-12345")

	gootel.RecordBusinessMetric(span, "order.amount", 150.50)
	gootel.RecordBusinessMetric(span, "order.items.count", 3)

	gootel.AddBusinessEvent(span, "order.created", map[string]interface{}{
		"order_id":      "ORD-12345",
		"customer_id":   "CUST-67890",
		"amount":        150.50,
		"payment_method": "credit_card",
		"items_count":   3,
	})

	gootel.AddBusinessContext(span, "customer", map[string]string{
		"id":    "CUST-67890",
		"email": "customer@example.com",
		"tier":  "premium",
	})

	gootel.AddBusinessContext(span, "payment", map[string]string{
		"method":     "credit_card",
		"provider":   "stripe",
		"currency":   "USD",
		"status":     "approved",
	})
}
