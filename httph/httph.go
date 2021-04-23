package httph

import (
	"context"
	"net/http"

	"github.com/opentracing/opentracing-go/ext"

	"github.com/opentracing/opentracing-go"
)

func TracingHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var span opentracing.Span
		tracer := opentracing.GlobalTracer()
		spCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		if err != nil {
			span = opentracing.StartSpan(
				r.URL.Path,
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
				opentracing.Tag{Key: "http.url", Value: r.URL},
				opentracing.Tag{Key: "http.form", Value: r.Form},
				opentracing.Tag{Key: "http.method", Value: r.Method},
				ext.SpanKindRPCServer,
			)
			defer span.Finish()
		} else {
			span = opentracing.StartSpan(
				r.URL.Path,
				opentracing.ChildOf(spCtx),
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
				opentracing.Tag{Key: "http.url", Value: r.URL},
				opentracing.Tag{Key: "http.form", Value: r.Form},
				opentracing.Tag{Key: "http.method", Value: r.Method},
				ext.SpanKindRPCServer,
			)
			defer span.Finish()
		}
		ctx := opentracing.ContextWithSpan(context.TODO(), span)
		next(w, r.WithContext(ctx))
	}
}
