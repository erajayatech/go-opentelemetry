package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	gootel "github.com/erajayatech/go-opentelemetry/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	fatalIfErr(err)

	tp, err := gootel.NewTraceProvider(context.Background())
	fatalIfErr(err)
	defer func() {
		err := tp.Shutdown(context.Background())
		warnIfErr(err)
	}()

	go runHTTPServerGin()
	go runGRPCServer()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	slog.Info("listens for the interrupt or terminate signal from the OS")
	<-ctx.Done()
	stop()
}
