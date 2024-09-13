package constant

import (
	"errors"
	"math"
	"time"
)

var (
	Request_Time_Key_NotFound = time.Time{} //key未找到
	Err_Type                  = errors.New("类型错误")
)

const (
	Time_Layout                = "2006-01-02 15:04:05"
	Time_Layout_Ignore_Seconds = "2006-01-02 15:04"
	Time_Layout_Minute         = "15:04"
	Date_Layout                = "2006-01-02"
	Date_Layout_Chinese        = "2006年01月02日"
	Date_Layout_Ignore_day     = "2006-01"
	Date_Layout_Version        = "2006-0102"
	Date_Layout_Ignore         = "20060102"
	Date_Layout_Period         = "2006.01.02"

	Request_String_Key_NotFound  = "NIL"           //key未找到
	Request_Uint64_Key_NotFound  = math.MaxUint64  //key未找到
	Request_Uint32_Key_NotFound  = math.MaxUint32  //key未找到
	Request_Uint16_Key_NotFound  = math.MaxUint16  //key未找到
	Request_Uint8_Key_NotFound   = math.MaxUint8   //key未找到
	Request_Int64_Key_NotFound   = math.MinInt64   //key未找到
	Request_Int32_Key_NotFound   = math.MinInt32   //key未找到
	Request_Int16_Key_NotFound   = math.MinInt16   //key未找到
	Request_Int8_Key_NotFound    = math.MinInt8    //key未找到
	Request_Float64_Key_NotFound = math.MaxFloat64 //key未找到
	Request_Float32_Key_NotFound = math.MaxFloat32 //key未找到

	Access_Cache_Key  = "access"
	Refresh_Cache_Key = "refresh"

	B_My_Rules  = "b_my_rules"
	B_Man_Rules = "b_man_rules"

	// 合付宝手续费 单位百分号
	Hefubao_Rate float64 = 4.9
)
