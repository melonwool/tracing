package ginh

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// OpenTracingHandler open tracing gin handler
func OpenTracingHandler(tracer opentracing.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var span opentracing.Span
		spCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
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
		c.Set("Tracer", tracer)
		c.Set("ParentSpanContext", span.Context())
		c.Next()
	}
}
