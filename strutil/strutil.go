package strutil

import (
	"regexp"
	"strings"
)

var pascalCase *regexp.Regexp

func init() {
	pascalCase = regexp.MustCompile(`[A-Z]+[a-z\d]+`)
}

func Hyphenate(pascal string) string {
	matches := []string{}
	for _, match := range pascalCase.FindAllString(pascal, -1) {
		matches = append(matches, strings.ToLower(match))
	}
	return strings.Join(matches, "-")
}
