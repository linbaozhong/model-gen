package models

import "sync"

// TableName black_mobile
type BlackMobile struct {
	ID     int64 `gorm:"column:id;PRIMARY_KEY"`
	Mobile int64 `gorm:"column:mobile"`
}

var (
	blackmobilePool = sync.Pool{
		New: func() interface{} {
			return &BlackMobile{}
		},
	}
)

func NewBlackMobile() *BlackMobile {
	return blackmobilePool.Get().(*BlackMobile)
}

func (p *BlackMobile) Free() {
	//todo:初始化每个字段

	p.ID = ""

	p.Mobile = ""

	blackmobilePool.Put(p)
}

func (*BlackMobile) TableName() string {
	return table.BlackMobile.TableName
}
