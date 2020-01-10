package isuperagent

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"
	"time"

	"ptapp.cn/util/isuperagent/bodyParser"
)

type RequestInterface interface {
	Method(method string) *Request
	Get(url string) *Request
	Post(url string) *Request
	Head(url string) *Request
	Put(url string) *Request
	Update(url string) *Request
	Delete(url string) *Request

	Url(url string) *Request

	Header(name, value string) *Request
	Headers(map[string]string) *Request

	Query(name string, value string) *Request
	Queries(map[string]string) *Request

	Body(v interface{}) *Request

	Timeout(d time.Duration) *Request
	Retry(times int) *Request

	Cert(certPath string) *Request
	Key(keyPath string) *Request
	BasicAuth(name, pass string) *Request

	Context(ctx context.Context) *Request

	Middleware(middleware ...MiddlewareInterface) *Request

	Do() (*Response, error)
}

type Request struct {
	ctx         context.Context
	host        string
	method      string
	url         string
	contentType string
	timeout     time.Duration
	retry       int

	headers http.Header
	queries url.Values

	// Optionally override the trusted CA certificates.
	// Default is to trust the well-known CAs curated by Mozilla
	ca string
	// Cert chains in PEM format.
	// One cert chain should be provided per private key
	cert string
	// Private keys in PEM format.
	// PEM allows the option of private keys being encrypted
	key string
	// If false, the server certificate is verified against the list of supplied CAs.
	// An 'error' event is emitted if verification fails; err.code contains the OpenSSL error code.
	// Default: false.
	insecureSkipVerify bool

	tlsConfig *tls.Config

	body    interface{}
	bodyRaw []byte

	middleware []MiddlewareInterface
}

const (
	Method_GET    = "GET"
	Method_POST   = "POST"
	Method_HEAD   = "HEAD"
	Method_PUT    = "PUT"
	Method_UPDATE = "UPDATE"
	Method_DELETE = "DELETE"
)

func NewRequest() *Request {
	return &Request{ctx: context.Background(), queries: url.Values{}, headers: http.Header{}}
}

func NewRequestWithContext(ctx context.Context) *Request {
	return &Request{ctx: ctx, queries: url.Values{}, headers: http.Header{}}
}

func (r *Request) Context(ctx context.Context) *Request {
	r.ctx = ctx

	return r
}

// Set request options, method, url, body, header, query string
// The first argument is url, it is required.
// All of other arguments are not required, they can be set by other functions, such as Header(), Body() and so on.
// options definition as the following:
// 1. Set request url if first options exists, it must be string type, ignore is otherwise..
// 2. Set request body if second options exists.
// 3. Set request headers if third options exists, it must be map[string]string type, ignore is otherwise.
// 4. Set request queries if fourth options exists, it must be map[string]string type, ignore is otherwise.
func (r *Request) Method(method string, options ...interface{}) *Request {
	r.method = strings.ToUpper(method)

	if len(options) > 0 {
		if v, ok := options[0].(string); ok {
			r.Url(v)
		}
	}

	if len(options) > 1 {
		r.Body(options[1])
	}

	if len(options) > 2 {
		if v, ok := options[2].(map[string]string); ok {
			r.Headers(v)
		}
	}

	if len(options) > 3 {
		if v, ok := options[3].(map[string]string); ok {
			r.Queries(v)
		}
	}

	return r
}

func (r *Request) GetMethod() string {
	return r.method
}

// Set request URL
func (r *Request) Url(url string) *Request {
	r.url = url

	return r
}

func (r *Request) GetUrl() string {
	return r.url
}

// Set request Host
func (r *Request) Host(host string) *Request {
	r.host = host

	return r
}

func (r *Request) GetHost() string {
	return r.host
}

// Set request method to GET and request URL
func (r *Request) Get(url string) *Request {
	r.Method(Method_GET, url)

	return r
}

// Set POST options, method, url, body, header, query string
// The first argument is url, it is required.
// All of other arguments are not required, they can be set by other functions, such as Header(), Body() and so on.
// options definition as the following:
// 1. Set request body if first options exists.
// 2. Set request headers if second options exists, it must be map[string]string type, ignore is otherwise.
// 3. Set request queries if third options exists, it must be map[string]string type, ignore is otherwise.
func (r *Request) Post(url string, options ...interface{}) *Request {
	options = append([]interface{}{url}, options...)
	r.Method(Method_POST, options...)

	return r
}

