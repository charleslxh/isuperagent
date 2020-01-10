package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"ptapp.cn/util/isuperagent"
)

const OMS_SNIPER_HEADER = "Authorization"

// Middleware: OMS 人群画像系统验签服务
// @see http://sx.eolinker.mengtuiapp.com/#/home/project/inside/doc/detail?groupID=-1&documentID=4&projectName=%E4%BA%BA%E7%BE%A4%E6%9C%8D%E5%8A%A1&projectID=48
type OmsSniperAuthMiddleware struct {
	accessKey  string
	secretKey  string
	headerName string
}

func NewOmsSniperAuthMiddlewareFactory(v ...interface{}) (isuperagent.MiddlewareInterface, error) {
	if len(v) < 2 {
		return nil, errors.New("access_key and secret_key is required, and the first argument must be access_key, next secret_key")
	}

	middleware := &OmsSniperAuthMiddleware{}

	if ak, ok := v[0].(string); !ok {
		return nil, errors.New(fmt.Sprintf("excepted access_key is string, but got %v(%s)", v[0], reflect.TypeOf(v[0])))
	} else {
		middleware.accessKey = ak
	}

	if sk, ok := v[1].(string); !ok {
		return nil, errors.New(fmt.Sprintf("excepted secret_key is string, but got %v(%s)", v[1], reflect.TypeOf(v[1])))
	} else {
		middleware.secretKey = sk
	}

	if len(v) > 2 {
		if name, ok := v[2].(string); !ok {
			return nil, errors.New(fmt.Sprintf("excepted header_name is string, but got %v(%s)", v[2], reflect.TypeOf(v[2])))
		} else {
			middleware.headerName = name
		}
	} else {
		middleware.headerName = OMS_SNIPER_HEADER
	}

	return middleware, nil
}

func init() {
	isuperagent.RegisterMiddlewareFactory("oms_sniper_auth", NewOmsSniperAuthMiddlewareFactory)
}

func (m *OmsSniperAuthMiddleware) Name() string {
	return "oms_sniper_auth"
}

func (m *OmsSniperAuthMiddleware) Run(ctx context.Context, req *isuperagent.Request, next isuperagent.Next) (*isuperagent.Response, error) {
	token, err := m.GenerateToken(ctx, req)
	if err != nil {
		return nil, err
	}

	req.Header(m.headerName, token)

	res, err := next()

	return res, err
}

func (m *OmsSniperAuthMiddleware) GenerateToken(ctx context.Context, req *isuperagent.Request) (string, error) {
	var data strings.Builder

	if len(req.GetUrl()) != 0 {
		data.WriteString(req.GetMethod())
		data.WriteString(req.GetUrl())
	}

	if len(req.GetQueries()) != 0 {
		data.WriteString(req.GetQueries().Encode())
	}

	data.WriteString("\n")
	data.WriteString(req.GetHost())
	data.WriteString("\n")
	data.WriteString(req.GetHeader("Content-Type"))
	data.WriteString("\n\n")

	requestBody, err := req.GetBodyRaw()
	if err != nil {
		return "", err
	}
	if len(requestBody) != 0 {
		data.WriteString(string(requestBody))
	}

	token := "mt" + m.accessKey + m.Digest([]byte(data.String()))

	return token, nil
}

func (m *OmsSniperAuthMiddleware) Digest(plainText []byte) string {
	mac := hmac.New(sha1.New, []byte(m.secretKey))
	mac.Write(plainText)
	encryptedText := mac.Sum(nil)

	signature := base64.URLEncoding.EncodeToString(encryptedText)

	return signature
}
