package python2

import (
	"fmt"
	"entangle/utils"
	"strings"
)

// Code writer.
//
// Utility function for writing code line by line with easy indentation
// management.
type codeWriter struct {
	indent int
	buffer *safeBuffer
}

// New code writer.
func newCodeWriter() *codeWriter {
	return &codeWriter {
		buffer: new(safeBuffer),
	}
}

// Indent.
func (w *codeWriter) Indent() {
	w.indent++
}

// Unindent.
func (w *codeWriter) Unindent() {
	if w.indent > 0 {
		w.indent--
	}
}

// Blank line.
func (w *codeWriter) BlankLine() {
	w.buffer.Write([]byte { '\n' })
}

// Write line.
func (w *codeWriter) Line(line string) {
	w.buffer.WriteString(strings.Repeat("    ", w.indent))
	w.buffer.WriteString(line)
	w.buffer.Write([]byte { '\n' })
}

// Write formatted line.
func (w *codeWriter) Linef(format string, a ...interface{}) {
	w.Line(fmt.Sprintf(format, a...))
}

// Write documentation.
func (w *codeWriter) Documentation(docs []string) {
	if docs == nil || len(docs) == 0 {
		return
	}

	wrapper := utils.NewSimpleTextWrapper(79 - w.indent * 4 - 3)
	lines := make([]string, 0, len(docs)*2)

	for _, paragraph := range docs {
		if len(lines) > 0 {
			lines = append(lines, "")
		} else {
			paragraph = fmt.Sprintf(`"""%s`, strings.TrimSpace(paragraph))
		}

		lines = append(lines, wrapper.Wrap(paragraph)...)
	}

	if len(lines) == 0 {
		return
	}

	for _, l := range lines {
		w.Line(l)
	}

	w.Line(`"""`)
	w.BlankLine()
}

// Write parentherized statement with arguments.
//
// Writes `<prefix>(<args...>)<suffix>` or breaks it into multiple lines with
// following lines indented by the length of the prefix and parenthesis.
func (w *codeWriter) ParentherizedWithArguments(prefix, suffix string, args ...string) {
	// Attempt with a single line version first.
	singleLine := fmt.Sprintf("%s(%s)%s", prefix, strings.Join(args, ", "), suffix)

	if w.Fits(singleLine) || len(args) == 0 {
		w.Line(singleLine)
		return
	}

	// Determine the maximum argument length.
	maxArgLength := 0
	for _, a := range args {
		if len(a) > maxArgLength {
			maxArgLength = len(a)
		}
	}

	// Build the multi line version.
	if len(prefix) + 2 + maxArgLength > w.availableSpace() {
		w.Linef("%s(", prefix)
		for i, a := range args {
			if i == len(args) - 1 {
				w.Linef("    %s", a)
			} else {
				w.Linef("    %s,", a)
			}
		}
		w.Linef(")%s", suffix)
	} else {
		indent := strings.Repeat(" ", len(prefix) + 1)

		w.Linef("%s(%s,", prefix, args[0])
		for i, a := range args {
			if i == 0 {
				continue
			}

			if i == len(args) - 1 {
				w.Linef("%s%s)%s", indent, a, suffix)
			} else {
				w.Linef("%s%s,", indent, a)
			}
		}
	}
}

// Write parentherized definition.
//
// Writes `<definition> = (<args...>)` or breaks it into multiple lines with
// following lines indented by the length of the prefix and parenthesis.
func (w *codeWriter) ParentherizedDefinition(definition string, args ...string) {
	// Attempt with a single line version first.
	singleLine := fmt.Sprintf("%s = (%s)", definition, strings.Join(args, ", "))

	if w.Fits(singleLine) || len(args) == 0 {
		w.Line(singleLine)
		return
	}

	// Build the multi line version.
	w.Linef("%s = (", definition)
	for _, a := range args {
		w.Linef("    %s,", a)
	}
	w.Line(")")
}

// Write exception raising statement.
func (w *codeWriter) RaiseException(name, message string) {
	// Attempt with a single line version first.
	if singleLine := fmt.Sprintf("raise %s('%s')", name, message); w.Fits(singleLine) {
		w.Line(singleLine)
		return
	}

	// Build the multi line version.
	messagePartLength := w.availableSpace() - len(name) - len("raise (''")
	indent := strings.Repeat(" ", len("raise (") + len(name))

	i := 0
	for i < len(message) {
		prefix := indent
		if i == 0 {
			prefix = fmt.Sprintf("raise %s(", name)
		}

		residual := len(message) - i
		if residual < messagePartLength {
			w.Linef("%s'%s')", prefix, message[i:])
		} else if residual == messagePartLength {
			w.Linef("%s'%s'", prefix, message[i:i + messagePartLength - 4])
			w.Linef("%s'%s')", indent, message[i + messagePartLength - 4:])
		} else {
			w.Linef("%s'%s'", prefix, message[i:i + messagePartLength])
		}

		i += messagePartLength
	}
}

// Write a comment.
func (w *codeWriter) Comment(comment string) {
	w.Linef("# %s", comment)
}

// Test if a line will fit within the boundary.
func (w *codeWriter) Fits(input string) bool {
	return len(input) < (80 - w.indent * 4)
}

// Available line space.
func (w *codeWriter) availableSpace() int {
	return 79 - w.indent * 4
}

// Bytes.
func (w *codeWriter) Bytes() []byte {
	return w.buffer.Bytes()
}
