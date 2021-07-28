package tests

//tablename user
type User struct {
	ID     uint64 `json:"id" xorm:"'id' pk"`
	Name   string `json:"name" xorm:"'name'"`
	Depart string `json:"depart"`
	Age    int    `json:"age" xorm:"'age'"`
	IS     bool   `json:"is"`
}
