// @Title 雪花算法生成乱序数字id,支持1024个节点
// @Description
// @Seller
// @Update
package utils

import (
	"golang.org/x/sync/singleflight"

	"github.com/bwmarrin/snowflake"
)

const (
	No_User    = "user"    //用户
	No_Resume  = "resume"  //简历号
	No_Project = "project" //项目号
	No_Job     = "job"     //工作职位
)

var (
	nodes = map[string]int64{
		No_User:    0, //用户
		No_Resume:  1, //简历号
		No_Project: 2, //项目号
		No_Job:     3, //工作职位
	}
	generateNodes = make(map[string]*snowflake.Node)
	sg            singleflight.Group
)

func getGenerateNode(key string) *snowflake.Node {
	if node, ok := generateNodes[key]; ok && node != nil {
		return node
	}
	r, e, _ := sg.Do(key, func() (r interface{}, e error) {
		r, e = snowflake.NewNode(nodes[key])
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

//GetUserNo 读取一个新的用户编号
func GetUserNo() int64 {
	node := getGenerateNode(No_User)
	return node.Generate().Int64()
}

//GetResumeNo 读取一个新的简历编号
func GetResumeNo() int64 {
	node := getGenerateNode(No_Resume)
	return node.Generate().Int64()
}

//GetProjectNo 生成一个新的项目编号
func GetProjectNo() int64 {
	node := getGenerateNode(No_Project)
	return node.Generate().Int64()
}

//GetJobNo 生成新的职位编号
func GetJobNo() int64 {
	node := getGenerateNode(No_Job)
	return node.Generate().Int64()
}
