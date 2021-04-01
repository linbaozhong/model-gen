package main

import (
	"sync"

	"github.com/linbaozhong/model-gen/table"
)

var (
	userPool = sync.Pool{
		New: func() interface{} {
			return &User{}
		},
	}
)

func NewUser() *User {
	return userPool.Get().(*User)
}

func (p *User) Free() {
	p.ID = 0

	userPool.Put(p)
}

func (*User) TableName() string {
	return table.User.TableName
}

//func (p *User) ToMap() map[string]interface{} {
//	m := make(map[string]interface{}, 1)
//	m[table.User.ID.Name] = p.ID
//
//	return m
//}
