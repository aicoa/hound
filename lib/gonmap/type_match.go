package gonmap

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type match struct {
	service       string
	pattern       string
	versionInfo   *FingerPrint
	patternRegexp *regexp.Regexp
	soft          bool
}

var matchVersionInfoHelperRegxP = regexp.MustCompile(`\$P\((\d)\)`)

// matchVersionInfoHelperRegxP 匹配形如 $P(1) 这样的模式，其中 \d 是一个数字占位符。

var matchVersionInfoHelperRegxV = regexp.MustCompile(`\$(\d)`)

// matchVersionInfoHelperRegxV 匹配形如 $1 这样的模式，用于从捕获的文本中提取具体值。

var matchLoadRegexps = []*regexp.Regexp{
	regexp.MustCompile("^([a-zA-Z0-9-_./]+) m\\|([^|]+)\\|([is]{0,2})(?: (.*))?$"),
	// 匹配形如 "service m|pattern|options (version info)" 的字符串。

	regexp.MustCompile("^([a-zA-Z0-9-_./]+) m=([^=]+)=([is]{0,2})(?: (.*))?$"),
	// 匹配形如 "service m=pattern=options (version info)" 的字符串。

	regexp.MustCompile("^([a-zA-Z0-9-_./]+) m%([^%]+)%([is]{0,2})(?: (.*))?$"),
	// 匹配形如 "service m%pattern%options (version info)" 的字符串。

	regexp.MustCompile("^([a-zA-Z0-9-_./]+) m@([^@]+)@([is]{0,2})(?: (.*))?$"),
	// 匹配形如 "service m@pattern@options (version info)" 的字符串。
}

var matchVersionInfoRegexps = map[string]*regexp.Regexp{
	"PRODUCTNAME": regexp.MustCompile("p/([^/]+)/"),
	// PRODUCTNAME 匹配产品名，形如 "p/product-name/" 的字符串。

	"VERSION": regexp.MustCompile("v/([^/]+)/"),
	// VERSION 匹配版本号，形如 "v/version-number/" 的字符串。

	"INFO": regexp.MustCompile("i/([^/]+)/"),
	// INFO 匹配附加信息，形如 "i/additional-info/" 的字符串。

	"HOSTNAME": regexp.MustCompile("h/([^/]+)/"),
	// HOSTNAME 匹配主机名，形如 "h/hostname/" 的字符串。

	"OS": regexp.MustCompile("o/([^/]+)/"),
	// OS 匹配操作系统，形如 "o/operating-system/" 的字符串。

	"DEVICE": regexp.MustCompile("d/([^/]+)/"),
	// DEVICE 匹配设备类型，形如 "d/device-type/" 的字符串。
}

var matchVersionInfoHelperRegx = regexp.MustCompile(`\$(\d)`)

func parseMatch(s string, soft bool) *match {
	m := &match{soft: soft}
	var regx *regexp.Regexp

	for _, r := range matchLoadRegexps {
		if r.MatchString(s) {
			regx = r
		}
	}

	if regx == nil {
		panic(errors.New("match 语句参数不正确"))
	}

	args := regx.FindStringSubmatch(s)
	//m.service = FixProtocol(args[1])
	m.service = FixProtocol(args[1])
	m.pattern = args[2]
	m.patternRegexp = m.getPatternRegexp(m.pattern, args[3])
	m.versionInfo = parseVersionInfo(s)

	return m
}

func (this *match) getPatternRegexp(pattern string, options string) *regexp.Regexp {
	pattern = strings.ReplaceAll(pattern, `\0`, `\x00`)
	if options != "" {
		if strings.Contains(options, "i") == false {
			options += "i"
		}
		if pattern[:1] == "^" {
			pattern = fmt.Sprintf("^(?%s:%s", options, pattern[1:])
		} else {
			pattern = fmt.Sprintf("(?%s:%s", options, pattern)
		}
		if pattern[len(pattern)-1:] == "$" {
			pattern = fmt.Sprintf("%s)$", pattern[:len(pattern)-1])
		} else {
			pattern = fmt.Sprintf("%s)", pattern)
		}
	}
	return regexp.MustCompile(pattern)
}

func (this *match) getVersionInfo(s string, regId string) string {
	if matchVersionInfoRegexps[regId].MatchString(s) {
		return matchVersionInfoRegexps[regId].FindStringSubmatch(s)[1]
	} else {
		return ""
	}
}
func (m *match) makeVersionInfo(s string, f *FingerPrint) {
	f.Info = m.makeVersionInfoSubHelper(s, m.versionInfo.Info)
	f.DeviceType = m.makeVersionInfoSubHelper(s, m.versionInfo.DeviceType)
	f.Hostname = m.makeVersionInfoSubHelper(s, m.versionInfo.Hostname)
	f.Os = m.makeVersionInfoSubHelper(s, m.versionInfo.Os)
	f.Version = m.makeVersionInfoSubHelper(s, m.versionInfo.Version)
	f.Service = m.makeVersionInfoSubHelper(s, m.versionInfo.Service)
}

func (this *match) makeVersionInfoSubHelper(s string, pattern string) string {
	if len(this.patternRegexp.FindStringSubmatch(s)) == 1 {
		return pattern
	}
	if pattern == "" {
		return pattern
	}
	sArr := this.patternRegexp.FindStringSubmatch(s)

	if matchVersionInfoHelperRegxP.MatchString(pattern) {
		pattern = matchVersionInfoHelperRegxP.ReplaceAllStringFunc(pattern, func(s string) string {
			a := matchVersionInfoHelperRegxP.FindStringSubmatch(s)[1]
			return "$" + a
		})
	}
	if matchVersionInfoHelperRegx.MatchString(pattern) {
		pattern = matchVersionInfoHelperRegx.ReplaceAllStringFunc(pattern, func(s string) string {
			a, _ := strconv.Atoi(matchVersionInfoHelperRegx.FindStringSubmatch(s)[1])
			return sArr[a]
		})
	}
	pattern = strings.ReplaceAll(pattern, "\n", "")
	pattern = strings.ReplaceAll(pattern, "\r", "")
	return pattern

}
