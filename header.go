package isuperagent

import (
	"regexp"
	"strings"
)

type ContentType struct {
	MediaType string
	Charset   string
	Boundary  string
}

func ParseContentType(raw string) ContentType {
	c := ContentType{
		MediaType: "text/plain",
		Charset:   "utf-8",
	}
	exp := regexp.MustCompile(`\s*;\s*`)
	contentType := exp.Split(strings.TrimSpace(raw), -1)

	if len(contentType) > 0 && strings.TrimSpace(contentType[0]) != "" {
		c.MediaType = contentType[0]
	}

	if len(contentType) > 1 && strings.TrimSpace(contentType[1]) != "" {
		c.Charset = contentType[1]
	}

	if len(contentType) > 2 && strings.TrimSpace(contentType[2]) != "" {
		c.Charset = contentType[2]
	}

	return c
}

func (c *ContentType) String() string {
	var buf strings.Builder

	buf.WriteString(c.MediaType)

	if c.Charset != "" {
		buf.WriteString("; ")
		buf.WriteString(c.Charset)
	}

	if c.Boundary != "" {
		buf.WriteString("; ")
		buf.WriteString(c.Boundary)
	}

	return buf.String()
}
