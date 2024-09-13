// @Title 时间类型转换
// @Description
// @Author 蔺保仲 2020/04/20
// @Update 蔺保仲 2020/04/20
package utils

import (
	"fmt"
	"github.com/linbaozhong/model-gen/pkg/constant"
	"strconv"
	"strings"
	"time"
)

func GetAgeByBirthday(birth time.Time) int {

	now := time.Now()

	if birth.After(now) {
		return 0
	}

	y := birth.Year()
	m := birth.Month()
	age := now.Year() - y
	if now.Month() < m {
		age--
	}
	return age
}

// 2006.01 格式生日的年龄
func GetAgeByBirthdayString(birth string) string {

	if birth == "" || len(birth) < 7 {
		return ""
	}

	birth = strings.Replace(birth, "-", ".", -1)
	birth = strings.Replace(birth, "/", ".", -1)

	t, _ := time.ParseInLocation("2006.01", birth[:7], time.Local)

	age := GetAgeByBirthday(t)

	return IFF(age == 0, "", IntToString(age)).(string)
}

// TimeToMilliSecond 时间转毫秒
func TimeToMilliSecond(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.UnixNano() / 1e6
}

// NowMilliSecond 当前时间的毫秒值
func NowMilliSecond() int64 {
	return TimeToMilliSecond(time.Now())
}

// MsecToTime 毫秒转时间
func MsecToTime(msec int64) time.Time {
	if msec == 0 {
		return time.Unix(0, 0)
	}
	t := msec / 1000
	return time.Unix(t, msec*1e6%t)
}

// MsecToTimeString 毫秒转日期字符串
func MsecToTimeString(msec int64) string {
	return TimeToString(MsecToTime(msec))
}

// TimeToString 时间转字符串
func TimeToString(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(constant.Time_Layout)
}

// SecToTime 秒转时间
func SecToTime(sec int64) time.Time {
	if sec == 0 {
		return time.Unix(0, 0)
	}
	t := sec
	return time.Unix(t, 0)
}

// ParseDateStrToUnixTime 日期格式转换为时间戳（秒
// @author xu.sun
// @example
//
//	ParseDateStrToUnixTime("2020-05-03 23:59:59")
func ParseDateStrToUnixTime(dateStr string) int64 {
	loc, _ := time.LoadLocation("Local")
	theTime, err := time.ParseInLocation(constant.Time_Layout, dateStr, loc)
	if err == nil {
		return theTime.Unix()
	}
	return 0
}

// ParseDateStrToTime 日期格式转换为时间格式
// @author xu.sun
// @example
//
//	ParseDateStrToTime("2020-05-03 23:59:59")
func ParseDateStrToTime(dateStr string) time.Time {
	unix := ParseDateStrToUnixTime(dateStr)
	return time.Unix(unix, 0)
}

// ParseUnixTimeToDateStr 时间戳转换为日期格式
// @author xu.sun
// @example
//
//	ParseUnixTimeToDateStr(1588521599)
//	ParseUnixTimeToDateStr("1588521599")
//	ParseUnixTimeToDateStr("1588521599", "2006-01-02")
func ParseUnixTimeToDateStr(ptime interface{}, layout ...string) string {
	var timeInt64 int64 = 0
	if ptimeStr, ok := ptime.(string); ok {
		timeInt64 = String2Int64(ptimeStr)
	} else if ptimeInt64, ok := ptime.(int64); ok {
		timeInt64 = ptimeInt64
	}

	var cusLayout = ""
	if len(layout) > 0 {
		cusLayout = layout[0]
	} else {
		cusLayout = constant.Time_Layout
	}
	return time.Unix(timeInt64, 0).Format(cusLayout)
}

// MinuteToString 日期分钟转字符串
func MinuteToString(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(constant.Time_Layout_Ignore_Seconds)
}

// MonthToString 年月转字符串
func MonthToString(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(constant.Date_Layout_Ignore_day)
}

