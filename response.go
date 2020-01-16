package isuperagent

import (
	"io/ioutil"
	"net/http"

	"github.com/charleslxh/isuperagent/bodyParser"
)

type Response interface {
	IsOk() bool
	GetHeaders() http.Header
	GetBody() *Body
	ParseBody(v interface{}) error
	GetStatusCode() int
	GetStatusText() string

	GetHttpRequest() *http.Request
	GetHttpResponse() *http.Response
}

type BodyInterface interface {
	GetRaw() string
	Unmarshal(v interface{}) error
}

type iresponse struct {
	StatusCode int
	StatusText string
	Body       *Body
	Headers    http.Header

	HttpReq  *http.Request
	HttpResp *http.Response
}

type Body struct {
	data        []byte
	contentType string
}

func (b *Body) GetData() []byte {
	return b.data
}

func (b *Body) Unmarshal(v interface{}) error {
	err := bodyParser.Unmarshal(b.contentType, b.data, v)
	if err != nil {
		return err
	}

	return nil
}

func NewResponse(req *http.Request, resp *http.Response) (Response, error) {
	res := &iresponse{}

	res.StatusCode = resp.StatusCode
	res.StatusText = resp.Status
	res.Headers = resp.Header

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	res.Body = &Body{data: content, contentType: resp.Header.Get("content-type")}

	res.HttpReq = req
	res.HttpResp = resp

	return res, nil
}

func (r *iresponse) IsOk() bool {
	return r.StatusCode == 200
}

func (r *iresponse) GetHeaders() http.Header {
	return r.Headers
}

func (r *iresponse) GetBody() *Body {
	return r.Body
}

func (r *iresponse) ParseBody(v interface{}) error {
	return r.Body.Unmarshal(v)
}

func (r *iresponse) GetStatusCode() int {
	return r.StatusCode
}

func (r *iresponse) GetStatusText() string {
	return r.StatusText
}

func (r *iresponse) GetHttpRequest() *http.Request {
	return r.HttpReq
}

func (r *iresponse) GetHttpResponse() *http.Response {
	return r.HttpResp
}
