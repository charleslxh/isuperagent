package middleware

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"

	"isuperagent"
)

// Middleware: send request
//
// The last of middleware chain, it will call after all middleware.
// This middleware will do the following duties:
// 1. Generate the request object.
// 2. Set all request options, include queries, headers, bodies, authorization.
// 3. Send request, Generate response object.
// 4. Retry strategy when failure.
// Then return *Response, error to previous middleware.
type ExecRequestMiddleware struct{}

func NewExecRequestMiddlewareFactory(v ...interface{}) (isuperagent.MiddlewareInterface, error) {
	return &ExecRequestMiddleware{}, nil
}

func init() {
	isuperagent.RegisterMiddlewareFactory("request_exec", NewExecRequestMiddlewareFactory)
}

func (m *ExecRequestMiddleware) Name() string {
	return "request_exec"
}

func (m *ExecRequestMiddleware) Run(ctx context.Context, r *isuperagent.Request, next isuperagent.Next) (*isuperagent.Response, error) {
	// generate request body
	requestBody, err := r.GetBodyRaw()
	if err != nil {
		return nil, err
	}

	// create request
	req, err := http.NewRequest(r.GetMethod(), r.GetUrl(), bytes.NewReader(requestBody))
	if err != nil {
		return nil, err
	}

	// set query string
	req.URL.RawQuery = r.GetQueries().Encode()

	// set headers
	if r.GetHeader("Host") == "" {
		r.Header("Host", r.GetHost())
	}
	req.Header = r.GetHeaders()

	// send request
	c := http.Client{
		Timeout: r.GetTimeout(),
	}

	// set https options
	if r.IsHttps() {
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
				return nil, err
			}

			pool := x509.NewCertPool()
			pool.AppendCertsFromPEM(cert)
			tr.TLSClientConfig.RootCAs = pool
		}

		// Set client certificate
		if cert, key := r.GetCert(); cert != "" && key != "" {
			clientCert, err := tls.LoadX509KeyPair(cert, key)
			if err != nil {
				return nil, err
			}
			tr.TLSClientConfig.Certificates = []tls.Certificate{clientCert}
		}

		c.Transport = tr
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	res, err := isuperagent.NewResponse(req, resp)
	if err != nil {
		return nil, err
	}

	return res, nil
}
