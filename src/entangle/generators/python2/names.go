package python2

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	camelCaseExpressionA = regexp.MustCompile(`[A-Z]+[A-Z][a-z]`)
	camelCaseExpressionB = regexp.MustCompile(`[a-z\d][A-Z]`)
)

// Snake case a string.
func snakeCaseString(name string) string {
	name = camelCaseExpressionA.ReplaceAllStringFunc(name, func(x string) string {
		return fmt.Sprintf("%s_%s", x[:len(x) - 2], x[len(x) - 2:])
	})

	name = camelCaseExpressionB.ReplaceAllStringFunc(name, func(x string) string {
		return fmt.Sprintf("%s_%s", x[:len(x) - 1], x[len(x) - 1:])
	})

	return strings.ToLower(name)
}

// Snake case a slice of strings.
func snakeCaseStrings(names []string) (result []string) {
	result = make([]string, len(names))
	for i, n := range names {
		result[i] = snakeCaseString(n)
	}
	return
}
