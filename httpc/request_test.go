package httpc

import (
	"context"
	"testing"

	"github.com/melonwool/tracing/jaegerc"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	jaegerc.NewJaegerTracer("127.0.0.1:6831", "testing", jaegerc.ConstType, 1)
	c := NewRequest()
	body, err := c.RestyRequest().Get(context.Background(), "https://httpbin.org/get")
	assert.Equal(t, err, nil)
	assert.NotEqual(t, body, nil)
}