// Set request method to HEAD and request URL
func (r *Request) Head(url string) *Request {
	r.Method(Method_HEAD, url)

	return r
}

// Set Put options, url, method, query string, body, header
// The options same as POST method.
func (r *Request) Put(url string, options ...interface{}) *Request {
	options = append([]interface{}{url}, options...)
	r.Method(Method_PUT, options...)

	return r
}

// Set Update options, url, method, query string, body, header
// The options same as POST method.
func (r *Request) Update(url string, options ...interface{}) *Request {
	options = append([]interface{}{url}, options...)
	r.Method(Method_UPDATE, options...)

	return r
}

// Set request method to HEAD and request URL
func (r *Request) Delete(url string) *Request {
	r.Method(Method_DELETE, url)

	return r
}

func (r *Request) Header(name, value string) *Request {
	r.headers.Add(name, value)

	return r
}

func (r *Request) GetHeader(name string) string {
	return r.headers.Get(name)
}

func (r *Request) Headers(kvs map[string]string) *Request {
	for k, v := range kvs {
		r.Header(k, v)
	}

	return r
}

func (r *Request) GetHeaders() http.Header {
	return r.headers
}

func (r *Request) Query(name, value string) *Request {
	r.queries.Add(name, value)

	return r
}

func (r *Request) GetQuery(name string) string {
	return r.queries.Get(name)
}

func (r *Request) Queries(kvs map[string]string) *Request {
	for k, v := range kvs {
		r.Query(k, v)
	}

	return r
}

func (r *Request) GetQueries() url.Values {
	return r.queries
}

func (r *Request) ContentType(contentType string) *Request {
	r.Header("Content-Type", contentType)
	r.contentType = contentType

	return r
}

func (r *Request) Timeout(d time.Duration) *Request {
	r.timeout = d

	return r
}

func (r *Request) GetTimeout() time.Duration {
	return r.timeout
}

func (r *Request) Retry(times int) *Request {
	r.retry = times

	return r
}

func (r *Request) GetRetry() int {
	return r.retry
}

func (r *Request) IsHttps() bool {
	return strings.HasPrefix(r.url, "https")
}

func (r *Request) Middleware(middleware ...MiddlewareInterface) *Request {
	r.middleware = append(r.middleware, middleware...)

	return r
}

func (r *Request) Body(v interface{}) *Request {
	r.body = v

	return r
}

func (r *Request) GetBody() interface{} {
	return r.body
}

func (r *Request) GetBodyRaw() ([]byte, error) {
	// generate request body
	requestBody, err := bodyParser.Marshal(r.contentType, r.body)
	if err != nil {
		return nil, err
	}

	return requestBody, nil
}

func (r *Request) Cert(certPath, keyPath string) *Request {
	r.cert = certPath
	r.key = keyPath

	return r
}

func (r *Request) GetCert() (string, string) {
	return r.cert, r.key
}

// Set server root certificate.
func (r *Request) Ca(caPath string) *Request {
	r.ca = caPath

	return r
}

func (r *Request) GetCa() string {
	return r.ca
}

// Ignore verify the self-signed certificate.
// If InsecureSkipVerify = false, and the server's certificate is self-signed,
// The request connection will failed with error:
//		x509: certificate signed by unknown authority.
// So you can set to true if you don't care about server's certificate.
func (r *Request) InsecureSkipVerify(insecureSkipVerify bool) *Request {
	r.insecureSkipVerify = insecureSkipVerify

	return r
}

func (r *Request) GetInsecureSkipVerify() bool {
	return r.insecureSkipVerify
}

// Set SSL config, see tls.Config
func (r *Request) TlsConfig(tlsConfig *tls.Config) *Request {
	r.tlsConfig = tlsConfig

	return r
}

func (r *Request) GetTlsConfig() *tls.Config {
	return r.tlsConfig
}

func (r *Request) Do() (*Response, error) {
	m, err := NewMiddleware("request_exec")
	if err != nil {
		return nil, err
	}

	middleware := append(r.middleware, m)

	return Compose(r.ctx, middleware, r)(r.ctx, r)
}
