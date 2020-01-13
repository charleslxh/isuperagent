package isuperagent_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/charleslxh/isuperagent"
	"github.com/charleslxh/isuperagent/middleware"
)

type ResponseData struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

func MockHttp(ast *assert.Assertions, ch <-chan struct{}) {
	mux := http.NewServeMux()

	mux.HandleFunc("/getQuery", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Server", runtime.Version())
		w.WriteHeader(200)

		data := &ResponseData{Code: 0, Msg: "OK"}
		data.Data = map[string]interface{}{}

		for key, value := range r.URL.Query() {
			data.Data[key] = value
		}

		bs, err := json.Marshal(data)
		ast.Nil(err)

		_, err = w.Write(bs)
		ast.Nil(err)
	})

	mux.HandleFunc("/getHeader", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Server", runtime.Version())
		w.WriteHeader(200)

		data := &ResponseData{Code: 0, Msg: "OK"}
		data.Data = map[string]interface{}{}

		for key, value := range r.Header {
			data.Data[key] = value
		}

		bs, err := json.Marshal(data)
		ast.Nil(err)

		_, err = w.Write(bs)
		ast.Nil(err)
	})

	mux.HandleFunc("/getHtml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("X-Server", runtime.Version())
		w.WriteHeader(200)

		html := `<html><head>Welcome</head><body>Welcome to golang world</body></html>`
		_, err := w.Write([]byte(html))
		ast.Nil(err)
	})

	mux.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		bs, err := ioutil.ReadAll(r.Body)
		ast.Nil(err)
		ast.Equal("POST", r.Method)

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("X-Server", runtime.Version())
		w.WriteHeader(200)

		_, err = w.Write(bs)
		ast.Nil(err)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("X-Server", runtime.Version())
		w.WriteHeader(200)

		text := `Hello world`
		_, err := w.Write([]byte(text))
		ast.Nil(err)
	})

	srv := http.Server{Addr: ":28080", Handler: mux}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println("Mock HTTP shutdown")
		}
	}()

	log.Println("Mock HTTPS server started on 28080")

	<-ch

	_ = srv.Shutdown(nil)
}

func MockHttps(ast *assert.Assertions, ch <-chan struct{}) {
	mux := http.NewServeMux()

	mux.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		bs, err := ioutil.ReadAll(r.Body)
		ast.Nil(err)
		ast.Equal("POST", r.Method)

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("X-Server", runtime.Version())
		w.WriteHeader(200)

		_, err = w.Write(bs)
		ast.Nil(err)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("X-Server", runtime.Version())
		w.WriteHeader(200)

		text := `Hello world HTTPS`
		_, err := w.Write([]byte(text))
		ast.Nil(err)
	})

	srv := http.Server{Addr: ":2443", Handler: mux}
	go func() {
		if err := srv.ListenAndServeTLS("./cert/cert.pem", "./cert/key.pem"); err != nil {
			log.Println("Mock HTTPS shutdown")
		}
	}()

	log.Println("Mock HTTPS server started on 2443")

	<-ch

	_ = srv.Shutdown(nil)
}

func startServer(ast *assert.Assertions, https bool) chan<- struct{} {
	ch := make(chan struct{})

	if https {
		go MockHttps(ast, ch)
	} else {
		go MockHttp(ast, ch)
	}

	return ch
}

func TestSuperAgent_Request(t *testing.T) {
	ast := assert.New(t)

	ch := startServer(ast, false)
	defer close(ch)

	// -----------------
	// 验证 QUERY 相关
	// -----------------
	queries := map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}

	res, err := isuperagent.NewRequest().Get("http://localhost:28080/getQuery").Queries(queries).Do()

	ast.Nil(err)
	ast.Equal(200, res.StatusCode)
	ast.True(res.IsOk())

	data1 := struct {
		Code int                 `json:"code"`
		Msg  string              `json:"msg"`
		Data map[string][]string `json:"data"`
	}{}
	err = res.Body.Unmarshal(&data1)
	ast.Nil(err)
	ast.Equal(map[string][]string{
		"a": {"1"},
		"b": {"2"},
		"c": {"3"},
	}, data1.Data)
	ast.Equal(0, data1.Code)
	ast.Equal("OK", data1.Msg)

	// -----------------
	// 验证 HEADER 相关
	// -----------------
	headers := map[string]string{
		"referer":         "https://golang.org/pkg/net/http/",
		"accept-encoding": "gzip, deflate, br",
		"cookie":          "_ga=GA1.2.242301321.1564383471; __utmc=110886291; __utmz=110886291.1578538723.31.21.utmcsr=google|utmccn=(organic)|utmcmd=organic|utmctr=(not%20provided); __utma=110886291.242301321.1564383471.1578538723.1578565071.32",
		"user-agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.117 Safari/537.36",
	}
	res, err = isuperagent.NewRequest().Get("http://localhost:28080/getHeader").Headers(headers).Do()
	ast.Nil(err)
	ast.Equal(200, res.StatusCode)
	ast.True(res.IsOk())

	data2 := struct {
		Code int                 `json:"code"`
		Msg  string              `json:"msg"`
		Data map[string][]string `json:"data"`
	}{}
	err = res.Body.Unmarshal(&data2)
	ast.Nil(err)
	ast.Equal(0, data2.Code)
	ast.Equal("OK", data2.Msg)
	ast.Equal(headers["accept-encoding"], data2.Data["Accept-Encoding"][0])
	ast.Equal(headers["referer"], data2.Data["Referer"][0])
	ast.Equal(headers["user-agent"], data2.Data["User-Agent"][0])
	ast.Equal(headers["cookie"], data2.Data["Cookie"][0])

	// -----------------
	// 验证 POST 请求
	// -----------------
	res, err = isuperagent.NewRequest().Post("http://localhost:28080/echo").Body("Hello World").Do()
	ast.Nil(err)
	ast.Equal(200, res.StatusCode)
	ast.True(res.IsOk())

	var data3 string
	err = res.Body.Unmarshal(&data3)
	ast.Nil(err)
	ast.Equal("Hello World", data3)

	res, err = isuperagent.NewRequest().Post("http://localhost:28080/echo", "Hello World Golang").Do()
	ast.Nil(err)
	ast.Equal(200, res.StatusCode)
	ast.True(res.IsOk())

	var data4 string
	err = res.Body.Unmarshal(&data4)
	ast.Nil(err)
	ast.Equal("Hello World Golang", data4)
}

