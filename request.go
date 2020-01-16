package isuperagent

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/charleslxh/isuperagent/bodyParser"
)

type Request interface {
	SetMethod(method string, options ...interface{}) Request
	GetMethod() string

	Get(url string, options ...interface{}) Request
	Post(url string, options ...interface{}) Request
	Head(url string, options ...interface{}) Request
	Put(url string, options ...interface{}) Request
	Update(url string, options ...interface{}) Request
	Delete(url string, options ...interface{}) Request

	SetUrl(url string) Request
	GetUrl() *URL
	GetRawUrl() string

	GetHeader(name string) string
	SetHeader(name, value string) Request
	GetHeaders() http.Header
	SetHeaders(kv map[string]string) Request

	SetContentType(contentType string) Request

	GetQuery(name string) string
	SetQuery(name string, value string) Request
	GetQueries() url.Values
	SetQueries(kv map[string]string) Request

	SetBody(v interface{}) Request
	GetBody() interface{}
	GetBodyRaw() ([]byte, error)

	SetTimeout(d time.Duration) Request
	GetTimeout() time.Duration
	SetRetry(times int) Request
	GetRetry() int

	SetInsecureSkipVerify(insecureSkipVerify bool) Request
	GetInsecureSkipVerify() bool
	SetTlsConfig(tlsConfig *tls.Config) Request
	GetTlsConfig() *tls.Config
	SetCa(caPath string) Request
	GetCa() string
	SetCert(certPath, keyPath string) Request
	GetCert() (string, string)
	BasicAuth(name, pass string) Request
	GetUsername() string
	GetPassword() string

	SetContext(ctx context.Context) Request

	Middleware(middleware ...Middleware) Request

	Do() (Response, error)
}

type irequest struct {
	Context context.Context

	Method string
	Url    *URL

	ContentType ContentType
	Timeout     time.Duration
	Retry       int

	Headers http.Header

	// Optionally override the trusted CA certificates.
	// Default is to trust the well-known CAs curated by Mozilla
	Ca string
	// Cert chains in PEM format.
	// One cert chain should be provided per private key
	Cert string
	// Private keys in PEM format.
	// PEM allows the option of private keys being encrypted
	Key string
	// If false, the server certificate is verified against the list of supplied CAs.
	// An 'error' event is emitted if verification fails; err.code contains the OpenSSL error code.
	// Default: false.
	InsecureSkipVerify bool

	TlsConfig *tls.Config

	Body    interface{}
	BodyRaw []byte

	// Basic Auth
	Username string
	Password string

	Middlewares []Middleware
}

const (
	Method_GET    = "GET"
	Method_POST   = "POST"
	Method_HEAD   = "HEAD"
	Method_PUT    = "PUT"
	Method_UPDATE = "UPDATE"
	Method_DELETE = "DELETE"
)

func NewRequest() Request {
	return &irequest{Context: context.Background(), Url: NewURL(), Headers: http.Header{}}
}

func NewRequestWithContext(ctx context.Context) Request {
	return &irequest{Context: ctx, Url: NewURL(), Headers: http.Header{}}
}

