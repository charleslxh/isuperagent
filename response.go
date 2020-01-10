package isuperagent

import (
	"io/ioutil"
	"net/http"

	"isuperagent/bodyParser"
)

type ResponseInterface interface {
	IsOk() bool
}

type BodyInterface interface {
	GetRaw() string
	Unmarshal(v interface{}) error
}

type Response struct {
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

func NewResponse(req *http.Request, resp *http.Response) (*Response, error) {
	res := &Response{}

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

func (r *Response) IsOk() bool {
	return r.StatusCode == 200
}
