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
		ServiceName:    "otlp-example-service",
		ServiceVersion: "1.0.0",
		Environment:    "production",
	}

	otlpEndpoint := os.Getenv("OTLP_ENDPOINT")
	if otlpEndpoint == "" {
		otlpEndpoint = "localhost:4317"
	}

	otlpHeaders := os.Getenv("OTLP_HEADERS")
	var headers map[string]string
	if otlpHeaders != "" {
		headers = parseHeaders(otlpHeaders)
	}

	useHTTP := os.Getenv("OTLP_USE_HTTP") == "true"

	tp, err := gootel.NewTraceProviderWithOTLP(
		ctx,
		config,
		otlpEndpoint,
		useHTTP,
		headers,
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

func parseHeaders(headerString string) map[string]string {
	headers := make(map[string]string)
	pairs := splitByComma(headerString)
	for _, pair := range pairs {
		kv := splitByEquals(pair)
		if len(kv) == 2 {
			headers[kv[0]] = kv[1]
		}
	}
	return headers
}

func splitByComma(s string) []string {
	var result []string
	start := 0
	for i, c := range s {
		if c == ',' {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	result = append(result, s[start:])
	return result
}

func splitByEquals(s string) []string {
	for i, c := range s {
		if c == '=' {
			return []string{s[:i], s[i+1:]}
		}
	}
	return []string{s}
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
	ctx, span := gootel.TraceHTTPRequest(ctx, "otlp-example-service", "GET", "https://api.example.com/products")
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
	})
}

func processDatabaseOperation(ctx context.Context) {
	ctx, span := gootel.TraceDBOperation(ctx, "postgresql", "appdb", "SELECT", "SELECT * FROM products WHERE category = $1")
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
