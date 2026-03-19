package tracer

import (
	"context"

	"github.com/go-resty/resty/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func SetDataToTracer(span trace.Span, resp *resty.Response) {
	span.SetAttributes(
		attribute.Int("http.response.status_code", resp.StatusCode()), // Lưu status code
		attribute.String("http.response.body", string(resp.Body())),   // Lưu body response
	)
}

func NewSpan(ctx context.Context, traceInstance trace.Tracer, name string) (context.Context, trace.Span) {
	return traceInstance.Start(ctx, name, trace.WithSpanKind(trace.SpanKindClient))
}
