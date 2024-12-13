package util

import (
	"regexp"
	"strings"
)

var (
	ReNameReg = regexp.MustCompile(`(?m)(\s?[(（【][^)）】(（【册卷套]{4,}[)）】])`)
)

func CleanTitle(title string) string {
	if len(ReNameReg.FindAllString(title, -1)) == 0 {
		return title
	}
	newTitle := ReNameReg.ReplaceAllString(title, "")
	newTitle = strings.TrimSpace(strings.ReplaceAll(newTitle, "\"", " "))
	return newTitle
}
