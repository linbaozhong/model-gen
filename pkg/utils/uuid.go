// @Title UUID生成器
// @Description
// @Author 蔺保仲 2020/04/20
// @Update 蔺保仲 2020/04/20
package utils

import (
	"encoding/hex"

	"github.com/google/uuid"
)

//GetUUID 返回去除连接线(-)的32位字符的uuid字符串
func GetUUID() string {
	id := uuid.New()
	buf := make([]byte, 32)
	hex.Encode(buf[:], id[:])
	return string(buf[:])
}
