package utils

import "strings"

func GetStringArrayFromSqlString(value string) []string {
	return strings.Split(strings.Trim(value, "{}"), ",")
}
