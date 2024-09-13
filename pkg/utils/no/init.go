package no

import (
	"strconv"

	"github.com/bwmarrin/snowflake"
	"golang.org/x/sync/singleflight"
)

type nodeKey int64

func (k nodeKey) String() string {
	return strconv.FormatInt(int64(k), 10)
}

const (
	No_Request_ID nodeKey = iota //用户
	No_Resume                    //简历号
	No_Project                   //项目号
	No_Job                       //工作职位

	No_B_Sign_Document  //
	No_B_Transaction_Id //

	No_A_Sign_Document //

	No_C_Transaction_Id //

	No_Staffing_Approval_No
	No_Staffing_Transaction_Id
	No_Staffing_Sign_Document_Id

	No_Biz_Transaction_Id
	No_Biz_Company_Man_Id
	No_Biz_Sign_Document_Id
	No_Biz_Company_Code_Id
)

var (
	generateNodes = make(map[nodeKey]*snowflake.Node)
	sg            singleflight.Group
)

func GetGenerateNode(key nodeKey) *snowflake.Node {
	if node, ok := generateNodes[key]; ok && node != nil {
		return node
	}
	r, e, _ := sg.Do(key.String(), func() (r interface{}, e error) {
		r, e = snowflake.NewNode(int64(key))
		if e == nil {
			generateNodes[key] = r.(*snowflake.Node)
		} else {
			generateNodes[key] = nil
		}
		return
	})
	if e != nil {
		panic(e)
	}
	return r.(*snowflake.Node)
}
