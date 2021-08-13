package ginh

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// OpenTracingHandler open tracing gin handler
func OpenTracingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !opentracing.IsGlobalTracerRegistered() {
			c.Next()
			return
		}
		var span opentracing.Span
		tracer := opentracing.GlobalTracer()
		spCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil {
			span = tracer.StartSpan(
				c.Request.URL.Path,
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
				ext.SpanKindRPCServer,
			)
		} else {
			span = tracer.StartSpan(
				c.Request.URL.Path,
				opentracing.ChildOf(spCtx),
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
				ext.SpanKindRPCServer,
			)
		}
		// 设置http.url, http.method
		ext.HTTPMethod.Set(span, c.Request.Method)
		ext.HTTPUrl.Set(span, c.Request.URL.String())

		defer span.Finish()
		// request context 中添加span 信息 用来传递
		c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), span))
		_ = tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		c.Next()
	}
}
