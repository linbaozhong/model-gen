package utils

import "strings"

// 按符号切割字符串
func SplitString(object, symbol string) []string {
	if object == "" {
		return make([]string, 0)
	} else {
		slice := strings.Split(object, symbol)
		return slice
	}
}
