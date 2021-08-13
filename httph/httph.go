package httph

import (
	"net/http"

	"github.com/opentracing/opentracing-go/ext"

	"github.com/opentracing/opentracing-go"
)

func TracingHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !opentracing.IsGlobalTracerRegistered() {
			next(w, r)
			return
		}
		var span opentracing.Span
		tracer := opentracing.GlobalTracer()
		spCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		if err != nil {
			span = tracer.StartSpan(
				r.URL.Path,
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
				ext.SpanKindRPCServer,
			)
		} else {
			span = tracer.StartSpan(
				r.URL.Path,
				opentracing.ChildOf(spCtx),
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
				ext.SpanKindRPCServer,
			)
		}
		// 设置http.url, http.method
		ext.HTTPMethod.Set(span, r.Method)
		ext.HTTPUrl.Set(span, r.URL.String())
		defer span.SetTag("http.form", r.PostForm.Encode()).Finish()
		next(w, r.WithContext(opentracing.ContextWithSpan(r.Context(), span)))
	}
}
