package middleware

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"

	"github.com/charleslxh/isuperagent"
)

// Middleware: Basic Auth
type BasicAuthMiddleware struct {
	username string
	password string
}

const BASIC_AUTH_HEADER = "Authorization"

func NewBasicAuthMiddlewareFactory(v ...interface{}) (isuperagent.MiddlewareInterface, error) {
	if len(v) < 2 {
		return nil, errors.New("excepted two arguments, the first is username, next is password")
	}

	middleware := &BasicAuthMiddleware{}

	if username, ok := v[0].(string); !ok {
		return nil, errors.New(fmt.Sprintf("excepted username is string, but got %v(%s)", v[0], reflect.TypeOf(v[0])))
	} else {
		middleware.username = username
	}

	if password, ok := v[1].(string); !ok {
		return nil, errors.New(fmt.Sprintf("excepted password is string, but got %v(%s)", v[0], reflect.TypeOf(v[0])))
	} else {
		middleware.password = password
	}

	return middleware, nil
}

func init() {
	isuperagent.RegisterMiddlewareFactory("basic_auth", NewBasicAuthMiddlewareFactory)
}

func (m *BasicAuthMiddleware) Name() string {
	return "basic_auth"
}

func (m *BasicAuthMiddleware) Run(ctx context.Context, req *isuperagent.Request, next isuperagent.Next) (*isuperagent.Response, error) {
	req.Header(BASIC_AUTH_HEADER, "Basic "+basicAuth(m.username, m.password))

	return next()
}

// See 2 (end of page 4) https://www.ietf.org/rfc/rfc2617.txt
// "To receive authorization, the client sends the userid and password,
// separated by a single colon (":") character, within a base64
// encoded string in the credentials."
// It is not meant to be urlencoded.
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}