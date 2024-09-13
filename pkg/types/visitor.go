package types

import jsoniter "github.com/json-iterator/go"

var (
	JSON = jsoniter.ConfigCompatibleWithStandardLibrary
)

//Visitor 当前操作者-客户
type Visitor struct {
	Peer    string
	From    string //产品的访问端名称：比如：手机端=mobile，小程序=mp，web端=web
	ID      BigUint
	Name    string
	AppID   BigUint //子账号id
	Own     BigUint //所有者id，B端数据：companyID，C端数据：userID
	IP      string
	AppType int8 // 应用程序所属
}
