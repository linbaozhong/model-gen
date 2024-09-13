package utils

import (
	"fmt"
	"math"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

func TestGetInviteCode(t *testing.T) {
	var id uint64 = 14329658718059274240
	code, _ := GetInviteCode(id)
	_id := GetIDFromInviteCode(code)
	t.Log(id, code, _id, id == _id)
}

func TestBirthday(t *testing.T) {

	birthTime := time.Date(2023, 10, 3, 0, 0, 0, 0, time.Local)
	span := birthTime.Sub(time.Now())
	fmt.Println(birthTime, math.Floor(span.Hours()/24))
}

func TestXORM(t *testing.T) {
	s := GetDbConnectionString("ssld_dev", "Cu83&sr66")
	x, e := xorm.NewEngine("mysql", s)
	if e != nil {
		panic(e)
	}

	x.ShowSQL(true)

	x.SetMaxOpenConns(100)
	x.SetMaxIdleConns(5)

	if e := x.Ping(); e != nil {
		panic(e)
	}

	cls := make([]interface{}, 0)
	e = x.Table("companys_mans_info").Cols("id").Where("not isnull(hire_date)").Limit(1000).Find(&cls)
	if e != nil {
		t.Log(e)
	}
	fmt.Println(cls)

	//cls = make([]interface{}, 0)
	//e = x.Table("companys_mans_info").Cols("id").Where("id_card <> ?", "").Limit(1000).Find(&cls)
	//if e != nil {
	//	t.Log(e)
	//}
	//fmt.Println(cls)
}

func GetDbConnectionString(user, pwd string) string {

	return user + ":" + pwd + "@" + "tcp(39.107.252.66:13306)" + "/" + "hr_biz" + "?" + "charset=utf8mb4&parseTime=True&loc=Local"

}
