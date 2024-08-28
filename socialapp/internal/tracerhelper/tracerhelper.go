package tracerhelper

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var appname = "socialapp"

func GetTracer() trace.Tracer {
	tracerProvider := otel.GetTracerProvider()
	tracer := tracerProvider.Tracer(appname)

	return tracer
}
