# go-opentelemetry v2

Go OpenTelemetry Helper.

## Things should be highlighted in `v2`

1. Many func from `v1` is retained in `v2` but we mark it as deprecated.
2. Main feature of `v2` is otel context propagation, so you need to setup server and client to handle it. 

- Server gin: Use `otelgin`, see [example](./example/server/server_gin.go), 
- Server grpc: Use `otelgrpc`, see [example](./example/server/server_grpc.go), 
- Client http: Use `otelhttp`, see [example](./example/client/http.go),
- Client grpc: Use `otelgrpc`, see [example](./example/client/grpc.go),

3. We recommend to import the library like this

```go
gootel "github.com/erajayatech/go-opentelemetry/v2"
```

instead of

```go
otel "github.com/erajayatech/go-opentelemetry/v2"
```

because one day you will see yourself need to import

```go
"go.opentelemetry.io/otel"
```
