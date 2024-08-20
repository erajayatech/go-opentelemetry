package main

import (
	"context"

	gootel "github.com/erajayatech/go-opentelemetry/v2"
)

func repoGetFoo(ctx context.Context) {
	ctx, span := gootel.RecordSpan(ctx) //nolint:staticcheck,ineffassign
	defer span.End()
}
