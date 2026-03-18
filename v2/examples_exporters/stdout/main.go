package main

import (
	"context"
	"log"
	"time"

	gootel "github.com/erajayatech/go-opentelemetry/v2"
	"go.opentelemetry.io/otel"
)

func main() {
	config := gootel.ExporterConfig{
		ServiceName:    "stdout-example-service",
		ServiceVersion: "1.0.0",
		Environment:    "development",
	}

	tp, err := gootel.NewTraceProviderWithStdout(config, true)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down trace provider: %v", err)
		}
	}()

	otel.SetTracerProvider(tp)

	runExampleOperations(context.Background())
}

func runExampleOperations(ctx context.Context) {
	ctx, span := gootel.RecordSpan(ctx)
	defer span.End()

	gootel.AddBusinessAttribute(span, "operation.type", "demo")
	gootel.AddBusinessAttribute(span, "environment", "development")

	processHTTPRequest(ctx)
	processDatabaseOperation(ctx)
	processCacheOperation(ctx)
	processBusinessLogic(ctx)
}

func processHTTPRequest(ctx context.Context) {
	ctx, span := gootel.TraceHTTPRequest(ctx, "stdout-example-service", "GET", "https://api.example.com/users")
	defer span.End()

	gootel.AddBusinessAttribute(span, "http.category", "external-api")
	gootel.AddBusinessAttribute(span, "http.resource", "users")

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
	gootel.RecordDBQueryEvent(span, "SELECT * FROM users WHERE id = 'user-123'")
}

func processCacheOperation(ctx context.Context) {
	ctx, span := gootel.TraceRedisOperation(ctx, "GET", "user:123", 0)
	defer span.End()

	gootel.AddBusinessAttribute(span, "cache.key", "user:123")
	gootel.AddBusinessAttribute(span, "cache.ttl", "3600")

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

func processBusinessLogic(ctx context.Context) {
	ctx, span := gootel.RecordSpan(ctx)
	defer span.End()

	gootel.AddBusinessAttribute(span, "business.process", "user.profile.update")
	gootel.AddBusinessAttribute(span, "user.id", "USER-123")

	gootel.RecordBusinessMetric(span, "user.session.duration", 300)
	gootel.RecordBusinessMetric(span, "user.api.calls", 15)

	gootel.AddBusinessEvent(span, "user.profile.updated", map[string]interface{}{
		"user_id":  "USER-123",
		"fields":   "name,email,phone",
		"source":   "web_ui",
		"duration": "250ms",
	})

	gootel.AddBusinessContext(span, "user", map[string]string{
		"id":       "USER-123",
		"username": "john.doe",
		"email":    "john.doe@example.com",
		"tier":     "premium",
	})

	gootel.AddBusinessContext(span, "session", map[string]string{
		"id":        "SESSION-456",
		"ip_address": "192.168.1.100",
		"user_agent": "Mozilla/5.0",
		"device":     "desktop",
	})
}