func TestSuperAgent_HttpsRequest(t *testing.T) {
	ast := assert.New(t)

	ch := startServer(ast, true)
	defer close(ch)

	// -------------------
	// 请求知名机构签发的证书
	// -------------------
	res, err := isuperagent.NewRequest().Get("https://www.baidu.com").Do()

	ast.Nil(err)
	ast.Equal(200, res.StatusCode)
	ast.True(res.IsOk())

	var html string
	err = res.Body.Unmarshal(&html)
	ast.Nil(err)

	// -------------------
	// 请求自签名证书
	// -------------------
	res, err = isuperagent.NewRequest().Get("https://localhost:2443/").Do()
	ast.NotNil(err)
	ast.Equal("Get https://localhost:2443/: x509: certificate signed by unknown authority", err.Error())

	res, err = isuperagent.NewRequest().Get("https://localhost:2443/").InsecureSkipVerify(true).Do()
	ast.Nil(err)
	ast.Equal(200, res.StatusCode)
	ast.True(res.IsOk())

	err = res.Body.Unmarshal(&html)
	ast.Nil(err)
	ast.Equal("Hello world HTTPS", html)

	// TODO：增加客户端配置服务端根证书认证 CASE（单向认证）
	// TODO：增加双向认证 CASE
}

func TestSuperAgent_RequestTimeMiddleware(t *testing.T) {
	ast := assert.New(t)

	ch := startServer(ast, false)
	defer close(ch)

	timeMiddleware, err := isuperagent.NewMiddleware("request_time", 1)
	ast.NotNil(err)
	ast.Nil(timeMiddleware)
	ast.Equal("excepted header_name is string, but got 1(int)", err.Error())

	timeMiddleware, err = isuperagent.NewMiddleware("request_time")
	ast.Nil(err)

	res, err := isuperagent.NewRequest().Get("http://localhost:28080/").Middleware(timeMiddleware).Do()
	ast.Nil(err)
	ast.Equal(200, res.StatusCode)
	ast.True(res.IsOk())

	var data string
	err = res.Body.Unmarshal(&data)
	ast.Nil(err)
	ast.Equal("Hello world", data)
	ast.NotEqual("0", res.Headers.Get(middleware.X_SUPERAGENT_DURATION))
}

func TestSuperAgent_Middleware(t *testing.T) {
	ast := assert.New(t)

	ch := startServer(ast, false)
	defer close(ch)

	nonExistsMiddleware, err := isuperagent.NewMiddleware("non_exists")
	ast.NotNil(err)
	ast.Nil(nonExistsMiddleware)
	ast.Equal("middleware non_exists not registered", err.Error())

	timeMiddleware, err := isuperagent.NewMiddleware("request_time")
	ast.Nil(err)

	basicAuthMiddleware, err := isuperagent.NewMiddleware("basic_auth")
	ast.NotNil(err)
	ast.Nil(basicAuthMiddleware)
	ast.Equal("excepted two arguments, the first is username, next is password", err.Error())

	basicAuthMiddleware, err = isuperagent.NewMiddleware("basic_auth", 1, 1)
	ast.NotNil(err)
	ast.Nil(basicAuthMiddleware)
	ast.Equal("excepted username is string, but got 1(int)", err.Error())

	basicAuthMiddleware, err = isuperagent.NewMiddleware("basic_auth", "MTIzNDU2", "YXNkZmdoamts")
	ast.Nil(err)

	debugMiddleware, err := isuperagent.NewMiddleware("debug", func(ctx context.Context, req *isuperagent.Request) {
		ast.Equal(isuperagent.Method_GET, req.GetMethod())
		ast.NotEqual(0, len(req.GetHeader(middleware.BASIC_AUTH_HEADER)))
		log.Println(fmt.Sprintf("req headers: %+v", req.GetHeaders()))
	})
	ast.Nil(err)

	res, err := isuperagent.NewRequest().Get("http://localhost:28080/").
		Middleware(timeMiddleware, basicAuthMiddleware, debugMiddleware).
		Do()
	ast.Nil(err)
	ast.Equal(200, res.StatusCode)
	ast.True(res.IsOk())
}
