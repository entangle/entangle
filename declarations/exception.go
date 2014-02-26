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
