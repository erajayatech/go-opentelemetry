package main

import (
	"context"

	gootel "github.com/erajayatech/go-opentelemetry/v2"
)

func serviceFoo(ctx context.Context) {
	ctx, span := gootel.NewSpan(ctx, "serviceFoo", "") //nolint:staticcheck
	defer span.End()

	repoGetFoo(ctx)
}
