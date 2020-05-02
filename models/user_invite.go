package models

import (
	"github.com/linbaozhong/model-gen/models/table"
	"sync"
)

//TableName user_invite
type UserInvite struct {
	ID         int64  `json:"id" gorm:"column:id"`                   //用户id
	Inviter    int64  `json:"inviter" gorm:"column:inviter"`         //邀请者id
	InviteCode string `json:"invite_code" gorm:"column:invite_code"` //邀请码
	Depth      int8   `json:"depth" gorm:"column:depth"`             //链深度
	Chains50   string `json:"chains50" gorm:"column:chains50"`       //邀请链50节点
}

var (
	userinvitePool = sync.Pool{
		New: func() interface{} {
			return &UserInvite{}
		},
	}
)

func NewUserInvite() *UserInvite {
	return userinvitePool.Get().(*UserInvite)
}

func (p *UserInvite) Free() {
	//todo:初始化每个字段

	p.Chains50 = ""

	p.Depth = ""

	p.ID = ""

	p.InviteCode = ""

	p.Inviter = ""

	userinvitePool.Put(p)
}

func (*UserInvite) TableName() string {
	return table.UserInvite.TableName
}
