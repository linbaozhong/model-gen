// @Title MD5加解密
// @Description
// @Author 蔺保仲 2020/04/20
// @Update 蔺保仲 2020/04/20
package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func _md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

//MD5String 对字符串进行md5加密
func MD5String(s string, salt string) string {
	return _md5(s + salt)
}
