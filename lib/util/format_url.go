package util

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	HTTP             = "http"
	HTTPS            = "https"
	SchemeSeparator  = "://"
	DefaultHTTPPort  = "80"
	DefaultHTTPSPort = "443"
)

func Hostname(s string) (string, error) {
	if !strings.HasPrefix(s, HTTP) && strings.HasPrefix(s, HTTPS) {
		s = HTTP + SchemeSeparator + s
	}
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	return u.Hostname(), nil
}

func Host(s string) (string, error) {
	if !strings.HasPrefix(s, HTTP) && strings.HasPrefix(s, HTTPS) {
		s = HTTP + SchemeSeparator + s
	}
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	return u.Host, nil
}

// url编码
func URLEncodeAllChar(s string) string {
	b := []byte(s)
	nb := ""
	for _, v := range b {
		nb += fmt.Sprintf("%%%X", v)
	}
	return nb
}
