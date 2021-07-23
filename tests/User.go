package tests

//tablename user
//cachedata time.Minute *  1
//  cachelist time.Minute * 1
//cachelimit   500
type User struct {
	Name   string `json:"name" xorm:"'name'"`
	Depart string `json:"depart"`
	Age    int    `json:"age" xorm:"'age'"`
	IS     bool   `json:"is"`
}
