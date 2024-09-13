// Copyright © 2023 Linbaozhong. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/linbaozhong/model-gen/pkg/types"
	"math"
	"strconv"
	"strings"
	"time"
)

var (
	JSON = jsoniter.ConfigCompatibleWithStandardLibrary
)

func Interface2String(s interface{}) string {
	switch v := s.(type) {
	case string:
		return v
	case types.BigUint:
		return v.String()
	case types.Money:
		return strconv.FormatUint(uint64(v), 10)
	case []byte:
		return string(v)
	case uint64:
		return strconv.FormatUint(v, 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int:
		return strconv.FormatInt(int64(v), 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case time.Time:
		if s.(time.Time).IsZero() {
			return ""
		}
		return s.(time.Time).Format(time.DateTime)
	case bool:
		b := s.(bool)
		if b {
			return "1"
		}
		return "0"
	default:
		return fmt.Sprintf("%+v", s)
	}
}

// String2Bool 字符串转bool
func String2Bool(s string, def ...bool) bool {
	if b, e := strconv.ParseBool(strings.TrimSpace(s)); e == nil {
		return b
	}
	if len(def) > 0 {
		return def[0]
	}
	return false
}

// String2Int8Ptr 字符串转int8指针
func String2Int8Ptr(s string, def ...int64) *int8 {
	intValue := int8(String2Int64(s, def...))
	return &intValue
}

// String2IntPtr 字符串转int指针
func String2IntPtr(s string, def ...int64) *int {
	intValue := int(String2Int64(s, def...))
	return &intValue
}

// String2Int32Ptr 字符串转int32指针
func String2Int32Ptr(s string, def ...int64) *int32 {
	intValue := int32(String2Int64(s, def...))
	return &intValue
}

// String2Int64Ptr 字符串转int64指针
func String2Int64Ptr(s string, def ...int64) *int64 {
	intValue := String2Int64(s, def...)
	return &intValue
}

func String2Uint64(s string, def ...uint64) uint64 {
	if i, e := strconv.ParseUint(strings.TrimSpace(s), 10, 64); e == nil {
		return i
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// String2Uint32 字符串转uint32
func String2Uint32(s string, def ...uint32) uint32 {
	if i, e := strconv.ParseUint(strings.TrimSpace(s), 10, 32); e == nil {
		return uint32(i)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// String2Uint 字符串转uint
func String2Uint(s string, def ...uint) uint {
	if i, e := strconv.ParseUint(strings.TrimSpace(s), 10, 64); e == nil {
		return uint(i)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func String2Uint16(s string, def ...uint16) uint16 {
	if i, e := strconv.ParseUint(strings.TrimSpace(s), 10, 16); e == nil {
		return uint16(i)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// String2Uint8 字符串转uint
func String2Uint8(s string, def ...uint8) uint8 {
	if i, e := strconv.ParseUint(strings.TrimSpace(s), 10, 8); e == nil {
		return uint8(i)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// String2Int
func String2Int(s string, def ...int) int {
	if i, e := strconv.Atoi(strings.TrimSpace(s)); e == nil {
		return i
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// String2Int8 字符串转int8
func String2Int8(s string, def ...int8) int8 {
	if b, e := strconv.ParseInt(strings.TrimSpace(s), 10, 8); e == nil {
		return int8(b)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// String2Int16 字符串转int16
func String2Int16(s string, def ...int16) int16 {
	if b, e := strconv.ParseInt(strings.TrimSpace(s), 10, 16); e == nil {
		return int16(b)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// String2Int32 字符串转int32
func String2Int32(s string, def ...int32) int32 {
	if b, e := strconv.ParseInt(strings.TrimSpace(s), 10, 32); e == nil {
		return int32(b)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// String2Int64 字符串转int64
func String2Int64(s string, def ...int64) int64 {
	if b, e := strconv.ParseInt(strings.TrimSpace(s), 10, 64); e == nil {
		return b
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// String2Float32 字符串转float32
func String2Float32(s string, def ...float32) float32 {
	if b, e := strconv.ParseFloat(strings.TrimSpace(s), 32); e == nil {
		return float32(b)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// String2Float64 字符串转float64
func String2Float64(s string, def ...float64) float64 {
	if b, e := strconv.ParseFloat(strings.TrimSpace(s), 64); e == nil {
		return b
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// String2Time 如果转换失败,返回 def时间(如果存在)
func String2Time(s string, def ...time.Time) time.Time {
	if b, e := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local); e == nil {
		return b
	}
	if len(def) > 0 {
		return def[0]
	}
	return time.Time{}
}

func Interface2Time(s interface{}, def ...time.Time) time.Time {
	if b, ok := s.(time.Time); ok {
		return b
	}
	if len(def) > 0 {
		return def[0]
	}
	return time.Time{}
}

func Interface2Int(s interface{}, def ...int) int {
	switch v := s.(type) {
	case uint64:
		return int(v)
	case uint32:
		return int(v)
	case uint16:
		return int(v)
	case uint8:
		return int(v)
	case uint:
		return int(v)
	case int64:
		return int(v)
	case int32:
		return int(v)
	case int16:
		return int(v)
	case int8:
		return int(v)
	case int:
		return v
	case string:
		return String2Int(v, def...)
	default:
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
}

func Interface2Uint8(s interface{}, def ...uint8) uint8 {
	switch v := s.(type) {
	case uint64:
		return uint8(v)
	case uint32:
		return uint8(v)
	case uint16:
		return uint8(v)
	case uint8:
		return v
	case uint:
		return uint8(v)
	case int64:
		return uint8(v)
	case int32:
		return uint8(v)
	case int16:
		return uint8(v)
	case int8:
		return uint8(v)
	case int:
		return uint8(v)
	case string:
		return String2Uint8(v, def...)
	default:
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
}
func Interface2Uint(s interface{}, def ...uint) uint {
	switch v := s.(type) {
	case types.BigUint:
		return v.Uint()
	case uint64:
		return uint(v)
	case uint32:
		return uint(v)
	case uint16:
		return uint(v)
	case uint8:
		return uint(v)
	case uint:
		return v
	case int64:
		return uint(v)
	case int32:
		return uint(v)
	case int16:
		return uint(v)
	case int8:
		return uint(v)
	case int:
		return uint(v)
	case string:
		return String2Uint(v, def...)
	default:
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
}

func Interface2Uint32(s interface{}, def ...uint32) uint32 {
	switch v := s.(type) {
	case types.BigUint:
		return uint32(v)
	case uint64:
		return uint32(v)
	case uint32:
		return v
	case uint16:
		return uint32(v)
	case uint8:
		return uint32(v)
	case uint:
		return uint32(v)
	case int64:
		return uint32(v)
	case int32:
		return uint32(v)
	case int16:
		return uint32(v)
	case int8:
		return uint32(v)
	case int:
		return uint32(v)
	case string:
		return String2Uint32(v, def...)
	default:
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
}

func Interface2Int64(s interface{}, def ...int64) int64 {
	if b, ok := s.(int64); ok {
		return b
	}
	switch v := s.(type) {
	case types.BigUint:
		return int64(v)
	case uint64:
		return int64(v)
	case uint32:
		return int64(v)
	case uint16:
		return int64(v)
	case uint8:
		return int64(v)
	case uint:
		return int64(v)
	case int32:
		return int64(v)
	case int16:
		return int64(v)
	case int8:
		return int64(v)
	case int:
		return int64(v)
	case string:
		return String2Int64(v, def...)
	default:
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
}

func Interface2Uint64(s interface{}, def ...uint64) uint64 {
	if v, ok := s.(types.BigUint); ok {
		return v.Uint64()
	}
	switch v := s.(type) {
	case uint64:
		return v
	case uint32:
		return uint64(v)
	case uint16:
		return uint64(v)
	case uint8:
		return uint64(v)
	case uint:
		return uint64(v)
	case int64:
		return uint64(v)
	case int32:
		return uint64(v)
	case int16:
		return uint64(v)
	case int8:
		return uint64(v)
	case int:
		return uint64(v)
	case []byte:
		return String2Uint64(string(v), def...)
	case string:
		return String2Uint64(v, def...)
	default:
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
}

func Interface2BigInt(s interface{}, def ...types.BigUint) types.BigUint {
	if v, ok := s.(types.BigUint); ok {
		return v
	}
	switch v := s.(type) {
	case uint64:
		return types.BigUint(v)
	case uint32:
		return types.BigUint(v)
	case uint16:
		return types.BigUint(v)
	case uint8:
		return types.BigUint(v)
	case uint:
		return types.BigUint(v)
	case int64:
		return types.BigUint(v)
	case int32:
		return types.BigUint(v)
	case int16:
		return types.BigUint(v)
	case int8:
		return types.BigUint(v)
	case int:
		return types.BigUint(v)
	case string:
		return types.BigUint(String2Uint64(v))
	default:
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
}

func Interface2Int8(s interface{}, def ...int8) int8 {
	switch v := s.(type) {
	case uint64:
		return int8(v)
	case uint32:
		return int8(v)
	case uint16:
		return int8(v)
	case uint8:
		return int8(v)
	case uint:
		return int8(v)
	case int64:
		return int8(v)
	case int32:
		return int8(v)
	case int16:
		return int8(v)
	case int8:
		return v
	case int:
		return int8(v)
	case string:
		return String2Int8(v, def...)
	default:
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
}
func Interface2Int16(s interface{}, def ...int16) int16 {
	if b, ok := s.(int16); ok {
		return b
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// Interface2StringSlice interface{}转[]string
func Interface2StringSlice(s interface{}, def ...[]string) []string {
	if b, ok := s.([]string); ok {
		return b
	}
	if len(def) > 0 {
		return def[0]
	}
	return []string{}
}

// Interface2IntSlice interface{}转[]int
func Interface2IntSlice(s interface{}, def ...[]int) []int {
	if b, ok := s.([]int); ok {
		return b
	}
	if len(def) > 0 {
		return def[0]
	}
	return []int{}
}

// Interface2StringMap interface{}转map[string]interface{}
func Interface2StringMap(s interface{}, def ...map[string]interface{}) map[string]interface{} {
	if b, ok := s.(map[string]interface{}); ok {
		return b
	}
	if len(def) > 0 {
		return def[0]
	}
	return map[string]interface{}{}
}

// IntToFloat64 IntToFloat64
func IntToFloat64(i int) float64 {
	intValueString := strconv.Itoa(i)
	value, err := strconv.ParseFloat(intValueString, 64)
	if err != nil {
		return 0
	}
	return value
}

func Uint8ToString(i uint8) string {
	return strconv.FormatUint(uint64(i), 10)
}

// Int16ToString Int16ToString
func Int16ToString(i int16) string {
	valueString := strconv.FormatInt(int64(i), 10)
	return valueString
}

// Int32ToString Int32ToString
func Int32ToString(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

// UintToString
func UintToString(i uint) string {
	return strconv.FormatUint(uint64(i), 10)
}

// Uint16ToString
func Uint16ToString(i uint16) string {
	return strconv.FormatUint(uint64(i), 10)
}

// Uint32ToString
func Uint32ToString(i uint32) string {
	return strconv.FormatUint(uint64(i), 10)
}

// Uint64ToString
func Uint64ToString(i uint64) string {
	return strconv.FormatUint(i, 10)
}

// IntToString IntToString
func IntToString(i int) string {
	valueString := strconv.Itoa(i)
	return valueString
}

// Int64ToString Int64ToString
func Int64ToString(i int64) string {
	valueString := strconv.FormatInt(i, 10)
	return valueString
}

// Int8ToString Int8ToString
func Int8ToString(i int8) string {
	valueString := strconv.FormatInt(int64(i), 10)
	return valueString
}

// Float32ToString Float32ToString
func Float32ToString(i float32) string {
	valueString := strconv.FormatFloat(float64(i), 'f', -1, 32)
	return valueString
}

// Float64ToString Float64ToString
func Float64ToString(i float64) string {
	valueString := strconv.FormatFloat(i, 'f', -1, 64)
	return valueString
}

// Wrap 将float64转成精确的int64
func Wrap(num float64, retain int) int64 {
	return int64(num * math.Pow10(retain))
}

// Unwrap 将int64恢复成正常的float64
func Unwrap(num int64, retain int) float64 {
	return float64(num) / math.Pow10(retain)
}

// WrapToFloat64 精准float64
func WrapToFloat64(num float64, retain int) float64 {
	return num * math.Pow10(retain)
}

// UnwrapToInt64 精准int64
func UnwrapToInt64(num int64, retain int) int64 {
	return int64(Unwrap(num, retain))
}

// // 处理float64精度,保留n位小数
// func Round(f float64, n int) float64 {
//	n10 := math.Pow10(n)
//	return math.Trunc((f+0.5/n10)*n10) / n10
// }

// 处理float64精度,保留n位小数,可以是负数
func Round(f float64, n int) float64 {
	n10 := math.Pow10(n)

	if f > 0 {
		return math.Trunc((f+0.5/n10)*n10) / n10
	}
	return math.Trunc((f-0.5/n10)*n10) / n10
}

// Hex2Dec 十六进制转十进制
func Hex2Dec(val string) uint64 {
	n, err := strconv.ParseUint(val, 16, 32)
	if err != nil {
		fmt.Println(err)
	}
	return n
}
