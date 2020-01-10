package bodyParser

import (
	"errors"
	"reflect"

	"ptapp.cn/util/filter/mtype"
)

type TextParser struct{}

func init() {
	Register("text", []string{
		"text/plain", "text/css",
		"text/csv", "text/javascript",
		"text/plain", "text/xml",
	}, &TextParser{})
}

func (p *TextParser) Unmarshal(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("txt: Unmarshal(non-pointer " + rv.Type().String() + ")")
	}

	*v.(*string) = string(data)

	return nil
}

func (p *TextParser) Marshal(v interface{}) ([]byte, error) {
	return mtype.GetByte(v), nil
}
