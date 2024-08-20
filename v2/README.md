# go-opentelemetry v2

Go OpenTelemetry Helper.

Why we need `v2`?

1. Span trace front to back (context propagation).
2. Upgrade go version to `v1.21.0` and otel version from `v1.10.0` to `v1.28.0`, see [why_need_upgrade_version](./why_need_upgrade_version.md).
3. Better library API. See [better_api.md](./better_api.md)

## Feature

- [x] Opentelemetry Trace
- [x] Opentelemetry Context Propagation

![context_propagation](./README_asset/context_propagation.png)


## Installation v2

```bash
go get github.com/erajayatech/go-opentelemetry/v2@v2.0.0-alpha.6
```

```go
import gootel "github.com/erajayatech/go-opentelemetry/v2"
```

## Usage

See [example server](./example/server/main.go) and [example client](./example/client/main.go).

In New Relic you will get.

![grpc-client-span](./README_asset/grpc_span.png)

![http-client-span](./README_asset/http_span.png)

## Migrate from v1

See [Migrate from v1](./migrate_from_v1.md)

## Things should be highlighted in `v2`

See [highlighted_in_v2.md](./highlighted_in_v2.md)
