package util

import "strings"

// RemoveIllegalChar 去除字符串中的空格、回车和换行。
func RemoveIllegalChar(s string) string {
	s = strings.ReplaceAll(s, " ", "")  // 去除空格
	s = strings.ReplaceAll(s, "\r", "") // 去除回车
	s = strings.ReplaceAll(s, "\n", "") // 去除换行
	return s
}
