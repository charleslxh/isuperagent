package bodyParser

import "encoding/json"

type JsonParser struct{}

func init() {
	Register("json", []string{
		"application/json",
		"application/javascript",
		"application/ld+json",
	}, &JsonParser{})
}

func (p *JsonParser) Unmarshal(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	return nil
}

func (p *JsonParser) Marshal(v interface{}) ([]byte, error) {
	bs, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return bs, nil
}
