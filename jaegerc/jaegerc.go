package jaegerc

import (
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

type JaegerType string

const (
	ConstType         JaegerType = "const"
	ProbabilisticType JaegerType = "probabilistic"
	RateLimitingType  JaegerType = "rateLimiting"
	RemoteType        JaegerType = "remote"
)

// NewJaegerTracer 初始化jaeger client
// jaegerType 默认使用 const  param 1 全采样 0 不采样 0.5 50%采样
func NewJaegerTracer(jaegerHostPort, appName string, jaegerType JaegerType, param float64) (opentracing.Tracer, io.Closer) {
	cfg := &jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  string(jaegerType), //固定采样
			Param: param,              //1=全采样、0=不采样
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: jaegerHostPort,
		},

		ServiceName: appName,
	}

	// tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.NullLogger))
	if err != nil {
		fmt.Printf("ERROR: cannot init Jaeger: %v\n", err)
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}
