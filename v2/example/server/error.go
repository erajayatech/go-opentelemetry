package main

import (
	"log/slog"
	"os"
)

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
