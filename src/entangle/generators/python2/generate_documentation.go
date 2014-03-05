package python2

import (
	"fmt"
	"entangle/utils"
	"strings"
)

// Write documentation.
//
// Documentation is written as individual lines followed by an empty line.
func writeDocumentation(w *safeBuffer, docs []string, indent int) {
	if docs == nil || len(docs) == 0 {
		return
	}

	prefix := strings.Repeat("    ", indent)
	wrapper := utils.NewSimpleTextWrapper(79 - len(prefix))
	lines := make([]string, 0, len(docs)*2)

	for _, paragraph := range docs {
		if len(lines) > 0 {
			lines = append(lines, "")
		} else {
			paragraph = fmt.Sprintf(`"""%s`, strings.TrimSpace(paragraph))
		}

		for _, line := range wrapper.Wrap(paragraph) {
			lines = append(lines, fmt.Sprintf("%s%s", prefix, line))
		}
	}

	if len(lines) == 0 {
		return
	}

	lines = append(lines, fmt.Sprintf(`%s"""`, prefix))

	w.Write([]byte(fmt.Sprintf("%s\n\n", strings.Join(lines, "\n"))))
}
