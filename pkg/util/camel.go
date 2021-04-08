package util

import "strings"

func Camel(str string) string {
	name := strings.Replace(str, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}
