package main

import (
	"net"

	"github.com/erajayatech/go-opentelemetry/v2/example_compatible_v1/pbfoo"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func runGRPCServer() {
	grpcServer := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
	pbfoo.RegisterExampleServer(grpcServer, &GRPCExampleServer{})
	lis, err := net.Listen("tcp", "localhost:4001")
	fatalIfErr(err)
	err = grpcServer.Serve(lis)
	fatalIfErr(err)
}
