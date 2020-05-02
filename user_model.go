package main

type User struct {
	ID       int64 `json:"id" grom:"column:id;primary_key"`
	Name     string
	Age      int8
	NickName string
}
