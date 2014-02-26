package declarations

// Struct declaration.
type Struct struct {
	// Struct name.
	Name string

	// Parent struct name.
	//
	// Empty if the struct does not inherit from a parent.
	ParentName string

	// Documentation paragraphs.
	Documentation []string

	// Fields.
	//
	// Do not modify this slice directly. Always use AddField.
	Fields []*Field

	// Field name mapping.
	fieldNameMapping map[string]*Field

	// Field index mapping.
	fieldIndexMapping map[uint]*Field
}

// New struct declaration.
func NewStruct(name string, documentation []string) *Struct {
	return &Struct{
		Name:              name,
		ParentName:        "",
		Documentation:     documentation,
		Fields:            []*Field{},
		fieldNameMapping:  map[string]*Field{},
		fieldIndexMapping: map[uint]*Field{},
	}
}

// Add a field to a struct declaration.
//
// The caller is expected to have validated that neither the name nor index are
// in use before calling AddField.
func (s *Struct) AddField(index uint, name string, documentation []string, fieldType Type) {
	field := &Field{
		Index:         index,
		Name:          name,
		Documentation: documentation,
		Type:          fieldType,
	}

	s.Fields = append(s.Fields, field)
	s.fieldNameMapping[name] = field
	s.fieldIndexMapping[index] = field
}

// Inherit from the current struct to a new struct.
func (s *Struct) Inherit(name string, documentation []string) *Struct {
	c := NewStruct(name, documentation)
	c.ParentName = s.Name

	for _, f := range s.Fields {
		c.AddField(f.Index, f.Name, f.Documentation, f.Type)
	}

	return c
}

// Determine if a field index is in use.
func (s *Struct) FieldIndexInUse(index uint) bool {
	_, inUse := s.fieldIndexMapping[index]
	return inUse
}

// Determine if a field name is in use.
func (s *Struct) FieldNameInUse(name string) bool {
	_, inUse := s.fieldNameMapping[name]
	return inUse
}
