package url

import (
	"net"
	"net/url"
	"strings"
)

func IsVal(str string) bool {
	url, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}

	address := net.ParseIP(url.Host)

	if address == nil {
		return strings.Contains(url.Host, ".")
	}

	return true
}
