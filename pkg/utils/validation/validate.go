// @Title Api校验入参格式
// @Description
// @Author 蔺保仲 2020/04/20
// @Update 蔺保仲 2020/04/20
package validation

import (
	"fmt"
	"github.com/linbaozhong/model-gen/pkg/types"
	"github.com/linbaozhong/model-gen/pkg/utils"
	"math"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

// 判定系统是32位还是64位
var wordsize = 32 << (^uint(0) >> 63)

func Required(obj interface{}) bool {
	if obj == nil {
		return false
	}
	switch v := obj.(type) {
	case string:
		return len(strings.TrimSpace(v)) > 0
	case bool:
		return true
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return v != 0
	case time.Time:
		return !v.IsZero()
	case types.BigUint, types.Money:
		return v != 0
	default:
		if v := reflect.ValueOf(obj); v.Kind() == reflect.Slice {
			return v.Len() > 0
		}
	}
	return false
}

// 字母 数字 特殊字符 必有其二
func AlphaDigitChar(obj interface{}) bool {
	return match(obj, alphaDigitChar)
}

var (
	// 邮政编码，有效类型：string
	zipCodePattern = regexp.MustCompile(`^[1-9]\d{5}$`)
	// 固定电话号，有效类型：string
	telPattern = regexp.MustCompile(`^(0\d{2,3}(\-)?)?\d{7,8}$`)
	// 手机号，有效类型：string   手机号更新太快，放开第2和3位限制
	mobilePattern = regexp.MustCompile(`^((\+86)|(86))?(1([0-9][0-9]))\d{8}$`)
	// base64 编码，有效类型：string
	base64Pattern = regexp.MustCompile(`^(?:[A-Za-z0-99+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$`)
	// IP 格式，目前只支持 IPv4 格式验证，有效类型：string
	ipPattern = regexp.MustCompile(`^((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)$`)
	// alpha 字符或数字或横杠 -_，有效类型：string
	alphaDashPattern = regexp.MustCompile(`[^\d\w-_]`)
	// 邮箱格式，有效类型：string
	emailPattern = regexp.MustCompile(`^[\w!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[\w!#$%&'*+/=?^_` + "`" + `{|}~-]+)*@(?:[\w](?:[\w-]*[\w])?\.)+[a-zA-Z0-9](?:[\w-]*[\w])?$`)
	// 字母 数字 特殊字符二选一
	alphaDigitChar = regexp.MustCompile("(\\d+.*[a-zA-Z_]+)|([a-zA-Z_]+.*\\d+)|([\\W_]+.*\\d+)|(\\d+.*\\W+)|([\\W_]+.*[a-zA-Z_]+)|([a-zA-Z_]+.*\\W+)")
)

func AlphaNumeric(obj interface{}) bool {
	if obj == nil {
		return false
	}
	if str, ok := obj.(string); ok {
		for _, v := range str {
			if ('Z' < v || v < 'A') && ('z' < v || v < 'a') && ('9' < v || v < '0') {
				return false
			}
		}
		return true
	}
	return false
}

func Numeric(obj interface{}) bool {
	if obj == nil {
		return false
	}
	if str, ok := obj.(string); ok {
		for _, v := range str {
			if '9' < v || v < '0' {
				return false
			}
		}
		return true
	}
	return false
}

func Alpha(obj interface{}) bool {
	if obj == nil {
		return false
	}
	if str, ok := obj.(string); ok {
		for _, v := range str {
			if ('Z' < v || v < 'A') && ('z' < v || v < 'a') {
				return false
			}
		}
		return true
	}
	return false
}

func MaxSize(obj interface{}, max int) bool {
	if obj == nil {
		return false
	}
	if str, ok := obj.(string); ok {
		return utf8.RuneCountInString(str) <= max
	}

	if v := reflect.ValueOf(obj); v.Kind() == reflect.Slice {
		return v.Len() <= max
	}
	return false
}

func MinSize(obj interface{}, min int) bool {
	if obj == nil {
		return false
	}
	if str, ok := obj.(string); ok {
		return utf8.RuneCountInString(str) >= min
	}

	if v := reflect.ValueOf(obj); v.Kind() == reflect.Slice {
		return v.Len() >= min
	}
	return false
}

func Range(obj interface{}, min, max int) bool {
	return Min(obj, min) && Max(obj, max)
}

func Max(obj interface{}, max int) bool {
	if obj == nil {
		return false
	}

	switch v := obj.(type) {
	case string:
		return utils.String2Int(v, math.MaxInt32) <= max
	case int64:
		if wordsize == 32 {
			return false
		}
		return int(v) <= max
	case int, int32, int16, int8:
		return utils.Interface2Int(v, 0) <= max
	default:
		return false
	}
}

func Min(obj interface{}, min int) bool {
	if obj == nil {
		return false
	}

	switch v := obj.(type) {
	case string:
		return utils.String2Int(v, math.MinInt32) >= min
	case int64:
		if wordsize == 32 {
			return false
		}
		return int(v) >= min
	case int, int32, int16, int8:
		return utils.Interface2Int(v, 0) >= min
	default:
		return false
	}
}

func Length(obj interface{}, n int, max ...int) bool {
	if obj == nil {
		return false
	}

	if str, ok := obj.(string); ok {
		slen := utf8.RuneCountInString(str)
		if len(max) == 1 {
			return slen >= n && slen <= max[0]
		}
		return slen == n
	}

	if v := reflect.ValueOf(obj); v.Kind() == reflect.Slice {
		slen := v.Len()
		if len(max) == 1 {
			return slen >= n && slen <= max[0]
		}
		return slen == n
	}
	return false
}

func ZipCode(obj interface{}) bool {
	return match(obj, zipCodePattern)
}

func Tel(obj interface{}) bool {
	return match(obj, telPattern)
}

func Mobile(obj interface{}) bool {
	return match(obj, mobilePattern)
}

func Base64(obj interface{}) bool {
	return match(obj, base64Pattern)
}

func IP(obj interface{}) bool {
	return match(obj, ipPattern)
}

func Email(obj interface{}) bool {
	return match(obj, emailPattern)
}

func AlphaDash(obj interface{}) bool {
	return match(obj, alphaDashPattern)
}

func match(obj interface{}, pattern *regexp.Regexp) bool {
	if obj == nil {
		return false
	}
	return pattern.MatchString(fmt.Sprintf("%v", obj))
}

func Match(obj interface{}, pattern string) bool {
	if obj == nil {
		return false
	}
	ok, _ := regexp.MatchString(pattern, fmt.Sprintf("%v", obj))
	return ok
}

// id为unit64类型，不能为负数
func IsId(id types.BigUint) bool {
	return id >= 0
}
