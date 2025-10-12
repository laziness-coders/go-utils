package logger

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// ExtractTraceFields extracts tracing information from context for both Datadog and OpenTelemetry.
func ExtractTraceFields(ctx context.Context) []zap.Field {
	if ctx == nil {
		return nil
	}

	var fields []zap.Field

	// OpenTelemetry tracing
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		spanCtx := span.SpanContext()
		fields = append(fields,
			zap.String("trace_id", spanCtx.TraceID().String()),
			zap.String("span_id", spanCtx.SpanID().String()),
		)
	}

	return fields
}
