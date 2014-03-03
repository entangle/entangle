package commands

import (
	"fmt"
	"entangle/utils"
	"strings"
)

type usageListElement struct {
	Name string
	Synopsis string
}

// Generate usage list.
func usageList(elements []usageListElement) string {
	// Determine the longest name.
	maxNameLength := 0

	for _, elem := range elements {
		if len(elem.Name) > maxNameLength {
			maxNameLength = len(elem.Name)
		}
	}

	// Determine the synopsis indentation.
	indentation := ((2 + maxNameLength + 5) / 4) * 4
	wrapper := utils.NewSimpleTextWrapper(79 - indentation)
	result := make([]string, 0, len(elements) * 2)
	subsequentLinePrefix := strings.Repeat(" ", indentation)

	for _, elem := range elements {
		for i, line := range wrapper.Wrap(elem.Synopsis) {
			if i == 0 {
				result = append(result, fmt.Sprintf("  %s%s%s", elem.Name, strings.Repeat(" ", indentation - 2 - len(elem.Name)), line))
			} else {
				result = append(result, fmt.Sprintf("%s%s", subsequentLinePrefix, line))
			}
		}
	}

	return strings.Join(result, "\n")
}
