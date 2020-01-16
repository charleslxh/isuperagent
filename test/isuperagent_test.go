package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/charleslxh/isuperagent"
)

type ResponseData struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

func MockHttp(ast *assert.Assertions, port int, ch <-chan struct{}) {
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

	srv := http.Server{Addr: fmt.Sprintf(":%d", port), Handler: mux}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println("Mock HTTP shutdown")
		}
	}()

	log.Println(fmt.Sprintf("Mock HTTPS server started on %d", port))

	<-ch

	_ = srv.Shutdown(nil)
}

func MockHttps(ast *assert.Assertions, port int, ch <-chan struct{}) {
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

	srv := http.Server{Addr: fmt.Sprintf(":%d", port), Handler: mux}
	go func() {
		if err := srv.ListenAndServeTLS("../cert/cert.pem", "../cert/key.pem"); err != nil {
			log.Println("Mock HTTPS shutdown")
		}
	}()

	log.Println(fmt.Sprintf("Mock HTTPS server started on %d", port))

	<-ch

	_ = srv.Shutdown(nil)
}

func startServer(ast *assert.Assertions, port int, https bool) chan<- struct{} {
	ch := make(chan struct{})

	if https {
		go MockHttps(ast, port, ch)
	} else {
		go MockHttp(ast, port, ch)
	}

	return ch
}

func TestSuperAgent_Request(t *testing.T) {
	ast := assert.New(t)

	port := 28080
	ch := startServer(ast, port, false)
	defer close(ch)

	// -----------------
	// 验证 QUERY 相关
	// -----------------
	queries := map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}

	url := fmt.Sprintf("http://localhost:%d/getQuery", port)
	res, err := isuperagent.NewRequest().Get(url).SetQueries(queries).Do()
	ast.Nil(err)
	ast.Equal(200, res.GetStatusCode())
	ast.True(res.IsOk())

	data1 := struct {
		Code int                 `json:"code"`
		Msg  string              `json:"msg"`
		Data map[string][]string `json:"data"`
	}{}
	err = res.GetBody().Unmarshal(&data1)
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
	url = fmt.Sprintf("http://localhost:%d/getHeader", port)
	res, err = isuperagent.NewRequest().Get(url).SetHeaders(headers).Do()
	ast.Nil(err)
	ast.Equal(200, res.GetStatusCode())
	ast.True(res.IsOk())

	data2 := struct {
		Code int                 `json:"code"`
		Msg  string              `json:"msg"`
		Data map[string][]string `json:"data"`
	}{}
	err = res.ParseBody(&data2)
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
	url = fmt.Sprintf("http://localhost:%d/echo", port)
	res, err = isuperagent.NewRequest().Post(url).SetBody("Hello World").Do()
	ast.Nil(err)
	ast.Equal(200, res.GetStatusCode())
	ast.True(res.IsOk())

	var data3 string
	err = res.GetBody().Unmarshal(&data3)
	ast.Nil(err)
	ast.Equal("Hello World", data3)

	url = fmt.Sprintf("http://localhost:%d/echo", port)
	res, err = isuperagent.NewRequest().Post(url, "Hello World Golang").Do()
	ast.Nil(err)
	ast.Equal(200, res.GetStatusCode())
	ast.True(res.IsOk())

	var data4 string
	err = res.GetBody().Unmarshal(&data4)
	ast.Nil(err)
	ast.Equal("Hello World Golang", data4)
}

func TestSuperAgent_HttpsRequest(t *testing.T) {
	ast := assert.New(t)

	port := 28081
	ch := startServer(ast, port, true)
	defer close(ch)

	// -------------------
	// 请求知名机构签发的证书
	// -------------------
	res, err := isuperagent.NewRequest().Get("https://www.baidu.com").Do()

	ast.Nil(err)
	ast.Equal(200, res.GetStatusCode())
	ast.True(res.IsOk())

	var html string
	err = res.GetBody().Unmarshal(&html)
	ast.Nil(err)

	// -------------------
	// 请求自签名证书
	// -------------------
	url := fmt.Sprintf("https://localhost:%d/", port)
	res, err = isuperagent.NewRequest().Get(url).Do()
	ast.NotNil(err)
	ast.Equal(fmt.Sprintf("Get https://localhost:%d/: x509: certificate signed by unknown authority", port), err.Error())

	res, err = isuperagent.NewRequest().Get(url).SetInsecureSkipVerify(true).Do()
	ast.Nil(err)
	ast.Equal(200, res.GetStatusCode())
	ast.True(res.IsOk())

	err = res.GetBody().Unmarshal(&html)
	ast.Nil(err)
	ast.Equal("Hello world HTTPS", html)

	// TODO：增加客户端配置服务端根证书认证 CASE（单向认证）
	// TODO：增加双向认证 CASE
}

func TestSuperAgent_RequestTimeMiddleware(t *testing.T) {
	ast := assert.New(t)

	port := 28082
	ch := startServer(ast, port, false)
	defer close(ch)

	timeMiddleware, err := isuperagent.NewMiddleware("request_time")
	ast.Nil(err)

	url := fmt.Sprintf("http://localhost:%d/", port)
	res, err := isuperagent.NewRequest().Get(url).Middleware(timeMiddleware).Do()
	ast.Nil(err)
	ast.Equal(200, res.GetStatusCode())
	ast.True(res.IsOk())

	var data string
	err = res.GetBody().Unmarshal(&data)
	ast.Nil(err)
	ast.Equal("Hello world", data)
	ast.NotEqual("0", res.GetHeaders().Get("X-SuperAgent-Duration"))
}

func TestSuperAgent_Middleware(t *testing.T) {
	ast := assert.New(t)

	port := 28083
	ch := startServer(ast, port, false)
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

	debugMiddleware, err := isuperagent.NewMiddleware("debug", func(ctx isuperagent.Context) {
		ast.Equal(isuperagent.Method_GET, ctx.GetReq().GetMethod())
		ast.NotEqual(0, len(ctx.GetReq().GetHeader("Authorization")))
		log.Println(fmt.Sprintf("req headers: %+v", ctx.GetReq().GetHeaders()))
	})
	ast.Nil(err)

	url := fmt.Sprintf("http://localhost:%d", port)
	res, err := isuperagent.NewRequest().Get(url).
		Middleware(timeMiddleware, basicAuthMiddleware, debugMiddleware).
		Do()
	ast.Nil(err)
	ast.Equal(200, res.GetStatusCode())
	ast.True(res.IsOk())
}
