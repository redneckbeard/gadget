package strutil

import (
	"bufio"
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
	scanner := bufio.NewScanner(strings.NewReader(pascal))
	scanner.Split(bufio.ScanBytes)
	var match, lastSeen string
	for scanner.Scan() {
		c := scanner.Text()
		if strings.ToUpper(c) == c && strings.ToLower(lastSeen) == lastSeen && !strings.ContainsAny(c, "1234568790") && lastSeen != "" {
			matches = append(matches, strings.ToLower(match+lastSeen))
			match, lastSeen = "", c
		} else if strings.ToLower(c) == c && strings.ToUpper(lastSeen) == lastSeen && lastSeen != "" && len(match) > 0 && match[len(match)-1:] == strings.ToUpper(match[len(match)-1:]) {
			matches = append(matches, strings.ToLower(match))
			match, lastSeen = lastSeen, c	
		} else {
			match, lastSeen = match+lastSeen, c
		}
	}
	matches = append(matches, strings.ToLower(match+lastSeen))
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
