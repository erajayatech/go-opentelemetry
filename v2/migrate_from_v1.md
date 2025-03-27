# go-opentelemetry v2

Go OpenTelemetry Helper.

This document is for you who use go-opentelemetry v1.

## Migrate from v1

V2 maintained compatibility with version 1.

You can see how compatible `v2` with `v1` is in full example, [example_compatible_v1 server](./example_compatible_v1/server/main.go) and [example_compatible_v1 client](./example_compatible_v1/client/main.go).

## Installation v2

For you who already use v1, you need to search and replace your import, from

```bash
"github.com/erajayatech/go-opentelemetry"
```

to

```bash
"github.com/erajayatech/go-opentelemetry/v2"
```

then

```bash
go get github.com/erajayatech/go-opentelemetry/v2@edge
```

then

```bash
go mod tidy
```

## Things should be highlighted in `v2`

See [highlighted_in_v2.md](./highlighted_in_v2.md)
