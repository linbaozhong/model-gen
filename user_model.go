package main

type User struct {
	ID       int64 `json:"id" grom:"column:id;primary_key"`
	Name     string
	Age      int8
	NickName string
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
	p.Age = 0
	p.ID = 0
	p.Name = ""
	p.NickName = ""

	userPool.Put(p)
}

func (*User) TableName() string {
	return table.User.TableName
}
