package utils

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"runtime"
)

func StartSpan(ctx context.Context) (context.Context, trace.Span) {
	pc, _, _, _ := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	return otel.GetTracerProvider().Tracer("").Start(ctx, details.Name())
}
