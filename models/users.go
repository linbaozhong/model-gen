package models

import (
	"time"
)

//tablename users
type Users struct {
	ID     uint64    `json:"id" xorm:"'id' pk autoincr"` //id
	Mobile uint64    `json:"mobile" xorm:"'mobile'"`     //手机号码
	Email  string    `json:"email" xorm:"'email'"`       //email
	Nick   string    `json:"nick" xorm:"'nick'"`         //昵称
	Pwd    string    `json:"pwd" xorm:"'pwd'"`           //密码
	Ctime  time.Time `json:"ctime" xorm:"'ctime'"`
	IP     uint64    `json:"ip" xorm:"'ip' ->"` //ip地址
}

//tablename wallet
type Wallet struct {
	ID       uint64 //用户id
	Currency uint8  `json:"currency" xorm:"'currency'"` //币种:
	Amount   uint64 `json:"amount" xorm:"'amount'"`     //金额:单位:分
	Fee      uint64 `json:"fee" xorm:"'fee'"`           //服务费,单位:分
	JSON     string `json:"json" xorm:"'jn'"`
}
