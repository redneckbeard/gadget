package strutil

import (
	"regexp"
	"strings"
)

var pascalCase, pascalSegment *regexp.Regexp

func init() {
	pascalCase = regexp.MustCompile(`^([A-Z]+[a-z\d]+)+$`)
	pascalSegment = regexp.MustCompile(`[A-Z]+[a-z\d]+`)
}

func depascal(pascal, separator string) string {
	matches := []string{}
	for _, match := range pascalSegment.FindAllString(pascal, -1) {
		matches = append(matches, strings.ToLower(match))
	}
	return strings.Join(matches, separator)
}

func IsPascalCase(s string) bool {
	return pascalCase.MatchString(s)
}

func Hyphenate(pascal string) string {
	return depascal(pascal, "-")
}

func Snakify(pascal string) string {
	return depascal(pascal, "_")
}
