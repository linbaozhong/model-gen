package utils

// 求两个切片交集
func Intersect1(s1 []string, s2 []string) []string {

	result := make([]string, 0)

	if len(s1) == 0 || len(s2) == 0 {
		return result
	}

	tempMap := make(map[string]bool)
	for _, v := range s2 {
		tempMap[v] = true
	}

	for _, v := range s1 {
		if c, ok := tempMap[v]; ok && c {
			result = append(result, v)
			tempMap[v] = false
		}
	}
	return result
}

// 求两个切片差集
func Difference(s1 []string, s2 []string) []string {

	result := make([]string, 0)

	if len(s1) == 0 {
		return result
	}

	if len(s2) == 0 {
		return s1
	}

	tempMap := make(map[string]bool)
	for _, v := range s2 {
		tempMap[v] = true
	}

	for _, v := range s1 {
		if _, ok := tempMap[v]; !ok {
			result = append(result, v)
		}
	}
	return result
}
