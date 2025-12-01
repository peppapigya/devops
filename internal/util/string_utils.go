package util

import "strings"

// SplitString 分割字符串
func SplitString(str string, sep string) []string {
	// 默认通过换行符去分割
	if sep == "" {
		return strings.Split(str, "\n")
	}
	return strings.Split(str, sep)
}
