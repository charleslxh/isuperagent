package isuperagent

import (
	"context"
	"errors"
	"fmt"
)

type MiddlewareInterface interface {
	Name() string
	Run(ctx context.Context, req *Request, next Next) (*Response, error)
}

// Middleware factory pool
var middlewareFactories map[string]MiddlewareFactory

// Middleware factory method
// it is convenient to create an middleware
type MiddlewareFactory func(v ...interface{}) (MiddlewareInterface, error)

// The method to call the next middleware
// Only middleware call next() method, the next middleware will invoke
// At last, the request will send to target server
type Next func() (*Response, error)

// Register your middleware to factory pool
func RegisterMiddlewareFactory(name string, factory MiddlewareFactory) {
	if middlewareFactories == nil {
		middlewareFactories = make(map[string]MiddlewareFactory, 0)
	}

	middlewareFactories[name] = factory
}

// Composer all middleware
// Return the start middleware
func Compose(ctx context.Context, middleware []MiddlewareInterface, req *Request) func(ctx context.Context, req *Request) (*Response, error) {
	i := 0
	var next Next
	next = func() (*Response, error) {
		i++
		if i >= len(middleware) {
			return nil, nil
		}

		if middleware[i] == nil {
			return nil, nil
		}

		return middleware[i].Run(ctx, req, next)
	}

	return func(ctx context.Context, req *Request) (*Response, error) {
		return middleware[i].Run(ctx, req, next)
	}
}

// The factory method to create a new middleware
// tip: you must register your middleware firstly by invoke isuperagent.RegisterMiddleware() method
func NewMiddleware(name string, v ...interface{}) (MiddlewareInterface, error) {
	if factory, ok := middlewareFactories[name]; ok {
		return factory(v...)
	}

	return nil, errors.New(fmt.Sprintf("middleware %s not registered", name))
}
