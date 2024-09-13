// @Title Aes加密解密
// @Description
// @Author 蔺保仲 2020/04/20
// @Update 蔺保仲 2020/04/20
package utils

import (
	"encoding/base64"

	crypter "github.com/sekrat/aescrypter"
)

func AesEncrypt(src, key string) (string, error) {
	crypter := crypter.New()
	buf, err := crypter.Encrypt(key, []byte(src))
	return base64.URLEncoding.EncodeToString(buf), err
}

func AesDecrypt(src, key string) (string, error) {
	buf, err := base64.URLEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}
	crypter := crypter.New()
	buf, err = crypter.Decrypt(key, buf)
	return string(buf), err
}

func AesEncryptBytes(src []byte, key string) ([]byte, error) {
	crypter := crypter.New()
	return crypter.Encrypt(key, src)
}

func AesDecryptBytes(src []byte, key string) ([]byte, error) {
	crypter := crypter.New()
	return crypter.Decrypt(key, src)
}
