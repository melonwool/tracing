package tracing

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

// SpanContextFromContext 从context 获取 span 和 context
func SpanContextFromContext(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
	tracer := opentracing.GlobalTracer()
	span, newCtx := opentracing.StartSpanFromContextWithTracer(ctx, tracer, operationName)
	return span, newCtx
}
