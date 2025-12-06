package util

import (
	"strconv"
	"strings"
)

// SplitString 分割字符串
func SplitString(str string, sep string) []string {
	// 默认通过换行符去分割
	if sep == "" {
		return strings.Split(str, "\n")
	}
	return strings.Split(str, sep)
}

// ArrayToString 将数组转化为字符串用逗号分割
func ArrayToString(arr []int64) string {
	if len(arr) == 0 {
		return ""
	}
	var ans strings.Builder
	for i := 0; i < len(arr); i++ {
		if i > 0 {
			ans.WriteByte(',')
		}
		ans.WriteString(strconv.FormatInt(arr[i], 10))
	}
	return ans.String()
}
