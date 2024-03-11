package util

import "regexp"

var (
	// RegIP 匹配标准IPv4地址
	RegIP = regexp.MustCompile(`((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}`)

	// RegDomain 匹配合法域名
	RegDomain = regexp.MustCompile(`[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+\.?`)

	// RegCompliance 匹配合规字符串;匹配空间引擎输入内容是否合规
	RegCompliance = regexp.MustCompile(`(\w+)[!,=]{1,3}"([^"]+)"`)

	// RegIPCompleteMask 匹配IPv4地址和完整子网掩码
	RegIPCompleteMask = regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}/\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)

	// RegIPCIDR 匹配CIDR表示的IP地址范围
	RegIPCIDR = regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})/(\d{1,2})`)
)
