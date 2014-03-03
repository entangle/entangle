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

// Fields by index.
type fieldsByIndex []*Field

func (l fieldsByIndex) Len() int {
	return len(l)
}

func (l fieldsByIndex) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l fieldsByIndex) Less(i, j int) bool {
	return l[i].Index < l[j].Index
}
