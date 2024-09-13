package request

import (
	"html"
	"libs/constant"
	"libs/types"
	"libs/utils"
	"strconv"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
)

func getRequestValue(c iris.Context, k string) (string, bool) {
	vals := c.FormValues()
	if v, ok := vals[k]; ok {
		if len(v) > 0 {
			return v[0], true
		}
		return "", true
	}
	return "", false
}

func PostString(c iris.Context, k string, def ...string) string {
	s, b := getRequestValue(c, k)
	if b {
		return html.EscapeString(s)
	}
	if len(def) > 0 {
		return def[0]
	}
	return constant.Request_String_Key_NotFound
}

func PostMoney(c iris.Context, k string, def ...types.Money) types.Money {
	s, b := getRequestValue(c, k)
	if b {
		return types.Money(utils.YuanString2Fen(s))
	}
	if len(def) > 0 {
		return def[0]
	}
	return constant.Request_Int64_Key_NotFound
}

func PostFloat64(c iris.Context, k string, def ...float64) float64 {
	s, b := getRequestValue(c, k)
	if b {
		return utils.String2Float64(s)
	}
	if len(def) > 0 {
		return def[0]
	}
	return constant.Request_Float64_Key_NotFound
}

func PostFloat32(c iris.Context, k string, def ...float32) float32 {
	s, b := getRequestValue(c, k)
	if b {
		return float32(utils.String2Float64(s))
	}
	if len(def) > 0 {
		return def[0]
	}
	return constant.Request_Float32_Key_NotFound
}

func PostInt64(c iris.Context, k string, def ...int64) int64 {
	s, b := getRequestValue(c, k)
	if b {
		return utils.String2Int64(s)
	}
	if len(def) > 0 {
		return def[0]
	}
	return constant.Request_Int64_Key_NotFound
}

func PostInt16(c iris.Context, k string, def ...int16) int16 {
	s, b := getRequestValue(c, k)
	if b {
		return utils.String2Int16(s)
	}
	if len(def) > 0 {
		return def[0]
	}
	return constant.Request_Int16_Key_NotFound
}

func PostInt8(c iris.Context, k string, def ...int8) int8 {
	s, b := getRequestValue(c, k)
	if b {
		return utils.String2Int8(s)
	}
	if len(def) > 0 {
		return def[0]
	}
	return constant.Request_Int8_Key_NotFound
}

func PostInt(c iris.Context, k string, def ...int) int {
	s, b := getRequestValue(c, k)
	if b {
		return utils.String2Int(s)
	}
	if len(def) > 0 {
		return def[0]
	}
	return constant.Request_Int64_Key_NotFound

}

func PostBigUint(c iris.Context, k string, def ...types.BigUint) types.BigUint {
	s, b := getRequestValue(c, k)
	if b {
		return types.BigUint(utils.String2Uint64(s))
	}
	if len(def) > 0 {
		return def[0]
	}
	return constant.Request_Uint64_Key_NotFound
}

func PostUint64(c iris.Context, k string, def ...uint64) uint64 {
	s, b := getRequestValue(c, k)
	if b {
		return utils.String2Uint64(s)
	}
	if len(def) > 0 {
		return def[0]
	}
	return constant.Request_Uint64_Key_NotFound
}

func PostUint(c iris.Context, k string, def ...uint) uint {
	s, b := getRequestValue(c, k)
	if b {
		return utils.String2Uint(s)
	}
	if len(def) > 0 {
		return def[0]
	}
	return constant.Request_Uint64_Key_NotFound

}

func PostUint32(c iris.Context, k string, def ...uint32) uint32 {
	s, b := getRequestValue(c, k)
	if b {
		return utils.String2Uint32(s)
	}
	if len(def) > 0 {
		return def[0]
	}
	return constant.Request_Uint32_Key_NotFound
}

func PostUint16(c iris.Context, k string, def ...uint16) uint16 {
	s, b := getRequestValue(c, k)
	if b {
		return utils.String2Uint16(s)
	}
	if len(def) > 0 {
		return def[0]
	}
	return constant.Request_Uint16_Key_NotFound
}

func PostUint8(c iris.Context, k string, def ...uint8) uint8 {
	s, b := getRequestValue(c, k)
	if b {
		return utils.String2Uint8(s)
	}
	if len(def) > 0 {
		return def[0]
	}
	return constant.Request_Uint8_Key_NotFound
}

