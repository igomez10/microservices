package tracerhelper

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func GetTracerWithAppName(appname string) trace.Tracer {
	tracerProvider := otel.GetTracerProvider()
	tracer := tracerProvider.Tracer(appname)
	return tracer
}
