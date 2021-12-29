package helper

import "strings"

func BbyteToSString(b [][]byte) []string {
	res := make([]string, 0)
	for i := 0; i < len(b); i++ {
		res = append(res, string(b[i]))
	}
	return res
}

func ContainWithoutCaseSensitive(s []string, e string) int {
	// 这里不要使用 range， 会造成 string 拷贝
	for i := 0; i < len(s); i++ {
		if strings.ToLower(e) == strings.ToLower(s[i]) {
			return i
		}
	}
	return -1
}
