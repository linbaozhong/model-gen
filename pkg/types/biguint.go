package types

import (
	"bytes"
	"strconv"
)

type BigUint uint64

func (i BigUint) MarshalJSON() ([]byte, error) {

	//if i == 0 {
	//	return []byte(""), nil
	//}

	return []byte(i.String()), nil
}

//func (i BigUint) MarshalText() ([]byte, error) {
//
//	//if i == 0 {
//	//	return []byte(""), nil
//	//}
//
//	return []byte(i.String()), nil
//}

func (i *BigUint) UnmarshalJSON(b []byte) error {
	c := string(bytes.Trim(b, "\""))

	if c == "" {
		*i = BigUint(0)
		return nil
	}
	tem, e := strconv.ParseUint(c, 10, 64)

	*i = BigUint(tem)
	return e
}

//String
func (i BigUint) String() string {
	return strconv.FormatUint(uint64(i), 10)
}

//Uint64
func (i BigUint) Uint64() uint64 {
	return uint64(i)
}

//Uint
func (i BigUint) Uint() uint {
	return uint(i)
}
