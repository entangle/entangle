package token

// Position.
type Position struct {
	// Line.
	Line int

	// Character.
	Character int
}

// Get the position immediately before the current position.
//
// Only valid if the character position is greater than 1.
func (p Position) Before() Position {
	return Position{
		Line:      p.Line,
		Character: p.Character - 1,
	}
}
