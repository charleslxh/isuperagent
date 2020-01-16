package isuperagent

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type Middleware func(ctx Context, next Next) error

// Middleware factory pool
var middlewareFactories map[string]MiddlewareFactory

// Middleware factory method
// it is convenient to create an middleware
type MiddlewareFactory func(v ...interface{}) (Middleware, error)

// The method to call the next middleware
// Only middleware call next() method, the next middleware will invoke
// At last, the request will send to target server
type Next func() error

// Register your middleware to factory pool
func RegisterMiddlewareFactory(name string, factory MiddlewareFactory) {
	if middlewareFactories == nil {
		middlewareFactories = make(map[string]MiddlewareFactory, 0)
	}

	middlewareFactories[name] = factory
}

// Composer all middleware
// Return the start middleware
func Compose(ctx Context, middleware []Middleware) func() error {
	i := 0
	var next Next
	next = func() error {
		i++
		if i >= len(middleware) {
			return nil
		}

		if middleware[i] == nil {
			return nil
		}

		return middleware[i](ctx, next)
	}

	return func() error {
		return middleware[i](ctx, next)
	}
}

// The factory method to create a new middleware
// tip: you must register your middleware firstly by invoke isuperagent.RegisterMiddleware() method
func NewMiddleware(name string, v ...interface{}) (Middleware, error) {
	if factory, ok := middlewareFactories[name]; ok {
		return factory(v...)
	}

	return nil, errors.New(fmt.Sprintf("middleware %s not registered", name))
}

// ==============================================================================================================
//                                             Middleware
// ==============================================================================================================

func init() {
	RegisterMiddlewareFactory("request_time", NewTimeMiddlewareFactory)
	RegisterMiddlewareFactory("debug", NewDebugMiddlewareFactory)
	RegisterMiddlewareFactory("basic_auth", NewBasicAuthMiddlewareFactory)
	RegisterMiddlewareFactory("request_exec", NewRequestExecMiddlewareFactory)
}

// Middleware: record the request duration
func NewTimeMiddlewareFactory(v ...interface{}) (Middleware, error) {
	return func(ctx Context, next Next) error {
		start := time.Now()
		defer func() {
			duration := fmt.Sprintf("%s", time.Now().Sub(start))
			ctx.Set("request_time", duration)
			ctx.GetRes().GetHeaders().Set("X-SuperAgent-Duration", duration)
		}()

		return next()
	}, nil
}

// Middleware: this middleware allow you to debug request information and response
func NewDebugMiddlewareFactory(v ...interface{}) (Middleware, error) {
	var cb func(ctx Context)
	if len(v) < 0 {
		cb = func(ctx Context) {}
	} else {
		if fn, ok := v[0].(func(ctx Context)); !ok {
			return nil, errors.New(fmt.Sprintf("excepted first argument is a func(ctx Context), but got %s", reflect.TypeOf(v[0])))
		} else {
			cb = fn
		}
	}

	return func(ctx Context, next Next) error {
		cb(ctx)
		return next()
	}, nil
}

// Middleware: HTTP Basic Auth
func NewBasicAuthMiddlewareFactory(v ...interface{}) (Middleware, error) {
	if len(v) < 2 {
		return nil, errors.New("excepted two arguments, the first is username, next is password")
	}

	var username, password string

	if user, ok := v[0].(string); !ok {
		return nil, errors.New(fmt.Sprintf("excepted username is string, but got %v(%s)", v[0], reflect.TypeOf(v[0])))
	} else {
		username = user
	}

	if pass, ok := v[1].(string); !ok {
		return nil, errors.New(fmt.Sprintf("excepted password is string, but got %v(%s)", v[1], reflect.TypeOf(v[1])))
	} else {
		password = pass
	}

	return func(ctx Context, next Next) error {
		token := "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
		ctx.GetReq().SetHeader("Authorization", token)

		return next()
	}, nil
}

// Middleware: send request
//
// The last of middleware chain, it will call after all middleware.
// This middleware will do the following duties:
// 1. Generate the request object.
// 2. Set all request options, include queries, headers, bodies, authorization.
// 3. Send request, Generate response object.
// 4. Retry strategy when failure.
// Then return *Response, error to previous middleware.
func NewRequestExecMiddlewareFactory(v ...interface{}) (Middleware, error) {
	return func(ctx Context, next Next) error {
		r := ctx.GetReq()

		// generate request body
		requestBody, err := r.GetBodyRaw()
		if err != nil {
			return err
		}

		// create request
		req, err := http.NewRequest(r.GetMethod(), r.GetRawUrl(), bytes.NewReader(requestBody))
		if err != nil {
			return err
		}

		// set query string
		req.URL.RawQuery = r.GetQueries().Encode()

		// set headers
		if r.GetHeader("Host") == "" {
			r.SetHeader("Host", r.GetUrl().Host)
		}
		req.Header = r.GetHeaders()

		// send request
		c := http.Client{
			Timeout: r.GetTimeout(),
		}

		// Set basic auth
		if r.GetUsername() != "" && r.GetPassword() != "" {
			req.SetBasicAuth(r.GetUsername(), r.GetPassword())
		}

		// set https options
		if "https" == strings.ToLower(r.GetUrl().Scheme) {
			tr := &http.Transport{}

			// set tls options
			if tlsConfig := r.GetTlsConfig(); tlsConfig != nil {
				tr.TLSClientConfig = tlsConfig
			} else {
				tr.TLSClientConfig = &tls.Config{
					InsecureSkipVerify: r.GetInsecureSkipVerify(),
				}
			}

			// Add server's root ca cert, verify the server certificate
			if ca := r.GetCa(); ca != "" {
				cert, err := ioutil.ReadFile(ca)
				if err != nil {
					return err
				}

				pool := x509.NewCertPool()
				pool.AppendCertsFromPEM(cert)
				tr.TLSClientConfig.RootCAs = pool
			}

			// Set client certificate
			if cert, key := r.GetCert(); cert != "" && key != "" {
				clientCert, err := tls.LoadX509KeyPair(cert, key)
				if err != nil {
					return err
				}
				tr.TLSClientConfig.Certificates = []tls.Certificate{clientCert}
			}

			c.Transport = tr
		}

		// Do request, retry it again if failed
		var resp *http.Response
		var e error
		if r.GetRetry() > 0 {
			for times := 0; times < r.GetRetry(); times++ {
				resp, e = c.Do(req)
				if e == nil {
					break
				}
			}

		} else {
			resp, e = c.Do(req)
		}
		if e != nil {
			return e
		}

		res, err := NewResponse(req, resp)
		if err != nil {
			return err
		}
		ctx.SetRes(res)

		return nil
	}, nil
}
