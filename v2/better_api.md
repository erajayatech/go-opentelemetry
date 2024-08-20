# go-opentelemetry v2

Go OpenTelemetry Helper.

## Better API

### `NewSpan`

From:

```go
ctx, span := otel.NewSpan(ctx, helper.MyCaller(1), "")
```

to:

```go
ctx, span := otel.NewSpan(ctx)
```

No need to set span name with `helper.MyCaller(1)` and operations with `""`.

**Even better**, just use `RecordSpan` in everywhere.

```go
ctx, span := gootel.RecordSpan(ctx)
```

### `Start`

Previously func `Start` only applicable for `*gin.Context`.

```go
ctx, span := otel.Start(c) // c *gin.Context
```

You can not use this func if you use other library for controller, let say fiber, grpc, echo, mux.

In `v2` it's possible.

```go
ctx, span := gootel.Start(ctx) // ctx context.Context
```

**Even better**, just use `RecordSpan` in controller (just like in layer service and repository).

<hr/>

No need to think wheter to use `Start` or `NewSpan`, just always use `RecordSpan` in every layer.

```go
ctx, span := gootel.RecordSpan(c)   // c *gin.Context
ctx, span := gootel.RecordSpan(ctx) // ctx context.Context
```

