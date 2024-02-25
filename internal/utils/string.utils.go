package utils

import "strings"

func SplitByFirstString(str string, splitter string) (string, string) {
	spt := strings.Split(str, splitter)

	if len(spt) > 1 {
		return spt[0], strings.Join(spt[1:], splitter)
	} else {
		return spt[0], ""
	}
}
