# go-opentelemetry

## Getting Started

You can find a getting started guide on [opentelemetry.io](https://opentelemetry.io/docs/go/getting-started/).

OpenTelemetry's goal is to provide a single set of APIs to capture distributed
traces and metrics from your application and send them to an observability
platform. This project allows you to do just that for applications written in
Go. There are two steps to this process: instrument your application, and
configure an exporter.

## Dependency
* [Gin Web Framework](https://github.com/gin-gonic/gin)
* [Open Telemetry](https://pkg.go.dev/go.opentelemetry.io/otel)
* [GoDotEnv by joho](https://github.com/joho/godotenv)

## Install
Go Version 1.16+
```
go get github.com/erajayatech/go-opentelemetry
```

## Setup Environment
- Set the following environment variables:

* `MODE=<your_application_mode>`
  * Example : `prod`
* `APP_VERSION=<your_application_version>`
* `APP_NAME=<your_application_name>`
* `OTEL_EXPORTER_OTLP_ENDPOINT=https://otlp.nr-data.net:4317`
* `OTEL_EXPORTER_OTLP_HEADERS="api-key=<your_license_key>"`
  * Replace `<your_license_key>` with your
    [Account License Key](https://one.newrelic.com/launcher/api-keys-ui.launcher).
* `OTEL_SAMPLED=true`
  * Be careful about using this sampler in a production application with significant traffic: a new trace will be started and exported for every request. If you won't sampling for every request just set to `false`

## How To Use
- Import package on main.go
```go
otel "github.com/erajayatech/go-opentelemetry"
```

- Add below code to `main.go`
```go
otelTracerService := otel.ConstructOtelTracer()
otelTracerServiceErr := otelTracerService.SetTraceProviderNewRelic(ctx)
if otelTracerServiceErr != nil {
    panic(otelTracerServiceErr)
}
```

- Tracing `Controller`
    * Import package
    ```go
    otel "github.com/erajayatech/go-opentelemetry"
    ```
    * Start Span Tracer
    ```go
    ctx, span := otel.Start(context)
    defer span.End()
    ```
    * Use code below to add span's tags
    ```go
    otel.AddSpanTags(span, map[string]string{"traceId": singleton.Trace().ID, "RegisterRequest": string(requestMarshalled)})
    ```
    * Use code below to add span error 
    ```go
    if err != nil {
        otel.AddSpanError(span, err)
        otel.FailSpan(span, "Request not valid")
        httpresponse.BadRequest(context, "Request not valid")
        return
    }
    ```

- Tracing `Service`
    * Import package
    ```go
    otel "github.com/erajayatech/go-opentelemetry"
    ```
    * Add context argument to function that we want to trace, example :
    ```go
    func (service *RegisterService) Register(context context.Context, request RegisterRequest, platform string) *httpresponse.HTTPError {}
    ```
    * Use code below for add span
    ```go
    ctx, span := otel.NewSpan(context, "authregister.service.Register", "")
    defer span.End()
    ```
    * Use code below to add span error 
    ```go
    if registeredPlatforms.Id == 0 {
        httpError.Code = http.StatusBadRequest
        httpError.Message = "Platform not valid"
        otel.AddSpanError(span, valErr)
        otel.FailSpan(span, httpError.Message)
        return &httpError
    }
    ```

- Tracing `Query`
  - Add code below on `database.go`
    * Import package
    ```go
    "gorm.io/gorm/logger"
    "gorm.io/plugin/opentelemetry/logging/logrus"
    "gorm.io/plugin/opentelemetry/tracing"
    ```
    * Set `gorm logger`
    ```go
    logger := logger.New(
        logrus.NewWriter(),
        logger.Config{
            SlowThreshold: time.Millisecond,
            LogLevel:      logger.Warn,
            Colorful:      false,
        },
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger})

    if err != nil {
        log.Println("Connected to database Failed:", err)
    }

    if err := db.Use(tracing.NewPlugin()); err != nil {
        panic(err)
    }
    ```
  - Add code below on `repository`
    * Import package
    ```go
    otel "github.com/erajayatech/go-opentelemetry"
    ```
    * Add context argument to function that we want to trace, example :
    ```go
    func (repo *RegisterRepository) getCustomerByEmail(context context.Context, customer *model.Customer, email string) {}
    ```
    * Use code below for add span
    ```go
    ctx, span := otel.NewSpan(context, "authregister.RegisterRepository.getCustomerByEmail", "")
	defer span.End()
    ```
    * Use code below for trace a query
    ```go
    repo.db.WithContext(ctx).First(customer, "email = ?", strings.ToLower(email))
    ```

- Tracing `External Service`
    * Import package
    ```go
    otel "github.com/erajayatech/go-opentelemetry"
    ```
    * Add context argument to function that we want to trace, example :
    ```go
    func (capillary *Capillary) sendToStagingCustomer(context context.Context, url string, payload interface{}, stagingResponse *StagingCustomerResponse) (*StagingCustomerResponse, error) {}
    ```
    * Use code below for add span
    ```go
    httpSpanAttribute := otel.HttpSpanAttribute{}
    httpSpanAttribute.Method = request.Method
    httpSpanAttribute.Url = request.RequestURI
    httpSpanAttribute.IP = request.Host

    _, span := otel.NewHttpSpan(context, "capillary.sendToStagingCustomer", url, httpSpanAttribute)
    defer span.End()
    ```
    * Use code below for add event span
    ```go
    otel.AddSpanEvents(span, "capillary.sendToStagingCustomer", map[string]string{"traceId": capillary.traceID, "StagCust Endpoint": url, "StagCust Request": string(payloadByte), "StagCustResponse": compactedBuffer.String(), "http.status_code": response.Status})
    ```


