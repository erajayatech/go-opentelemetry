# OpenTelemetry Exporter Examples

This directory contains complete working examples for each OpenTelemetry exporter supported by the go-opentelemetry v2 library.

## Available Examples

- [New Relic](./newrelic) - Production-ready observability platform
- [Datadog](./datadog) - Cloud monitoring and security platform
- [Jaeger](./jaeger) - Distributed tracing system (local development)
- [Stdout](./stdout) - Console output for debugging
- [OTLP](./otlp) - Generic OpenTelemetry Protocol exporter
- [Multiple Exporters](./multiple) - Export to multiple destinations simultaneously

## Quick Start

Each example includes:
- Complete working code
- Environment configuration file (.env.example)
- Detailed usage instructions

### New Relic Example

```bash
cd newrelic
cp .env.example .env
# Edit .env with your New Relic API key
go run main.go
```

### Datadog Example

```bash
cd datadog
cp .env.example .env
# Edit .env with your Datadog API key
go run main.go
```

### Jaeger Example

```bash
# Start Jaeger locally with Docker
docker run -d --name jaeger \
  -e COLLECTOR_OTLP_ENABLED=true \
  -p 16686:16686 \
  -p 4317:4317 \
  -p 4318:4318 \
  jaegertracing/all-in-one:latest

cd jaeger
cp .env.example .env
go run main.go

# Open Jaeger UI at http://localhost:16686
```

### Stdout Example

```bash
cd stdout
go run main.go
```

### OTLP Example

```bash
cd otlp
cp .env.example .env
# Edit .env with your OTLP endpoint
go run main.go
```

### Multiple Exporters Example

```bash
cd multiple
cp .env.example .env
# Edit .env with your exporter configurations
go run main.go
```

## Features Demonstrated

Each example demonstrates:

1. **HTTP Tracing** - External API calls with proper attributes
2. **Database Tracing** - SQL operations with query events
3. **Redis Tracing** - Cache operations with hit/miss tracking
4. **Business Attributes** - Domain-specific context and metrics
5. **Error Handling** - Proper error recording and span status

## Best Practices

1. **Environment Variables** - Use .env files for sensitive configuration
2. **Graceful Shutdown** - Always shutdown trace providers properly
3. **Context Propagation** - Maintain context throughout the request
4. **Semantic Attributes** - Follow OpenTelemetry conventions
5. **Business Context** - Add domain-specific attributes for better observability

## Running All Examples

To run all examples:

```bash
for dir in */; do
  echo "Running example in $dir"
  cd "$dir"
  if [ -f ".env.example" ] && [ ! -f ".env" ]; then
    cp .env.example .env
    echo "Please configure .env in $dir before running"
  else
    go run main.go
  fi
  cd ..
done
```

## Troubleshooting

### Connection Issues

- Verify endpoints are accessible from your network
- Check firewall settings for port access
- Ensure API keys are valid and have proper permissions

### Jaeger Not Showing Traces

- Ensure Jaeger is running: `docker ps | grep jaeger`
- Check Jaeger logs: `docker logs jaeger`
- Verify endpoint matches your Jaeger configuration

### New Relic/Datadog Not Receiving Traces

- Verify API keys are correct
- Check account permissions
- Review service name matches your observability platform

## Additional Resources

- [Main Documentation](../README.md)
- [OpenTelemetry Specification](https://opentelemetry.io/docs/reference/specification/)
- [New Relic OTLP Setup](https://docs.newrelic.com/docs/more-integrations/open-source-telemetry-integrations/opentelemetry/opentelemetry-setup/)
- [Datadog APM](https://docs.datadoghq.com/tracing/)
- [Jaeger Documentation](https://www.jaegertracing.io/docs/)
