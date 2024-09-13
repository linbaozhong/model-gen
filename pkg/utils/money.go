package utils

import (
	"strings"
)

//YuanString2Fen 金额元字符串 转 分
func YuanString2Fen(s string) int64 {
	s += "00"
	pos := strings.IndexByte(s, '.')
	if pos == 0 { //第0位是。
		return String2Int64(s[1:3])
	} else if pos > 0 {
		return String2Int64(s[:pos] + s[pos+1:pos+3])
	}
	return String2Int64(s)
}
