package bodyParser

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
)

type FormParser struct{}

func init() {
	Register("form", []string{
		"application/x-www-form-urlencoded",
	}, &FormParser{})
}

func (p *FormParser) Unmarshal(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("txt: Unmarshal(non-pointer " + rv.Type().String() + ")")
	}

	if _, ok := v.(url.Values); ok {
		val, err := url.ParseQuery(string(data))
		if err != nil {
			return err
		}

		*v.(*url.Values) = val

		return nil
	}

	if _, ok := v.(map[string][]string); ok {
		val, err := url.ParseQuery(string(data))
		if err != nil {
			return err
		}

		*v.(*map[string][]string) = val

		return nil
	}

	return errors.New(fmt.Sprintf("Unmarshal dest target must type of map[string][]string, but got %s", reflect.TypeOf(v)))
}

func (p *FormParser) Marshal(v interface{}) ([]byte, error) {
	if val, ok := v.(url.Values); ok {
		return []byte(val.Encode()), nil
	}

	if val, ok := v.(map[string][]string); ok {
		return []byte(url.Values(val).Encode()), nil
	}

	s := GetString(v)
	val, err := url.ParseQuery(s)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("invalid data, invalid request form data, data: %+v, type: %s", v, reflect.TypeOf(v)))
	}

	return []byte(val.Encode()), nil
}
