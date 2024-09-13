package utils

import (
	"strconv"
	"strings"
)

func Uint64Slice2Pieces(ids []uint64, size int) [][]uint64 {
	pieces := make([][]uint64, 0, len(ids)/size+1)
	for i := 0; i < len(ids); i += size {
		if len(ids) > i+size {
			pieces = append(pieces, ids[i:i+size])
		} else {
			pieces = append(pieces, ids[i:])
		}
	}
	return pieces
}

func IntSlice2String(src []int) string {
	var buf = make([]string, 0, len(src))
	for i := 0; i < len(src); i++ {
		buf = append(buf, strconv.Itoa(src[i]))
	}
	return strings.Join(buf, ",")
}
