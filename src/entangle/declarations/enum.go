package declarations

import (
	"sort"
)

// Enumeration value declaration.
type EnumValue struct {
	// Value.
	Value int64

	// Name.
	Name string

	// Documentation paragraphs.
	Documentation []string
}

// Enumeration declaration.
type Enum struct {
	// Struct name.
	Name string

	// Documentation paragraphs.
	Documentation []string

	// Values.
	//
	// Mapping of values to representation.
	Values map[int64]EnumValue
}

// New enumeration declaration.
func NewEnum(name string, documentation []string) *Enum {
	return &Enum{
		Name:          name,
		Documentation: documentation,
		Values:        make(map[int64]EnumValue),
	}
}

// Add a value.
func (e *Enum) AddValue(value int64, name string, documentation []string) {
	e.Values[value] = EnumValue{
		Value:         value,
		Name:          name,
		Documentation: documentation,
	}
}

// Determine if a value is in use.
func (e *Enum) ValueInUse(value int64) bool {
	_, inUse := e.Values[value]
	return inUse
}

// Enums by name.
type enumsByName []*Enum

func (l enumsByName) Len() int {
	return len(l)
}

func (l enumsByName) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l enumsByName) Less(i, j int) bool {
	return l[i].Name < l[j].Name
}

// Enum values by value.
type enumValuesByValue []EnumValue

func (l enumValuesByValue) Len() int {
	return len(l)
}

func (l enumValuesByValue) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l enumValuesByValue) Less(i, j int) bool {
	return l[i].Value < l[j].Value
}

// Sorted list of enum values by value.
func (e *Enum) ValuesSortedByValue() []EnumValue {
	unsorted := make([]EnumValue, len(e.Values))

	idx := 0
	for _, val := range e.Values {
		unsorted[idx] = val
		idx++
	}

	sort.Sort(enumValuesByValue(unsorted))

	return unsorted
}
