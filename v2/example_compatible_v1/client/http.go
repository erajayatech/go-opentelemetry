package main

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	gootel "github.com/erajayatech/go-opentelemetry/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func extapiHTTPFoo() {
	ctx, span := gootel.NewSpan(context.Background(), "extapiHTTPFoo", "") //nolint:staticcheck
	defer span.End()

	client := &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:4000/foo", nil)
	fatalIfErr(err)
	res, err := client.Do(req)
	fatalIfErr(err)
	resBody, err := io.ReadAll(res.Body)
	fatalIfErr(err)
	body := map[string]any{}
	err = json.Unmarshal(resBody, &body)
	fatalIfErr(err)
	slog.Info("success", "body", body)
}