func (r *irequest) SetContext(ctx context.Context) Request {
	r.Context = ctx

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
func (r *irequest) SetMethod(method string, options ...interface{}) Request {
	r.Method = strings.ToUpper(method)

	if len(options) > 0 {
		if v, ok := options[0].(string); ok {
			r.SetUrl(v)
		}
	}

	if len(options) > 1 {
		r.SetBody(options[1])
	}

	if len(options) > 2 {
		if v, ok := options[2].(map[string]string); ok {
			r.SetHeaders(v)
		}
	}

	if len(options) > 3 {
		if v, ok := options[3].(map[string]string); ok {
			r.SetQueries(v)
		}
	}

	return r
}

func (r *irequest) GetMethod() string {
	return r.Method
}

// Set request URL
func (r *irequest) SetUrl(url string) Request {
	// TODO: ... ...
	return r
}

func (r *irequest) GetUrl() *URL {
	return r.Url
}

func (r *irequest) GetRawUrl() string {
	return r.Url.String()
}

// Set request Host
func (r *irequest) Host(host string) Request {
	r.Url.Host = host

	return r
}

func (r *irequest) GetHost() string {
	return r.Url.Host
}

// Set request method to GET and request URL
// The options same as POST method.
// But body fields is nil.
func (r *irequest) Get(url string, options ...interface{}) Request {
	options = append([]interface{}{url, nil}, options...)
	r.SetMethod(Method_GET, url)

	return r
}

// Set POST options, method, url, body, header, query string
// The first argument is url, it is required.
// All of other arguments are not required, they can be set by other functions, such as Header(), Body() and so on.
// options definition as the following:
// 1. Set request body if first options exists.
// 2. Set request headers if second options exists, it must be map[string]string type, ignore is otherwise.
// 3. Set request queries if third options exists, it must be map[string]string type, ignore is otherwise.
func (r *irequest) Post(url string, options ...interface{}) Request {
	options = append([]interface{}{url}, options...)
	r.SetMethod(Method_POST, options...)

	return r
}

// Set request method to HEAD and request URL
// The options same as POST method.
// But body fields is nil.
func (r *irequest) Head(url string, options ...interface{}) Request {
	options = append([]interface{}{url, nil}, options...)
	r.SetMethod(Method_HEAD, url)

	return r
}

// Set Put options, url, method, query string, body, header
// The options same as POST method.
func (r *irequest) Put(url string, options ...interface{}) Request {
	options = append([]interface{}{url}, options...)
	r.SetMethod(Method_PUT, options...)

	return r
}

// Set Update options, url, method, query string, body, header
// The options same as POST method.
func (r *irequest) Update(url string, options ...interface{}) Request {
	options = append([]interface{}{url}, options...)
	r.SetMethod(Method_UPDATE, options...)

	return r
}

// Set request method to HEAD and request URL
// The options same as POST method.
// But body fields is nil.
func (r *irequest) Delete(url string, options ...interface{}) Request {
	options = append([]interface{}{url, nil}, options...)
	r.SetMethod(Method_DELETE, url)

	return r
}

func (r *irequest) SetHeader(name, value string) Request {
	r.Headers.Add(name, value)

	return r
}

func (r *irequest) GetHeader(name string) string {
	return r.Headers.Get(name)
}

func (r *irequest) SetHeaders(kvs map[string]string) Request {
	for k, v := range kvs {
		r.SetHeader(k, v)
	}

	return r
}

func (r *irequest) GetHeaders() http.Header {
	return r.Headers
}

func (r *irequest) SetQuery(name, value string) Request {
	r.Url.Queries.Add(name, value)

	return r
}

func (r *irequest) GetQuery(name string) string {
	if v := r.Url.Queries.Get(name); v != "" {
		return v
	}

	return r.Url.Query().Get(name)
}

func (r *irequest) SetQueries(kvs map[string]string) Request {
	for k, v := range kvs {
		r.SetQuery(k, v)
	}

	return r
}

func (r *irequest) GetQueries() url.Values {
	return r.Url.Queries
}

func (r *irequest) SetContentType(contentType string) Request {
	r.ContentType = ParseContentType(contentType)

	return r
}

func (r *irequest) SetTimeout(d time.Duration) Request {
	r.Timeout = d

	return r
}

func (r *irequest) GetTimeout() time.Duration {
	return r.Timeout
}

func (r *irequest) SetRetry(times int) Request {
	r.Retry = times

	return r
}

func (r *irequest) GetRetry() int {
	return r.Retry
}

func (r *irequest) IsHttps() bool {
	return "https" == r.Url.Scheme
}

func (r *irequest) Middleware(middleware ...Middleware) Request {
	r.Middlewares = append(r.Middlewares, middleware...)

	return r
}

func (r *irequest) SetBody(v interface{}) Request {
	r.Body = v

	return r
}

func (r *irequest) GetBody() interface{} {
	return r.Body
}

func (r *irequest) GetBodyRaw() ([]byte, error) {
	// generate request body
	requestBody, err := bodyParser.Marshal(r.ContentType.MediaType, r.Body)
	if err != nil {
		return nil, err
	}

	return requestBody, nil
}

func (r *irequest) SetCert(certPath, keyPath string) Request {
	r.Cert = certPath
	r.Key = keyPath

	return r
}

func (r *irequest) GetCert() (string, string) {
	return r.Cert, r.Key
}

// Set server root certificate.
func (r *irequest) SetCa(caPath string) Request {
	r.Ca = caPath

	return r
}

func (r *irequest) GetCa() string {
	return r.Ca
}

func (r *irequest) BasicAuth(user, pass string) Request {
	r.Username = user
	r.Password = pass

	return r
}

// Set username for basic auth.
func (r *irequest) SetUsername(username string) Request {
	r.Username = username

	return r
}

func (r *irequest) GetUsername() string {
	return r.Username
}

// Set password for basic auth.
func (r *irequest) SetPassword(password string) Request {
	r.Password = password

	return r
}

func (r *irequest) GetPassword() string {
	return r.Password
}

// Ignore verify the self-signed certificate.
// If InsecureSkipVerify = false, and the server's certificate is self-signed,
// The request connection will failed with error:
//		x509: certificate signed by unknown authority.
// So you can set to true if you don't care about server's certificate.
func (r *irequest) SetInsecureSkipVerify(insecureSkipVerify bool) Request {
	r.InsecureSkipVerify = insecureSkipVerify

	return r
}

func (r *irequest) GetInsecureSkipVerify() bool {
	return r.InsecureSkipVerify
}

// Set SSL config, see tls.Config
func (r *irequest) SetTlsConfig(tlsConfig *tls.Config) Request {
	r.TlsConfig = tlsConfig

	return r
}

func (r *irequest) GetTlsConfig() *tls.Config {
	return r.TlsConfig
}

func (r *irequest) Do() (Response, error) {
	m, err := NewMiddleware("request_exec")
	if err != nil {
		return nil, err
	}

	ctx := NewContext(r.Context, r, nil)

	middleware := append(r.Middlewares, m)

	err = Compose(ctx, middleware)(ctx)
	if err != nil {
		return nil, err
	}

	return ctx.GetRes(), nil
}
