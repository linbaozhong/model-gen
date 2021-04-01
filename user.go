package main

//tablename users
type User struct {
	ID int64 `json:"id" xorm:"id"`
}
