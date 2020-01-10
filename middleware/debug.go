package middleware

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/charleslxh/isuperagent"
)

// Middleware: this middleware allow you to debug request information and response
type DebugMiddleware struct {
	cb func(ctx context.Context, req *isuperagent.Request)
}

func DefaultCallback(ctx context.Context, req *isuperagent.Request) {
	// do nothing
}

func NewDebugMiddlewareFactory(v ...interface{}) (isuperagent.MiddlewareInterface, error) {
	if len(v) < 1 {
		return &DebugMiddleware{cb: DefaultCallback}, nil
	}

	if fn, ok := v[0].(func(ctx context.Context, req *isuperagent.Request)); ok {
		return &DebugMiddleware{cb: fn}, nil
	}

	return nil, errors.New(fmt.Sprintf("excepted first argument is a func(ctx context.Context, req *isuperagent.Request), but got %s", reflect.TypeOf(v[0])))
}

func init() {
	isuperagent.RegisterMiddlewareFactory("debug", NewDebugMiddlewareFactory)
}

func (m *DebugMiddleware) Name() string {
	return "debug"
}

func (m *DebugMiddleware) Run(ctx context.Context, req *isuperagent.Request, next isuperagent.Next) (*isuperagent.Response, error) {
	m.cb(ctx, req)

	return next()
}
