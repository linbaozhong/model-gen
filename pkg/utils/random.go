package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.NewSource(int64(time.Now().Nanosecond()))
}

// GetRandString 生成随机字符串
func GetRandString(l int) string {
	chars := "ABCDEFGHIJKMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz23456789"
	charsLen := len(chars)
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = chars[rand.Intn(charsLen)]
	}
	return string(bytes)
}

//GetRandDigit 生成范围随机数字
func GetRandDigit(min, max int64) int64 {
	if min >= max || min == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}
