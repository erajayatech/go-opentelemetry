package goopentelemetry

func EnvironmentMode() string {
	return GetEnv("MODE")
}

func AppName() string {
	return GetEnv("APP_NAME")
}

func AppVersion() string {
	return GetEnv("APP_VERSION")
}

func OtelSampled() bool {
	return AnyToBool(GetEnv("OTEL_SAMPLED"))
}

func OtelJaegerURL() string {
	return GetEnv("OTEL_JAEGER_URL")
}
