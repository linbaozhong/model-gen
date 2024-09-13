package utils

import (
	"crypto/sha1"
	"fmt"
)

// Signature sha1签名
func Signature(s []byte) string {
	return fmt.Sprintf("%x", SignatureByte(s))
}

func SignatureByte(s []byte) []byte {
	sha := sha1.New()
	sha.Write(s)
	return sha.Sum(nil)
}
