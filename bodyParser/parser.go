package bodyParser

import (
	"strings"
)

type BodyParserInterface interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

var parsers map[string]BodyParserInterface

var contentTypeAlias = map[string]string{}

func Register(alias string, contentTypes []string, parser BodyParserInterface) {
	if parsers == nil {
		parsers = make(map[string]BodyParserInterface, 0)
	}

	alias = strings.ToLower(alias)

	for _, c := range contentTypes {
		contentTypeAlias[strings.ToLower(c)] = alias
	}

	parsers[alias] = parser
}

func getParser(contentType string) BodyParserInterface {
	if parserName, ok := contentTypeAlias[strings.ToLower(contentType)]; ok {
		if parser, ok := parsers[parserName]; ok {
			return parser
		}
	}

	return parsers["text"]
}

func Unmarshal(contentType string, data []byte, v interface{}) error {
	return getParser(contentType).Unmarshal(data, v)
}

func Marshal(contentType string, v interface{}) ([]byte, error) {
	return getParser(contentType).Marshal(v)
}