func PostTime(c iris.Context, k string, def ...time.Time) time.Time {
	s, b := getRequestValue(c, k)
	if b {
		if t, e := time.ParseInLocation(constant.Time_Layout, s, time.Local); e == nil {
			return t
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return constant.Request_Time_Key_NotFound
}

func PostDate(c iris.Context, k string, def ...time.Time) time.Time {
	s, b := getRequestValue(c, k)
	if b {
		if t, e := time.ParseInLocation(constant.Date_Layout, s, time.Local); e == nil {
			return t
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return constant.Request_Time_Key_NotFound
}

// ////////////////////////////////
func getParamString(c iris.Context, k string) (string, bool) {
	en, b := c.Params().Store.GetEntry(k)
	if b {
		return en.StringTrim(), true
	}
	return "", false
}

func ParamString(c iris.Context, k string, def ...string) string {
	v, b := getParamString(c, k)
	if b {
		return html.EscapeString(v)
	}

	if len(def) > 0 {
		return def[0]
	}
	return constant.Request_String_Key_NotFound
}

func ParamFloat64(c iris.Context, k string, def ...float64) float64 {
	v, b := getParamString(c, k)
	if b {
		return utils.String2Float64(v, def...)
	}

	if len(def) > 0 {
		return def[0]
	}
	return constant.Request_Float64_Key_NotFound
}

func ParamInt64(c iris.Context, k string, def ...int64) int64 {
	v, b := getParamString(c, k)
	if b {
		return utils.String2Int64(v, def...)
	}

	if len(def) > 0 {
		return def[0]
	}
	return constant.Request_Int64_Key_NotFound
}

func ParamInt(c iris.Context, k string, def ...int) int {
	v, b := getParamString(c, k)
	if b {
		return utils.String2Int(v, def...)
	}

	if len(def) > 0 {
		return def[0]
	}
	return constant.Request_Int64_Key_NotFound
}

func ParamUint64(c iris.Context, k string, def ...uint64) uint64 {
	v, b := getParamString(c, k)
	if b {
		return utils.String2Uint64(v, def...)
	}

	if len(def) > 0 {
		return def[0]
	}
	return constant.Request_Uint64_Key_NotFound
}

func ParamUint(c iris.Context, k string, def ...uint) uint {
	v, b := getParamString(c, k)
	if b {
		return utils.String2Uint(v, def...)
	}

	if len(def) > 0 {
		return def[0]
	}
	return constant.Request_Uint64_Key_NotFound
}

func ParamUint16(c iris.Context, k string, def ...uint16) uint16 {
	v, b := getParamString(c, k)
	if b {
		return utils.String2Uint16(v, def...)
	}

	if len(def) > 0 {
		return def[0]
	}
	return constant.Request_Uint16_Key_NotFound
}

func ParamTime(c iris.Context, k string, def ...time.Time) time.Time {
	v, b := getParamString(c, k)
	if b {
		t, e := time.ParseInLocation(constant.Time_Layout, v, time.Local)
		if e == nil {
			return t
		}
	}

	if len(def) > 0 {
		return def[0]
	}
	return constant.Request_Time_Key_NotFound
}

// /////////////////////
// 切片
func PostInterfaceSlice(c iris.Context, k string) []interface{} {
	_vs := PostStringSlice(c, k)
	vsLen := len(_vs)
	vs := make([]interface{}, 0, vsLen)
	for _, s := range _vs {
		if i, e := strconv.Atoi(s); e == nil {
			vs = append(vs, i)
		}
	}
	return vs
}
func PostIntSlice(c iris.Context, k string) []int {
	_vs := PostStringSlice(c, k)
	vsLen := len(_vs)
	vs := make([]int, 0, vsLen)
	for _, s := range _vs {
		if i, e := strconv.Atoi(s); e == nil {
			vs = append(vs, i)
		}
	}
	return vs
}

func PostInt64Slice(c iris.Context, k string) []int64 {
	_vs := PostStringSlice(c, k)
	vsLen := len(_vs)
	vs := make([]int64, 0, vsLen)
	for _, s := range _vs {
		vs = append(vs, utils.String2Int64(s))
	}
	return vs
}
func PostStringSlice(c iris.Context, k string) []string {
	v := c.PostValue(k)
	if v == "" {
		return []string{}
	}
	return strings.Split(v, ",")
}
