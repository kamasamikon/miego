package filter

import (
	"regexp"
	"strings"
)

// untouched when error
func filterList(arr []string, pattern string) []string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return arr
	}

	result := make([]string, 0, len(arr))
	for _, s := range arr {
		if re.MatchString(s) {
			result = append(result, s)
		}
	}
	return result
}

func filter(text string, pattern string) string {
	return strings.Join(filterList(strings.Split(text, "\n"), pattern), "\n")
}