// NumberToString 年月日转字符串，不带符号
func NumberToString(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(constant.Date_Layout_Ignore)
}

// DateToString 日期转字符串
func DateToString(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(constant.Date_Layout)
}

// DateVersionToString 日期转版本号日期字符串
func DateVersionToString(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(constant.Date_Layout_Version)
}

// DateToStringChinese 日期转字符串
func DateToStringChinese(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(constant.Date_Layout_Chinese)
}

// 合同日期 2006年01月02日
func DateContract(t time.Time) string {
	if t.IsZero() {
		return "/"
	}
	return t.Format(constant.Date_Layout_Chinese)
}

// HourAndSecondsToString 时间分钟转字符串
func HourAndSecondsToString(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(constant.Time_Layout_Minute)
}

// DateToStringPeriod 日期转字符串,.分割
func DateToStringPeriod(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(constant.Date_Layout_Period)
}

// DateAndWeekToString 日期星期转字符串
func DateAndWeekToString(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	date := t.Format(constant.Date_Layout_Chinese)

	switch t.Weekday() {
	case time.Sunday:
		date += " 周日"
	case time.Monday:
		date += " 周一"
	case time.Tuesday:
		date += " 周二"
	case time.Wednesday:
		date += " 周三"
	case time.Thursday:
		date += " 周四"
	case time.Friday:
		date += " 周五"
	case time.Saturday:
		date += " 周六"
	default:
		date += ""
	}
	return date
}

// BeginMonth 获得下个月初的日期
func NextBeginMonth(t time.Time) time.Time {
	month := t.Month()
	year := t.Year()
	if month == 12 {
		year += 1
		month = 1
		return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
	}

	month += 1
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())

}

// 获得月初的日期
func BeginMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())

}

// 获得某日零点时间
func BeginDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// 天数增加
func AddDay(currentDay time.Time, days int64) time.Time {
	return currentDay.AddDate(0, 0, int(days))
}

// 根据出生日期得到年龄字符串
func GetAgeByTime(t time.Time) string {

	if t.Year() <= 1 {
		return ""
	}

	age := time.Now().Year() - t.Year()

	if age <= 0 {
		return ""
	}
	return strconv.Itoa(age)
}

// 根据出生日期得到年龄
func GetAgeByTimeInt(t time.Time) int {

	if t.Year() <= 1 {
		return 0
	}

	age := time.Now().Year() - t.Year()

	if age < 0 {
		return 0
	}
	return age
}

// 根据参加工作时间得到工作年数月数
func GetWorkingByTime(t time.Time) string {

	if t.Year() <= 1 {
		return "0年"
	}

	nowTime := time.Now()
	sub := nowTime.Sub(t)
	dt := time.Time{}.Add(sub)
	y := dt.Year() - 1
	m := int(dt.Month()) - 1
	var f strings.Builder
	if y > 0 {
		f.WriteString(IntToString(y) + "年 ")
	}
	if m > 0 {
		f.WriteString(IntToString(m) + "个月")
	}
	return f.String()

	// year := nowTime.Year() - t.Year()
	//
	// month := nowTime.Month() - t.Month()
	//
	// if month < 0 {
	//	year -= 1
	//
	//	month += 12
	// }
	//
	// if year < 0 {
	//	return "0年"
	// }
	//
	// if year == 0 && month == 0 {
	//	return "0年"
	// }
	//
	// if year == 0 {
	//	return fmt.Sprintf("%d个月", month)
	// }
	// if month == 0 {
	//	return fmt.Sprintf("%d年", year)
	// }
	//
	// return fmt.Sprintf("%d年%d个月", year, month)
}

