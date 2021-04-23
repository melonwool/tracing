package httpc

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"

	"github.com/getsentry/sentry-go"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

const ContentType = "Content-Type"

var (
	DefaultTimeout       = time.Second * 2
	DefaultRetryCount    = 1
	DefaultRetryWaitTime = time.Second * 1
)

type (
	Option struct {
		Timeout         time.Duration
		RetryWaitTime   time.Duration
		RetryCount      int
		ConnectionClose bool
	}
	OptFunc func(option *Request)
)

type Request struct {
	Option
}

func NewRequest(optFunc ...OptFunc) *Request {
	request := &Request{}
	for _, fn := range optFunc {
		fn(request)
	}
	return request
}

// Timeout 设置请求超时时间
func Timeout(duration time.Duration) OptFunc {
	return func(request *Request) {
		request.Timeout = duration
	}
}

// RetryWaitTime 重试等待时间
func RetryWaitTime(duration time.Duration) OptFunc {
	return func(request *Request) {
		request.RetryWaitTime = duration
	}
}

// RetryCount 重试次数
func RetryCount(count int) OptFunc {
	return func(request *Request) {
		request.RetryCount = count
	}
}

// ConnectionClose 是否关闭连接,长连/短连
func ConnectionClose(close bool) OptFunc {
	return func(request *Request) {
		request.ConnectionClose = close
	}
}

// RestyRequest request client
func (r *Request) RestyRequest() *resty.Request {
	client := resty.New()
	if r.Timeout > 0 {
		client.SetTimeout(r.Timeout)
	} else {
		client.SetTimeout(DefaultTimeout)
	}
	if r.RetryCount > 0 {
		client.SetRetryCount(r.RetryCount)
	} else {
		client.SetRetryCount(DefaultRetryCount)
	}
	if r.RetryWaitTime > 0 {
		client.SetRetryWaitTime(r.RetryWaitTime)
	} else {
		client.SetRetryWaitTime(DefaultRetryWaitTime)
	}
	client.AddRetryCondition(func(r *resty.Response, err error) bool {
		if r.IsError() {
			return true
		}
		if err != nil {
			sentry.CaptureException(errors.Errorf("%s,error:%s", "retry", err.Error()))
			return true
		}
		return false
	})
	client.SetCloseConnection(r.ConnectionClose)
	return client.R()
}

// GetResult 通过get 方法获取返回结果
func (r *Request) GetResult(ctx context.Context, url string, headers map[string]string, result interface{}) (err error) {
	tracer := opentracing.GlobalTracer()
	span, spCtx := opentracing.StartSpanFromContextWithTracer(ctx, tracer, url)
	defer span.Finish()
	restyReq := r.RestyRequest()
	//_ = tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(restyReq.Header))
	var response *resty.Response
	if response, err = restyReq.SetContext(spCtx).SetHeaders(headers).SetResult(&result).Get(url); err != nil {
		sentry.CaptureException(err)
	}
	if response.StatusCode() != http.StatusOK || !resty.IsJSONType(response.Header().Get(ContentType)) {
		if err = json.Unmarshal(response.Body(), &result); err != nil {
			sentry.CaptureException(err)
		}
	}
	return
}

// Get 通过get 方法获取返回结果
func (r *Request) Get(ctx context.Context, url string, headers map[string]string) (respBody []byte, err error) {
	tracer := opentracing.GlobalTracer()
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, tracer, url)
	defer span.Finish()
	restyReq := r.RestyRequest()
	if tracer != nil {
		err = tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(restyReq.EnableTrace().Header))
	}
	var response *resty.Response
	if response, err = restyReq.SetHeaders(headers).Get(url); err != nil {
		sentry.CaptureException(err)
	}
	respBody = response.Body()
	return
}

// PostResult 通过post 方法获取返回结果，并将结果存储到result 中
func (r *Request) PostResult(ctx context.Context, url string, body interface{}, headers map[string]string, result interface{}) (err error) {
	tracer := opentracing.GlobalTracer()
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, tracer, url)
	defer span.Finish()
	restyReq := r.RestyRequest()
	req := restyReq.RawRequest
	_ = tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	var response *resty.Response
	if response, err = restyReq.SetHeaders(headers).SetBody(body).SetResult(&result).Post(url); err != nil {
		sentry.CaptureException(err)
	}
	if response.StatusCode() != http.StatusOK || !resty.IsJSONType(response.Header().Get(ContentType)) {
		if err = json.Unmarshal(response.Body(), &result); err != nil {
			sentry.CaptureException(err)
		}
	}
	return
}

// Post 通过post 方法获取返回结果，并将结果存储到result 中
func (r *Request) Post(ctx context.Context, url string, body interface{}, headers map[string]string) (respBody []byte, err error) {
	tracer := opentracing.GlobalTracer()
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, tracer, url)
	defer span.Finish()
	restyReq := r.RestyRequest()
	req := restyReq.RawRequest
	_ = tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	var response *resty.Response
	if response, err = restyReq.SetHeaders(headers).SetBody(body).Post(url); err != nil {
		sentry.CaptureException(err)
	}
	respBody = response.Body()
	return
}
