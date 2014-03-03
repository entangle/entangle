package declarations

// Exception declaration.
type Exception struct {
	// Exception name.
	Name string

	// Documentation paragraphs.
	Documentation []string
}

// New exception declaration.
func NewException(name string, documentation []string) *Exception {
	return &Exception{
		Name:          name,
		Documentation: documentation,
	}
}

// Exceptions by name.
type exceptionsByName []*Exception

func (l exceptionsByName) Len() int {
	return len(l)
}

func (l exceptionsByName) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l exceptionsByName) Less(i, j int) bool {
	return l[i].Name < l[j].Name
}
