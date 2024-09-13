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
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"time"
)

type CompareFunc func(interface{}, interface{}) bool

func IndexOf(a []interface{}, e interface{}, cmp CompareFunc) int {
	n := len(a)
	for i := 0; i < n; i++ {
		if cmp(Interface2String(e), Interface2String(a[i])) {
			return i
		}
	}
	return -1
}

// 参数e是否包含在切片a中
func Contains(a []interface{}, e interface{}) bool {
	return IndexOf(a, e, func(i interface{}, i2 interface{}) bool {
		return i == i2
	}) != -1
}

func Clone(src interface{}) interface{} {
	var dst reflect.Value
	var get_dst = func(typ reflect.Type) reflect.Value {
		dst := reflect.New(typ).Elem()            // 创建对象
		b, _ := JSON.Marshal(src)                 // 导出json
		JSON.Unmarshal(b, dst.Addr().Interface()) // json序列化
		return dst
	}
	typ := reflect.TypeOf(src)
	if typ.Kind() == reflect.Ptr { // 如果是指针类型
		typ = typ.Elem() // 获取源实际类型(否则为指针类型)
		dst = get_dst(typ)
		return dst.Addr().Interface() // 返回指针
	}
	dst = get_dst(typ)
	return dst.Interface() // 返回值
}

// IsZero 是否零值
func IsZero(a interface{}) bool {
	switch v := a.(type) {
	case string:
		return v == ""
	case int, int64, int32, int16, int8, uint64, uint32, uint16, uint8:
		return v == 0
	case float64, float32:
		return v == 0.
	case bool:
		return v == false
	case time.Time:
		return v.IsZero()
	default:
		return v == nil
	}
}

// slice是否相等
func IsEqualSlice(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	if (a == nil) != (b == nil) {
		return false
	}

	b = b[:len(a)]
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

func GetAppPath() string {
	root, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return root
}

// CompressString 压缩字符串（空格、换行、回车、制表符）
func CompressString(s string) string {
	if s == "" {
		return s
	}
	reg := regexp.MustCompile("\\s+")
	return reg.ReplaceAllString(s, "")
}

func CompressXML(s string) string {
	if s == "" {
		return s
	}
	reg := regexp.MustCompile("[\\t\\n\\f\\r]+")
	return reg.ReplaceAllString(s, "")
}

// IFF 三元运算
func IFF(b bool, v1, v2 interface{}) interface{} {
	if b {
		return v1
	}
	return v2
}
