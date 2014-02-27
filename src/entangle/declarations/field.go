package declarations

// Field declaration.
type Field struct {
	// Field index.
	Index uint

	// Field name.
	Name string

	// Documentation paragraphs.
	Documentation []string

	// Type.
	Type Type
}
