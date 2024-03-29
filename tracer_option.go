package goopentelemetry

type OtelTracerOptionFunc func(*otelTracer)

func SetEnv(env string) OtelTracerOptionFunc {
	return func(ot *otelTracer) {
		ot.env = env
	}
}

func SetVersion(version string) OtelTracerOptionFunc {
	return func(ot *otelTracer) {
		ot.version = version
	}

}

func SetAppName(name string) OtelTracerOptionFunc {
	return func(ot *otelTracer) {
		ot.service = name
	}
}

func IsSampledEnable(isEnabled bool) OtelTracerOptionFunc {
	return func(ot *otelTracer) {
		ot.sampled = isEnabled
	}
}

func SetExporterEndpoint(endpoint string) OtelTracerOptionFunc {
	return func(ot *otelTracer) {
		ot.exporterEndpoint = &endpoint
	}
}

func SetApiKey(key string) OtelTracerOptionFunc {
	return func(ot *otelTracer) {
		ot.apiKey = &key
	}
}
