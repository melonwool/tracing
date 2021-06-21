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
	Client *resty.Client
	Option
}

func NewRequest(optFunc ...OptFunc) *Request {
	request := &Request{
		Client: resty.New(),
	}
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
	if r.Timeout > 0 {
		r.Client.SetTimeout(r.Timeout)
	} else {
		r.Client.SetTimeout(DefaultTimeout)
	}
	if r.RetryCount > 0 {
		r.Client.SetRetryCount(r.RetryCount)
	} else {
		r.Client.SetRetryCount(DefaultRetryCount)
	}
	if r.RetryWaitTime > 0 {
		r.Client.SetRetryWaitTime(r.RetryWaitTime)
	} else {
		r.Client.SetRetryWaitTime(DefaultRetryWaitTime)
	}
	r.Client.AddRetryCondition(func(r *resty.Response, err error) bool {
		if r.IsError() {
			return true
		}
		if err != nil {
			sentry.CaptureException(errors.Errorf("%s,error:%s", "retry", err.Error()))
			return true
		}
		return false
	})
	r.Client.SetCloseConnection(r.ConnectionClose)
	return r.Client.R()
}

// GetResult 通过get 方法获取返回结果
func (r *Request) GetResult(ctx context.Context, url string, headers map[string]string, result interface{}) (err error) {
	tracer := opentracing.GlobalTracer()
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, tracer, url)
	defer span.Finish()
	restyReq := r.RestyRequest()
	_ = tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(restyReq.EnableTrace().Header))
	var response *resty.Response
	if response, err = restyReq.SetHeaders(headers).SetResult(&result).Get(url); err != nil {
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
	_ = tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(restyReq.EnableTrace().Header))
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
func (r *Request) Post(ctx context.Context, url string, body interface{}, formData map[string]string, headers map[string]string) (respBody []byte, err error) {
	tracer := opentracing.GlobalTracer()
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, tracer, url)
	defer span.Finish()
	restyReq := r.RestyRequest()
	_ = tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(restyReq.EnableTrace().Header))
	var response *resty.Response
	if response, err = restyReq.SetHeaders(headers).SetBody(body).SetFormData(formData).Post(url); err != nil {
		sentry.CaptureException(err)
	}
	respBody = response.Body()
	return
}
