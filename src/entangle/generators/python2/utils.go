package python2

import (
	"fmt"
)

// Stringify a slice of strings.
//
// Turns [a b c] into ['a' 'b' 'c'].
func stringifyStrings(input []string) (output []string) {
	output = make([]string, len(input))
	for i, in := range input {
		output[i] = fmt.Sprintf("'%s'", in)
	}
	return
}
