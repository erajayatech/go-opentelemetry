package main

import (
	"context"
	"log/slog"

	gootel "github.com/erajayatech/go-opentelemetry/v2"
	"github.com/erajayatech/go-opentelemetry/v2/example_compatible_v1/pbfoo"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func extapiGRPCFoo() {
	ctx, span := gootel.NewSpan(context.Background(), "extapiGRPCFoo", "") //nolint:staticcheck
	defer span.End()

	conn, err := grpc.NewClient("localhost:4001", grpc.WithStatsHandler(otelgrpc.NewClientHandler()), grpc.WithTransportCredentials(insecure.NewCredentials()))
	fatalIfErr(err)
	client := pbfoo.NewExampleClient(conn)
	res, err := client.Foo(ctx, &pbfoo.ReqFoo{})
	fatalIfErr(err)
	slog.Info("success", "res", res)
}
