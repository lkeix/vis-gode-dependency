package utils

import "strings"

func IsUpperCase(s string) bool {
	return strings.ToUpper(string(s[0])) == string(s[0])
}
