package ginh

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// OpenTracingHandler open tracing gin handler
func OpenTracingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var span opentracing.Span
		tracer := opentracing.GlobalTracer()
		spCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil {
			span = opentracing.StartSpan(
				c.Request.URL.Path,
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
				opentracing.Tag{Key: "http.url", Value: c.Request.URL},
				opentracing.Tag{Key: "http.form", Value: c.Request.Form},
				opentracing.Tag{Key: "http.method", Value: c.Request.Method},
				ext.SpanKindRPCServer,
			)
		} else {
			span = opentracing.StartSpan(
				c.Request.URL.Path,
				opentracing.ChildOf(spCtx),
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
				opentracing.Tag{Key: "http.url", Value: c.Request.URL},
				opentracing.Tag{Key: "http.method", Value: c.Request.Method},
				ext.SpanKindRPCServer,
			)
		}
		defer span.Finish()
		c.Set("spanCtx", span.Context())
		c.Next()
	}
}

// GinSpanContext 获取span 和 context
func GinSpanContext(c *gin.Context, operationName string) (opentracing.Span, context.Context) {
	var span opentracing.Span
	var ctx context.Context
	tracer := opentracing.GlobalTracer()
	if spanCtxInf, exist := c.Get("spanCtx"); exist {
		if spanCtx, ok := spanCtxInf.(opentracing.SpanContext); ok {
			span, ctx = opentracing.StartSpanFromContextWithTracer(context.TODO(), tracer, operationName, opentracing.ChildOf(spanCtx))
		}
	}
	if span == nil {
		span, ctx = opentracing.StartSpanFromContextWithTracer(context.TODO(), tracer, operationName)
	}
	return span, ctx
}