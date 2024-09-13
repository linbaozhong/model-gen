package mask

import (
	"libs/utils/validation"
	"strings"
)

//脱敏手机号
func MaskMobile(m string) string {
	if validation.Mobile(m) {
		return m[:3] + "****" + m[7:]
	}
	return m
}

//脱敏邮箱
func MaskEmail(m string) string {
	if validation.Email(m) {
		pos := strings.Index(m, "@")
		return string([]rune(m)[:1]) + "****" + m[pos-1:]
	}
	return m
}

func MaskName(m string) string {
	ms := []rune(m)
	if len(ms) > 1 {
		return string(ms[:1]) + "**"
	}
	return m
}
