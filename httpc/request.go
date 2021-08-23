package httpc

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go/ext"

	"github.com/opentracing/opentracing-go"

	"github.com/go-resty/resty/v2"
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
type restyRequest struct {
	Request *resty.Request
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
func (r *Request) RestyRequest() *restyRequest {
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
			//sentry.CaptureException(errors.Errorf("%s,error:%s", "retry", err.Error()))
			return true
		}
		return false
	})
	r.Client.SetCloseConnection(r.ConnectionClose)
	return &restyRequest{Request: r.Client.R()}
}

// GetResult 通过get 方法获取返回结果
func (r *restyRequest) GetResult(ctx context.Context, requestURL string, result interface{}) (err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, url2func(requestURL))
	_ = span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Request.EnableTrace().Header))
	var response *resty.Response
	ext.HTTPMethod.Set(span, http.MethodGet)
	ext.HTTPUrl.Set(span, requestURL)
	defer span.SetTag("result", result).Finish()
	if response, err = r.Request.SetResult(&result).Get(requestURL); err != nil {
		return
	}
	if response.StatusCode() != http.StatusOK || !resty.IsJSONType(response.Header().Get(ContentType)) {
		if err = json.Unmarshal(response.Body(), &result); err != nil {
			return
		}
	}
	return
}

// Get 通过get 方法获取返回结果
func (r *restyRequest) Get(ctx context.Context, requestURL string) (respBody []byte, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, url2func(requestURL))
	_ = span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Request.EnableTrace().Header))
	var response *resty.Response
	if response, err = r.Request.Get(requestURL); err == nil {
		respBody = response.Body()
		ext.HTTPMethod.Set(span, r.Request.Method)
		ext.HTTPUrl.Set(span, r.Request.URL)
	}
	span.SetTag("result", string(respBody)).Finish()
	return
}

// PostResult 通过post 方法获取返回结果，并将结果存储到result 中
func (r *restyRequest) PostResult(ctx context.Context, requestURL string, result interface{}) (err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, url2func(requestURL))
	_ = span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Request.EnableTrace().Header))
	ext.HTTPMethod.Set(span, http.MethodPost)
	ext.HTTPUrl.Set(span, requestURL)
	defer span.SetTag("result", result).Finish()
	var response *resty.Response
	if response, err = r.Request.SetResult(&result).Post(requestURL); err != nil {
		return
	}
	if response.StatusCode() != http.StatusOK || !resty.IsJSONType(response.Header().Get(ContentType)) {
		if err = json.Unmarshal(response.Body(), &result); err != nil {
			return
		}
	}
	return
}

// Post 通过post 方法获取返回结果，并将结果存储到result 中
func (r *restyRequest) Post(ctx context.Context, requestURL string) (respBody []byte, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, url2func(requestURL))
	_ = span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Request.EnableTrace().Header))
	var response *resty.Response
	if response, err = r.Request.Post(requestURL); err == nil {
		respBody = response.Body()
		ext.HTTPMethod.Set(span, r.Request.Method)
		ext.HTTPUrl.Set(span, r.Request.URL)
	}
	span.SetTag("result", string(respBody)).Finish()
	return
}

// SetFormData 设置post 参数
func (r *restyRequest) SetFormData(data map[string]string) *restyRequest {
	r.Request.SetFormData(data)
	return r
}

// SetHeaders 设置header
func (r *restyRequest) SetHeaders(headers map[string]string) *restyRequest {
	r.Request.SetHeaders(headers)
	return r
}

// SetBody 设置body
func (r *restyRequest) SetBody(body interface{}) *restyRequest {
	r.Request.SetBody(body)
	return r
}

// SetQueryParams 设置query 参数
func (r *restyRequest) SetQueryParams(params map[string]string) *restyRequest {
	r.Request.SetQueryParams(params)
	return r
}

func url2func(requestURL string) string {
	values, err := url.Parse(requestURL)
	if err != nil {
		return requestURL
	}
	paths := strings.Split(values.Path, "/")
	funcName := paths[len(paths)-1]
	if funcName == "" {
		funcName = values.Host
	}
	return "HTTP_Request:" + funcName
}
