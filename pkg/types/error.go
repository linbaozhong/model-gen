package types

import (
	"fmt"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Info    string `json:"info"`
}

func (e Error) Error() string {
	if e.Info == "" {
		return e.Message
	}
	return e.Message + ":" + e.Info
}

func (e Error) Is(err error) bool {
	if err == nil {
		return false
	}
	if er, ok := err.(Error); ok {
		return er.Code == e.Code
	}
	return false
}

func (e Error) SetInfo(i interface{}) Error {
	if i == nil {
		return e
	}
	if err, ok := i.(error); ok {
		e.Info = err.Error()
	} else if s, ok := i.(string); ok {
		e.Info = s
	} else {
		e.Info = fmt.Sprint(i)
	}
	return e
}

func NewError(code int, message string) Error {
	return Error{
		Code:    code,
		Message: message,
	}
}
