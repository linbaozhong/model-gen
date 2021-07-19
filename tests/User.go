package tests

//tablename user
type User struct {
	Name   string `json:"name"`
	Depart string `json:"depart"`
	Age    int    `json:"age"`
	IS     bool   `json:"is"`
}
