package gootel

import (
	"fmt"
	"os"

	"github.com/erajayatech/go-opentelemetry/v2/internal/config"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

func getResource() (*resource.Resource, error) {
	fail := func(err error, msg string) (*resource.Resource, error) {
		return nil, fmt.Errorf("%s:: %w", msg, err)
	}

	serviceName, err := config.GetServiceName()
	if err != nil {
		return fail(err, "error get service name")
	}

	appVersion, err := config.GetAppVersion()
	if err != nil {
		return fail(err, "error get app version")
	}

	appEnv, err := config.GetAppEnvironment()
	if err != nil {
		return fail(err, "error get app environment")
	}

	_resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(appVersion),
		semconv.DeploymentEnvironmentKey.String(appEnv),
		semconv.K8SPodNameKey.String(os.Getenv("HOSTNAME")), // k8 will set this unique for every pod.
	)

	return _resource, nil
}
