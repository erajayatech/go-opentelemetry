package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetString(k string) (string, error) {
	v, ok := os.LookupEnv(k)
	if !ok {
		return "", fmt.Errorf("'%s' is not found in os env var, please set it", k)
	}

	return v, nil
}

func GetInt(k string) (int, error) {
	vStr, err := GetString(k)
	if err != nil {
		return 0, fmt.Errorf("error get string:: %w", err)
	}

	v, err := strconv.Atoi(vStr)
	if err != nil {
		return 0, fmt.Errorf("error convert string '%s' to int:: %w", vStr, err)
	}

	return v, nil
}

func GetBool(k string) (bool, error) {
	vStr, err := GetString(k)
	if err != nil {
		return false, fmt.Errorf("error get string:: %w", err)
	}

	v, err := strconv.ParseBool(vStr)
	if err != nil {
		return false, fmt.Errorf("error convert string '%s' to bool:: %w", vStr, err)
	}

	return v, nil
}

func GetServiceName() (string, error) {
	return GetString("APP_NAME")
}

func GetAppVersion() (string, error) {
	return GetString("APP_VERSION")
}

func GetAppEnvironment() (string, error) {
	return GetString("MODE")
}

func GetOtelOTLPNewrelicHost() (string, error) {
	v, err := GetString("OTEL_EXPORTER_OTLP_ENDPOINT")
	if err != nil {
		return "", err
	}
	v = strings.ReplaceAll(v, "https://", "")
	return v, nil
}

func GetOtelOTLPNewrelicHeaderAPIKey() (string, error) {
	v, err := GetString("OTEL_EXPORTER_OTLP_HEADERS")
	if err != nil {
		return "", err
	}
	v = strings.ReplaceAll(v, "api-key=", "")
	return v, nil
}

// GetJaegerEndpoint returns the Jaeger endpoint URL
func GetJaegerEndpoint() (string, error) {
	return GetString("JAEGER_ENDPOINT")
}

// IsJaegerEnabled returns whether the Jaeger exporter is enabled
func IsJaegerEnabled() bool {
	v, err := GetString("ENABLE_JAEGER_EXPORTER")
	if err != nil {
		return false
	}
	enabled, err := strconv.ParseBool(v)
	if err != nil {
		return false
	}
	return enabled
}

// IsNewRelicEnabled returns whether the New Relic exporter is enabled
func IsNewRelicEnabled() bool {
	v, err := GetString("ENABLE_NEWRELIC_EXPORTER")
	if err != nil {
		// Default to true for backward compatibility
		return true
	}
	enabled, err := strconv.ParseBool(v)
	if err != nil {
		// Default to true for backward compatibility
		return true
	}
	return enabled
}
