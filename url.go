package isuperagent

import (
	"net/url"
)

type URL struct {
	*url.URL
	Queries url.Values
}

func NewURL() *URL {
	return &URL{
		URL:     &url.URL{},
		Queries: url.Values{},
	}
}

func (u *URL) String() string {
	u.RawQuery = u.Queries.Encode()

	return u.String()
}
