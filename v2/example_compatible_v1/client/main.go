package main

import (
	"context"
	"log/slog"
	"os"

	gootel "github.com/erajayatech/go-opentelemetry/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	fatalIfErr(err)

	tp := gootel.ConstructOtelTracer()                      //nolint:staticcheck
	err = tp.SetTraceProviderNewRelic(context.Background()) //nolint:staticcheck
	fatalIfErr(err)
	defer func() {
		err := tp.Shutdown(context.Background())
		warnIfErr(err)
	}()

	extapiHTTPFoo()
	extapiGRPCFoo()
}

func warnIfErr(err error) {
	if err != nil {
		slog.Warn(err.Error())
	}
}

func fatalIfErr(err error) {
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
