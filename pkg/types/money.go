package types

import (
	"regexp"
	"strconv"
	"strings"
)

type Money int64

func (m *Money) MarshalJSON() ([]byte, error) {
	yuan := strconv.FormatFloat((float64(*m) / 100), 'f', -1, 64)
	return []byte(yuan), nil
}

func (m *Money) UnmarshalJSON(b []byte) error {
	fen, e := strconv.ParseFloat(string(b), 10)
	if e != nil {
		return e
	}
	i, e := strconv.ParseInt(strconv.FormatFloat(fen*100, 'f', 0, 64), 10, 64)
	if e == nil {
		*m = Money(i)
	}
	return e
}

// Yuan 金额分精确到元
func (m *Money) Yuan() float64 {
	return float64(*m) / 100
}

func (m *Money) ToCNY() string {
	if *m == 0 {
		return "零元整"
	}
	if *m < 0 {
		var mm Money = *m * -1
		return "负" + mm.ToCNY()
	}

	numstr := []rune(strconv.Itoa(int(*m)))
	numlen := len(numstr)
	moneyUnit := []string{"仟", "佰", "拾", "亿", "仟", "佰", "拾", "万", "仟", "佰", "拾", "元", "角", "分"}
	unit := moneyUnit[len(moneyUnit)-numlen:]
	num := map[rune]string{48: "零", 49: "壹", 50: "贰", 51: "叁", 52: "肆", 53: "伍", 54: "陆", 55: "柒", 56: "捌", 57: "玖"}

	var hasZero bool
	var buf strings.Builder
	for i := 0; i < numlen; i++ {
		if numstr[i] == 48 {
			if strings.Index("亿万", unit[i]) > -1 {
				buf.WriteString(unit[i])
			}
			if hasZero {
				continue
			}
			hasZero = true
		} else {
			if hasZero {
				buf.WriteString("零")
			}
			hasZero = false
			buf.WriteString(num[numstr[i]] + unit[i])
		}
	}

	result := buf.String()
	if strings.HasSuffix(result, "元") || strings.HasSuffix(result, "角") {
		buf.WriteString("整")
	} else if strings.HasSuffix(result, "分") {
	} else {
		buf.WriteString("元整")
	}

	result = strings.Replace(buf.String(), "亿万", "亿", -1)
	return result
}

// 金额分小写转中文大写金额
func (m *Money) ConvertNumToCny() (string, error) {
	strnum := strconv.Itoa(int(*m))

	sliceUnit := []string{"仟", "佰", "拾", "亿", "仟", "佰", "拾", "万", "仟", "佰", "拾", "元", "角", "分"}

	s := sliceUnit[len(sliceUnit)-len(strnum) : len(sliceUnit)]
	upperDigitUnit := map[string]string{"0": "零", "1": "壹", "2": "贰", "3": "叁", "4": "肆", "5": "伍", "6": "陆", "7": "柒", "8": "捌", "9": "玖"}
	str := ""
	for k, v := range strnum[:] {
		str = str + upperDigitUnit[string(v)] + s[k]
	}
	reg, err := regexp.Compile(`零角零分$`)
	str = reg.ReplaceAllString(str, "整")

	reg, err = regexp.Compile(`零角`)
	str = reg.ReplaceAllString(str, "零")

	reg, err = regexp.Compile(`零分$`)
	str = reg.ReplaceAllString(str, "整")

	reg, err = regexp.Compile(`零[仟佰拾]`)
	str = reg.ReplaceAllString(str, "零")

	reg, err = regexp.Compile(`零{2,}`)
	str = reg.ReplaceAllString(str, "零")

	reg, err = regexp.Compile(`零亿`)
	str = reg.ReplaceAllString(str, "亿")

	reg, err = regexp.Compile(`零万`)
	str = reg.ReplaceAllString(str, "万")

	reg, err = regexp.Compile(`零*元`)
	str = reg.ReplaceAllString(str, "元")

	reg, err = regexp.Compile(`亿零{0, 3}万`)
	str = reg.ReplaceAllString(str, "^元")

	reg, err = regexp.Compile(`零元`)
	str = reg.ReplaceAllString(str, "零")
	if err != nil {
		return "", err
	}
	return str, nil
}
