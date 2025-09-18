# Configuring Jaeger Exporter

In version 2, we've added support for exporting traces to both New Relic and Jaeger. You can configure either or both exporters using environment variables.

## Environment Variables

To enable Jaeger exporter, set the following environment variables:

```bash
# Enable Jaeger exporter (set to true to enable)
ENABLE_JAEGER_EXPORTER=true

# Jaeger endpoint (OTLP/gRPC endpoint)
JAEGER_ENDPOINT=localhost:4317

# Optional: Disable New Relic if you only want to use Jaeger
ENABLE_NEWRELIC_EXPORTER=false
```

## Jaeger Setup

1. Start Jaeger using Docker:

```bash
docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
  -e COLLECTOR_OTLP_ENABLED=true \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 4317:4317 \
  -p 4318:4318 \
  -p 14250:14250 \
  -p 14268:14268 \
  -p 14269:14269 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.45
```

2. Access the Jaeger UI at http://localhost:16686

## Implementation Example

No code changes are needed other than setting the environment variables above. The trace provider will automatically detect and configure the Jaeger exporter if enabled.

Here's an example of how to set up your application with both New Relic and Jaeger:

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

## Visualizing Traces

When you have both exporters enabled, your traces will be sent to both New Relic and Jaeger. This provides flexibility for local debugging with Jaeger while still sending production traces to New Relic.

![jaeger-ui-example](./jaeger_ui.png)