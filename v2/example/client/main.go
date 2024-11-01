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

	// create new trace provider
	tp, err := gootel.NewTraceProvider(context.Background())
	fatalIfErr(err)
	defer func() {
		err := tp.Shutdown(context.Background())
		warnIfErr(err)
	}()

	// create new meter provider
	mp, err := gootel.NewMeterProvider(context.Background())
	fatalIfErr(err)
	defer func() {
		err := mp.Shutdown(context.Background())
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
