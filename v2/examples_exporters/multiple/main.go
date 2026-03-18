package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	gootel "github.com/erajayatech/go-opentelemetry/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func main() {
	ctx := context.Background()

	config := gootel.ExporterConfig{
		ServiceName:    "multiple-exporters-service",
		ServiceVersion: "1.0.0",
		Environment:    "production",
	}

	exporters, err := createExporters(ctx, config)
	if err != nil {
		log.Fatal(err)
	}
	defer shutdownExporters(ctx, exporters)

	tp, err := gootel.NewTraceProviderWithMultipleExporters(ctx, config, exporters)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down trace provider: %v", err)
		}
	}()

	otel.SetTracerProvider(tp)

	log.Printf("Starting multiple exporters example with %d exporters", len(exporters))
	runExampleOperations(ctx)
}

func createExporters(ctx context.Context, config gootel.ExporterConfig) ([]trace.SpanExporter, error) {
	var exporters []trace.SpanExporter

	if os.Getenv("ENABLE_NEW_RELIC") == "true" {
		nrExporter, err := createNewRelicExporter(ctx, config)
		if err != nil {
			log.Printf("Warning: Failed to create New Relic exporter: %v", err)
		} else {
			exporters = append(exporters, nrExporter)
			log.Println("New Relic exporter enabled")
		}
	}

	if os.Getenv("ENABLE_DATADOG") == "true" {
		ddExporter, err := createDatadogExporter(ctx, config)
		if err != nil {
			log.Printf("Warning: Failed to create Datadog exporter: %v", err)
		} else {
			exporters = append(exporters, ddExporter)
			log.Println("Datadog exporter enabled")
		}
	}

	if os.Getenv("ENABLE_JAEGER") == "true" {
		jaegerExporter, err := createJaegerExporter(ctx, config)
		if err != nil {
			log.Printf("Warning: Failed to create Jaeger exporter: %v", err)
		} else {
			exporters = append(exporters, jaegerExporter)
			log.Println("Jaeger exporter enabled")
		}
	}

	if os.Getenv("ENABLE_OTLP") == "true" {
		otlpExporter, err := createOTLPExporter(ctx, config)
		if err != nil {
			log.Printf("Warning: Failed to create OTLP exporter: %v", err)
		} else {
			exporters = append(exporters, otlpExporter)
			log.Println("OTLP exporter enabled")
		}
	}

	if os.Getenv("ENABLE_STDOUT") == "true" || len(exporters) == 0 {
		stdoutExporter, err := createStdoutExporter(config)
		if err != nil {
			log.Printf("Warning: Failed to create stdout exporter: %v", err)
		} else {
			exporters = append(exporters, stdoutExporter)
			log.Println("Stdout exporter enabled")
		}
	}

	if len(exporters) == 0 {
		return nil, fmt.Errorf("no exporters could be created")
	}

	return exporters, nil
}

func createNewRelicExporter(ctx context.Context, config gootel.ExporterConfig) (trace.SpanExporter, error) {
	apiKey := os.Getenv("NEW_RELIC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("NEW_RELIC_API_KEY is required")
	}

	endpoint := os.Getenv("NEW_RELIC_ENDPOINT")
	if endpoint == "" {
		endpoint = "otlp.nr-data.net:4317"
	}

	return otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithHeaders(map[string]string{"api-key": apiKey}),
		otlptracegrpc.WithCompressor("gzip"),
	)
}

func createDatadogExporter(ctx context.Context, config gootel.ExporterConfig) (trace.SpanExporter, error) {
	apiKey := os.Getenv("DATADOG_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("DATADOG_API_KEY is required")
	}

	endpoint := os.Getenv("DATADOG_ENDPOINT")
	if endpoint == "" {
		endpoint = "trace-agent.datadoghq.com:4317"
	}

	useHTTP := os.Getenv("DATADOG_USE_HTTP") == "true"

	if useHTTP {
		return otlptracehttp.New(ctx,
			otlptracehttp.WithEndpoint(endpoint),
			otlptracehttp.WithHeaders(map[string]string{"DD-API-KEY": apiKey}),
			otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
			otlptracehttp.WithInsecure(),
		)
	}

	return otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithHeaders(map[string]string{"DD-API-KEY": apiKey}),
		otlptracegrpc.WithCompressor("gzip"),
		otlptracegrpc.WithInsecure(),
	)
}

func createJaegerExporter(ctx context.Context, config gootel.ExporterConfig) (trace.SpanExporter, error) {
	endpoint := os.Getenv("JAEGER_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:4317"
	}

	return otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
	)
}

func createOTLPExporter(ctx context.Context, config gootel.ExporterConfig) (trace.SpanExporter, error) {
	endpoint := os.Getenv("OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:4317"
	}

	useHTTP := os.Getenv("OTLP_USE_HTTP") == "true"

	headers := map[string]string{}
	if headerStr := os.Getenv("OTLP_HEADERS"); headerStr != "" {
		headers = parseHeaders(headerStr)
	}

	if useHTTP {
		return otlptracehttp.New(ctx,
			otlptracehttp.WithEndpoint(endpoint),
			otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithHeaders(headers),
		)
	}

	return otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithCompressor("gzip"),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithHeaders(headers),
	)
}

func createStdoutExporter(config gootel.ExporterConfig) (trace.SpanExporter, error) {
	return tracetest.NewInMemoryExporter(), nil
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

func shutdownExporters(ctx context.Context, exporters []trace.SpanExporter) {
	for _, exporter := range exporters {
		if err := exporter.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down exporter: %v", err)
		}
	}
}

func runExampleOperations(ctx context.Context) {
	ctx, span := gootel.RecordSpan(ctx)
	defer span.End()

	gootel.AddBusinessAttribute(span, "operation.type", "multi-exporter-demo")
	gootel.AddBusinessAttribute(span, "environment", "production")

	processHTTPRequest(ctx)
	processDatabaseOperation(ctx)
	processCacheOperation(ctx)
	processBusinessLogic(ctx)
}

func processHTTPRequest(ctx context.Context) {
	ctx, span := gootel.TraceHTTPRequest(ctx, "multiple-exporters-service", "GET", "https://api.example.com/users")
	defer span.End()

	gootel.AddBusinessAttribute(span, "http.category", "external-api")
	gootel.AddBusinessAttribute(span, "api.resource", "users")

	startTime := time.Now()
	time.Sleep(50 * time.Millisecond)
	duration := time.Since(startTime)

	gootel.RecordHTTPSuccess(span, 200, duration)
	gootel.AddEventToSpan(ctx, "http.request.completed", map[string]interface{}{
		"response_size": 1024,
		"cache_hit":     true,
		"exporters":     "all",
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

	gootel.AddBusinessAttribute(span, "business.process", "multi-exporter.workflow")
	gootel.AddBusinessAttribute(span, "workflow.id", "WORKFLOW-123")

	gootel.RecordBusinessMetric(span, "workflow.duration", 250)
	gootel.RecordBusinessMetric(span, "workflow.steps", 4)

	gootel.AddBusinessEvent(span, "workflow.completed", map[string]interface{}{
		"workflow_id": "WORKFLOW-123",
		"steps_count": 4,
		"duration_ms": 250,
		"exporters":   "all",
		"success":     true,
	})

	gootel.AddBusinessContext(span, "workflow", map[string]string{
		"id":     "WORKFLOW-123",
		"type":   "multi-step",
		"status": "completed",
		"source": "api",
	})
}
