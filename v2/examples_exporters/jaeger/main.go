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
		ServiceName:    "jaeger-example-service",
		ServiceVersion: "1.0.0",
		Environment:    "development",
	}

	jaegerEndpoint := os.Getenv("JAEGER_ENDPOINT")
	if jaegerEndpoint == "" {
		jaegerEndpoint = "localhost:4317"
	}

	tp, err := gootel.NewTraceProviderWithJaeger(
		ctx,
		config,
		jaegerEndpoint,
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
	gootel.AddBusinessAttribute(span, "environment", "development")

	processHTTPRequest(ctx)
	processDatabaseOperation(ctx)
	processCacheOperation(ctx)
}

func processHTTPRequest(ctx context.Context) {
	ctx, span := gootel.TraceHTTPRequest(ctx, "jaeger-example-service", "GET", "https://api.example.com/products")
	defer span.End()

	gootel.AddBusinessAttribute(span, "http.category", "external-api")
	gootel.AddBusinessAttribute(span, "api.resource", "products")

	startTime := time.Now()
	time.Sleep(60 * time.Millisecond)
	duration := time.Since(startTime)

	gootel.RecordHTTPSuccess(span, 200, duration)
	gootel.AddEventToSpan(ctx, "http.request.completed", map[string]interface{}{
		"response_size": 4096,
		"cache_hit":     true,
		"endpoint":      "/products",
	})
}

func processDatabaseOperation(ctx context.Context) {
	ctx, span := gootel.TraceDBOperation(ctx, "mysql", "appdb", "SELECT", "SELECT * FROM products WHERE category = $1")
	defer span.End()

	gootel.AddBusinessAttribute(span, "db.operation", "query")
	gootel.AddBusinessAttribute(span, "db.table", "products")

	startTime := time.Now()
	time.Sleep(35 * time.Millisecond)
	duration := time.Since(startTime)

	gootel.RecordDBSuccess(span, 25, duration)
	gootel.RecordBusinessMetric(span, "db.rows.returned", 25)
	gootel.RecordDBQueryEvent(span, "SELECT * FROM products WHERE category = 'electronics'")
}

func processCacheOperation(ctx context.Context) {
	ctx, span := gootel.TraceRedisOperation(ctx, "GET", "products:electronics", 0)
	defer span.End()

	gootel.AddBusinessAttribute(span, "cache.key", "products:electronics")
	gootel.AddBusinessAttribute(span, "cache.category", "products")

	startTime := time.Now()
	time.Sleep(12 * time.Millisecond)
	duration := time.Since(startTime)

	gootel.RecordRedisSuccess(span, duration, true)
	gootel.AddEventToSpan(ctx, "cache.hit", map[string]interface{}{
		"key":      "products:electronics",
		"ttl":      1800,
		"size":     8192,
		"location": "redis-cluster-02",
	})
}
