package models

import (
	"github.com/linbaozhong/model-gen/models/table"
	"sync"
)

// TableName log_reg
type LogReg struct {
	ID         int64  `gorm:"column:id;PRIMARY_KEY"`
	Mobile     int64  `gorm:"column:mobile"`
	Nick       string `gorm:"column:nick"`
	InviteCode string `gorm:"column:invite_code"`
	InviteUser int64  `gorm:"column:invite_user"`
	IP         int64  `gorm:"column:ip"`
	Device     int8   `gorm:"column:device"`
	DeviceID   string `gorm:"column:deviceid"`
}

var (
	logregPool = sync.Pool{
		New: func() interface{} {
			return &LogReg{}
		},
	}
)

func NewLogReg() *LogReg {
	return logregPool.Get().(*LogReg)
}

func (p *LogReg) Free() {
	//todo:初始化每个字段

	p.Device = ""

	p.DeviceID = ""

	p.ID = ""

	p.IP = ""

	p.InviteCode = ""

	p.InviteUser = ""

	p.Mobile = ""

	p.Nick = ""

	logregPool.Put(p)
}

func (*LogReg) TableName() string {
	return table.LogReg.TableName
}
