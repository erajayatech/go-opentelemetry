# Configuring Datadog Exporter

In version 2, we've added support for exporting traces to Datadog alongside New Relic and Jaeger. You can configure multiple exporters using environment variables.

## Environment Variables

To enable Datadog exporter, set the following environment variables:

```bash
# Enable Datadog exporter (set to true to enable)
ENABLE_DATADOG_EXPORTER=true

# Datadog endpoint (OTLP/gRPC endpoint)
DATADOG_ENDPOINT=localhost:4317

# Datadog API Key
DATADOG_API_KEY=your_datadog_api_key_here

# Optional: Disable New Relic if you only want to use Datadog
ENABLE_NEWRELIC_EXPORTER=false
```

## Datadog Agent Setup

1. Install and configure Datadog Agent with OTLP support:

```bash
# Using Docker
docker run -d --name datadog-agent \
  -e DD_API_KEY=your_datadog_api_key_here \
  -e DD_OTLP_CONFIG_RECEIVER_PROTOCOLS_GRPC_ENDPOINT=0.0.0.0:4317 \
  -e DD_OTLP_CONFIG_RECEIVER_PROTOCOLS_HTTP_ENDPOINT=0.0.0.0:4318 \
  -p 4317:4317 \
  -p 4318:4318 \
  -p 8126:8126 \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  -v /proc/:/host/proc/:ro \
  -v /sys/fs/cgroup/:/host/sys/fs/cgroup:ro \
  datadog/agent:latest
```

2. Or using Kubernetes with Helm:

```bash
helm repo add datadog https://helm.datadoghq.com
helm repo update

helm install datadog-agent datadog/datadog \
  --set datadog.apiKey=your_datadog_api_key_here \
  --set datadog.otlp.receiver.protocols.grpc.enabled=true \
  --set datadog.otlp.receiver.protocols.http.enabled=true
```

## Implementation Example

No code changes are needed other than setting the environment variables above. The trace provider will automatically detect and configure the Datadog exporter if enabled.

Here's an example of how to set up your application with Datadog:

```go
package main

import (
	"context"
	"log"
	
	gootel "github.com/erajayatech/go-opentelemetry/v2"
)

func main() {
	// Set up the trace provider
	tp, err := gootel.NewTraceProvider(context.Background())
	if err != nil {
		log.Fatalf("Error creating trace provider: %v", err)
	}
	defer tp.Shutdown(context.Background())
	
	// Your application code here...
}
```

## Multiple Exporters

You can enable multiple exporters simultaneously:

```bash
# Enable all three exporters
ENABLE_NEWRELIC_EXPORTER=true
ENABLE_JAEGER_EXPORTER=true
ENABLE_DATADOG_EXPORTER=true

# Configure each exporter
OTEL_EXPORTER_OTLP_ENDPOINT=https://otlp.nr-data.net:4317
OTEL_EXPORTER_OTLP_HEADERS="api-key=your_newrelic_license_key"
JAEGER_ENDPOINT=localhost:4317
DATADOG_ENDPOINT=localhost:4317
DATADOG_API_KEY=your_datadog_api_key_here
```

## Visualizing Traces in Datadog

Once configured, your traces will appear in the Datadog APM interface at https://app.datadoghq.com/apm/traces. You can:

- View distributed traces across your services
- Monitor application performance metrics
- Set up alerts based on trace data
- Correlate traces with logs and metrics

The traces will include service names, operation names, and custom tags configured in your application.