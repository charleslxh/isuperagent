package bodyParser

import (
	"errors"
	"reflect"
	"strconv"
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
	return []byte(GetString(v)), nil
}

func GetString(obj interface{}) string {
	if obj == nil {
		return ""
	}

	switch v := obj.(type) {
	case bool:
		if v {
			return "1"
		} else {
			return ""
		}
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatInt(int64(v), 10)
	case uint8:
		return strconv.FormatInt(int64(v), 10)
	case uint16:
		return strconv.FormatInt(int64(v), 10)
	case uint32:
		return strconv.FormatInt(int64(v), 10)
	case uint64:
		return strconv.FormatInt(int64(v), 10)
	case string:
		return v
	default:
		return ""
	}
}
