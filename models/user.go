package models

import (
	"github.com/linbaozhong/model-gen/models/table"
	"sync"
)

//TableName   user
type User struct {
	ID         int64  `json:"id" gorm:"column:id;PRIMARY_KEY"`
	Mobile     int64  `json:"mobile" gorm:"column:mobile"`
	Pwd        string `json:"pwd" gorm:"column:pwd"`
	Inviter    int64  `json:"inviter" gorm:"column:inviter"`
	InviteCode string `json:"invite_code" gorm:"column:invite_code"`
	Ukey       string `json:"ukey" gorm:"column:ukey"`
	UUID       string `json:"uuid" gorm:"column:uuid"`
}

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
	//todo:初始化每个字段
	p.ID = 0
	p.InviteCode = ""
	p.Inviter = 0
	p.Mobile = 0
	p.Pwd = ""
	p.UUID = ""
	p.Ukey = ""

	userPool.Put(p)
}

func (*User) TableName() string {
	return table.User.TableName
}
