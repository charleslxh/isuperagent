package middleware

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/charleslxh/isuperagent"
)

const X_SUPERAGENT_DURATION = "x-SuperAgent-Duration"

// Middleware: record the request duration
type TimeMiddleware struct {
	headerName string
}

func NewTimeMiddlewareFactory(v ...interface{}) (isuperagent.MiddlewareInterface, error) {
	if len(v) == 0 {
		return &TimeMiddleware{headerName: X_SUPERAGENT_DURATION}, nil
	}

	if h, ok := v[0].(string); ok {
		return &TimeMiddleware{headerName: h}, nil
	}

	return nil, errors.New(fmt.Sprintf("excepted header_name is string, but got %v(%s)", v[0], reflect.TypeOf(v[0])))
}

func init() {
	isuperagent.RegisterMiddlewareFactory("request_time", NewTimeMiddlewareFactory)
}

func (m *TimeMiddleware) Name() string {
	return "request_time"
}

func (m *TimeMiddleware) Run(ctx context.Context, req *isuperagent.Request, next isuperagent.Next) (*isuperagent.Response, error) {
	start := time.Now()

	res, err := next()
	if res != nil {
		res.Headers.Set(X_SUPERAGENT_DURATION, fmt.Sprintf("%s", time.Now().Sub(start)))
	}

	return res, err
}