// 获得特定时间到当前时间相隔的月数、天数和小时数,即是距离时间t又过了多久
func GetIntervalByTime(t time.Time) string {
	if t.IsZero() {
		return "未知开始时间"
	}
	nowTime := time.Now()
	if t.After(nowTime) {
		return "0天"
	}

	sub := nowTime.Sub(t)
	dt := time.Time{}.Add(sub)
	y := dt.Year() - 1
	m := int(dt.Month()) - 1
	d := dt.Day() - 1
	h := dt.Hour()

	var f strings.Builder
	if y > 0 {
		f.WriteString(IntToString(y) + "年 ")
	}
	if m > 0 {
		f.WriteString(IntToString(m) + "个月 ")
	}
	if d > 0 {
		f.WriteString(IntToString(d) + "天 ")
	}
	if h > 0 {
		f.WriteString(IntToString(h) + "小时")
	}
	return f.String()

	// year := nowTime.Year() - t.Year()
	//
	// month := int(nowTime.Month() - t.Month())
	//
	// day := nowTime.Day() - t.Day()
	//
	// hour := nowTime.Hour() - t.Hour()
	//
	// if hour < 0 {
	//	day -= 1
	//	hour += 24
	// }
	//
	// if day < 0 {
	//	month -= 1
	//
	//	var days int
	//	switch t.Month() {
	//
	//	case time.January, time.March, time.May, time.July, time.August, time.October, time.December:
	//		days = 31
	//	case time.February:
	//		if (t.Year()%4 == 0 && t.Year()%100 != 0) || t.Year()%400 == 0 {
	//			days = 29
	//		} else {
	//			days = 28
	//		}
	//	case time.April, time.June, time.September, time.November:
	//		days = 30
	//	default:
	//		days = 30
	//	}
	//
	//	day += days
	// }
	//
	// if month < 0 {
	//	year -= 1
	//	month += 12
	// }
	//
	// month = int(month) + year*12
	//
	// if month == 0 && day == 0 && hour == 0 {
	//	return "0小时"
	// }
	//
	// if month == 0 && day == 0 {
	//	return fmt.Sprintf("%d小时", hour)
	// }
	//
	// if month == 0 {
	//	return fmt.Sprintf("%d天 %d小时", day, hour)
	// }
	// return fmt.Sprintf("%d个月 %d天 %d小时", month, day, hour)
}

// WeekToString 星期转字符串
func WeekToString(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	switch t.Weekday() {
	case time.Sunday:
		return "星期日"
	case time.Monday:
		return "星期一"
	case time.Tuesday:
		return "星期二"
	case time.Wednesday:
		return "星期三"
	case time.Thursday:
		return "星期四"
	case time.Friday:
		return "星期五"
	case time.Saturday:
		return "星期六"

	}
	return ""
}

// 获得当前时间到特定时间相隔的月数、天数和小时数,即是当前时间距离时间t有多久
func GetInterval(t time.Time) string {

	nowTime := time.Now()
	if nowTime.After(t) {
		return "0天"
	}

	sub := t.Sub(nowTime)
	dt := time.Time{}.Add(sub)
	y := dt.Year() - 1
	m := int(dt.Month()) - 1
	d := dt.Day() - 1
	// h := dt.Hour()

	var f strings.Builder
	if y > 0 {
		f.WriteString(IntToString(y) + "年 ")
	}
	if m > 0 {
		f.WriteString(IntToString(m) + "个月 ")
	}
	if d > 0 {
		f.WriteString(IntToString(d) + "天 ")
	}
	// if h > 0 {
	//	f.WriteString(IntToString(h) + "小时")
	// }
	return f.String()

}

// 计算日期相差多少月,t1小时间t2大时间
func SubMonth(t1, t2 time.Time) (month int) {

	if t1.After(t2) {
		return 0
	}

	sub := t2.Sub(t1)
	dt := time.Time{}.Add(sub)

	fmt.Println(dt)
	y := dt.Year() - 1
	m := int(dt.Month()) - 1
	day := dt.Day() - 1

	month = y*12 + m

	if day > 15 {
		month++
	}
	return month
}
