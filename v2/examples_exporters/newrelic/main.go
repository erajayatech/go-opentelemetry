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
		ServiceName:    "newrelic-example-service",
		ServiceVersion: "1.0.0",
		Environment:    "production",
	}

	newRelicAPIKey := os.Getenv("NEW_RELIC_API_KEY")
	if newRelicAPIKey == "" {
		log.Fatal("NEW_RELIC_API_KEY environment variable is required")
	}

	newRelicEndpoint := os.Getenv("NEW_RELIC_ENDPOINT")
	if newRelicEndpoint == "" {
		newRelicEndpoint = "otlp.nr-data.net:4317"
	}

	tp, err := gootel.NewTraceProviderWithNewRelic(
		ctx,
		config,
		newRelicAPIKey,
		newRelicEndpoint,
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
}

func processHTTPRequest(ctx context.Context) {
	ctx, span := gootel.TraceHTTPRequest(ctx, "newrelic-example-service", "GET", "https://api.example.com/users")
	defer span.End()

	gootel.AddBusinessAttribute(span, "http.category", "external-api")

	startTime := time.Now()
	time.Sleep(50 * time.Millisecond)
	duration := time.Since(startTime)

	gootel.RecordHTTPSuccess(span, 200, duration)
	gootel.AddEventToSpan(ctx, "http.request.completed", map[string]interface{}{
		"response_size": 1024,
		"cache_hit":     true,
	})
}

func processDatabaseOperation(ctx context.Context) {
	ctx, span := gootel.TraceDBOperation(ctx, "postgresql", "appdb", "SELECT", "SELECT * FROM users WHERE id = $1")
	defer span.End()

	gootel.AddBusinessAttribute(span, "db.operation", "query")
	gootel.AddBusinessAttribute(span, "db.table", "users")

	startTime := time.Now()
	time.Sleep(30 * time.Millisecond)
	duration := time.Since(startTime)

	gootel.RecordDBSuccess(span, 1, duration)
	gootel.RecordBusinessMetric(span, "db.rows.affected", 1)
}

func processCacheOperation(ctx context.Context) {
	ctx, span := gootel.TraceRedisOperation(ctx, "GET", "user:123", 0)
	defer span.End()

	gootel.AddBusinessAttribute(span, "cache.key", "user:123")

	startTime := time.Now()
	time.Sleep(10 * time.Millisecond)
	duration := time.Since(startTime)

	gootel.RecordRedisSuccess(span, duration, true)
	gootel.AddEventToSpan(ctx, "cache.hit", map[string]interface{}{
		"key":      "user:123",
		"ttl":      3600,
		"size":     512,
		"location": "redis-cluster-01",
	})
}
