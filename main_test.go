package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestString(t *testing.T) {
	fmt.Println([]byte("az_"))
	str := "TableName"
	fmt.Println(getFieldName(str))
}

func getFieldName(name string) string {
	bs := bytes.NewBuffer([]byte{})
	for i, s := range name {
		if s >= 65 && s <= 90 {
			s += 32
			if i == 0 {
				bs.WriteByte(byte(s))
			} else {
				bs.WriteByte(byte(95))
				bs.WriteByte(byte(s))
			}
			continue
		}
		bs.WriteByte(byte(s))
	}
	return bs.String()
}
