package tracerhelper

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func GetTracer() trace.Tracer {
	return otel.Tracer("tracerhelper")
}
