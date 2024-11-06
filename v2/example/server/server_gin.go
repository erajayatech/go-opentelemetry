package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func runHTTPServerGin() {
	ginEngine := gin.Default()
	ginEngine.Use(
		otelgin.Middleware(""), // use otelgin to instrument http request
	)
	ginEngine.GET("foo", controllerGinFoo)
	httpServer := &http.Server{Addr: "localhost:4000", Handler: ginEngine}
	err := httpServer.ListenAndServe()
	fatalIfErr(err)
}
