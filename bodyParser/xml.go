package bodyParser

import (
	"encoding/xml"
)

type XmlParser struct{}

func init() {
	Register("xml", []string{
		"application/xml",
	}, &XmlParser{})
}

func (p *XmlParser) Unmarshal(data []byte, v interface{}) error {
	err := xml.Unmarshal(data, v)
	if err != nil {
		return err
	}

	return nil
}

func (p *XmlParser) Marshal(v interface{}) ([]byte, error) {
	bs, err := xml.Marshal(v)
	if err != nil {
		return nil, err
	}

	return bs, nil
}
