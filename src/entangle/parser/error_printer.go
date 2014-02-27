package parser

import (
	"entangle/errors"
	"entangle/term"
	"entangle/utils"
	"fmt"
	"strings"
)

const tabWidth = 4

// Print an error in a human readable format.
//
// Essentially a carbon copy of how Clang prints errors, because it's so darned
// helpful.
func PrintError(err errors.ParseError) {
	// Print each frame.
	for i, frame := range err.Frames() {
		// Print the source and description.
		term.Printf(term.BOLD, "%s:%d:%d: ", frame.Source.Path(), frame.Start.Line, frame.Start.Character)

		if i < len(err.Frames())-1 {
			term.Printf(term.BOLD|term.MAGENTA, "imported from here")
		} else {
			term.Printf(term.BOLD|term.RED, "error: ")
			term.Printf(term.BOLD, err.Description())
		}

		fmt.Println()

		// Print the problematic line.
		line := frame.Source.Line(frame.Start.Line)
		fmt.Println(utils.ExpandTabs(line, tabWidth))

		// Print the pointing arrow and curly marker.
		start := frame.Start.Character
		end := frame.End.Character

		if frame.End.Line > frame.Start.Line {
			end = len(line) + 1
		}

		if start > 1 {
			fmt.Print(utils.MaskWithWhitespaceExpanded(line[:start-1], tabWidth))
		}

		term.Printf(term.GREEN, "^")

		if end > start {
			term.Printf(term.GREEN, strings.Repeat("~", end-start))
		}

		fmt.Println()

		// Create an extra empty white line between frames.
		if i < len(err.Frames())-1 {
			fmt.Println()
		}
	}
}
