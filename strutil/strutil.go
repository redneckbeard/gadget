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

// IsPascalCase returns true if its argument s is a Pascal-cased string.
func IsPascalCase(s string) bool {
	return pascalCase.MatchString(s)
}

// Hyphenate converts Pascal-cased string to a hyphen-separated string.
func Hyphenate(pascal string) string {
	return depascal(pascal, "-")
}

// Snakify converts a Pascal-cased string to a snake-cased string.
func Snakify(pascal string) string {
	return depascal(pascal, "_")
}
