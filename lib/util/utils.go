package util

import (
	"hound/lib/proto"
	"io"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func ParseUrl(u *url.URL) *proto.UrlType {
	nu := &proto.UrlType{}
	nu.Scheme = u.Scheme
	nu.Domain = u.Hostname()
	nu.Host = u.Host
	nu.Port = u.Port()
	nu.Path = u.EscapedPath()
	nu.Query = u.RawQuery
	nu.Fragment = u.Fragment
	return nu
}

func GetNowDateTime() string {
	now := time.Now()
	return now.Format("01-02 15:04:05")
}

func GetNowDateTimeReportName() string {
	now := time.Now()
	return now.Format("20060102-150405")
}

func IsURL(input string) bool {
	u, err := url.ParseRequestURI(input)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func UrlTypeToString(u *proto.UrlType) string {
	var buf strings.Builder
	if u.Scheme != "" {
		buf.WriteString(u.Scheme)
		buf.WriteByte(':')
		buf.WriteString("//")
	}
	if h := u.Host; h != "" {
		buf.WriteString(u.Host)
	}
	if u.Path != "" && u.Path[0] != '/' && u.Host != "" {
		buf.WriteByte('/')
	}
	buf.WriteString(u.Path)
	if u.Query != "" {
		buf.WriteByte('?')
		buf.WriteString(u.Query)
	}
	if u.Fragment != "" {
		buf.WriteByte('#')
		buf.WriteString(u.Fragment)
	}
	return buf.String()
}

func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// 字符串转 utf 8
func Str2UTF8(str string) string {
	if len(str) == 0 {
		return ""
	}
	if !utf8.ValidString(str) {
		utf8Bytes, err := io.ReadAll(transform.NewReader(
			strings.NewReader(str),
			simplifiedchinese.GBK.NewDecoder(),
		))
		if err != nil {
			// 处理错误
			return ""
		}
		return string(utf8Bytes)
	}
	return str
}
