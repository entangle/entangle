package declarations

import (
	"entangle/utils"
	"sort"
)

// Interface declaration.
type Interface struct {
	// Interface name.
	Name string

	// Struct declarations.
	Structs map[string]*Struct

	// Exception declarations.
	Exceptions map[string]*Exception

	// Enumeration declarations.
	Enums map[string]*Enum

	// Service declarations.
	Services map[string]*Service

	// Used names.
	usedNames utils.StringSet
}

// New interface declaration.
func NewInterface() *Interface {
	return &Interface{
		Structs:    map[string]*Struct{},
		Exceptions: map[string]*Exception{},
		Enums:      map[string]*Enum{},
		Services:   map[string]*Service{},
		usedNames:  make(utils.StringSet),
	}
}

// Add a struct and mark its name as used.
func (i *Interface) AddStruct(decl *Struct) {
	i.Structs[decl.Name] = decl
	i.MarkNameAsUsed(decl.Name)
}

// Add an exception and mark its name as used.
func (i *Interface) AddException(decl *Exception) {
	i.Exceptions[decl.Name] = decl
	i.MarkNameAsUsed(decl.Name)
}

// Add an enumeration and mark its name as used.
func (i *Interface) AddEnum(decl *Enum) {
	i.Enums[decl.Name] = decl
	i.MarkNameAsUsed(decl.Name)
}

// Add a service and mark its name as used.
func (i *Interface) AddService(decl *Service) {
	i.Services[decl.Name] = decl
	i.MarkNameAsUsed(decl.Name)
}

// Test if a name is in use.
func (i *Interface) NameInUse(name string) bool {
	return i.usedNames.Contains(name)
}

// Mark a name as used.
func (i *Interface) MarkNameAsUsed(name string) {
	i.usedNames.Add(name)
}

// Sorted list of exceptions by name.
func (i *Interface) ExceptionsSortedByName() []*Exception {
	unsorted := make([]*Exception, len(i.Exceptions))

	idx := 0
	for _, exc := range i.Exceptions {
		unsorted[idx] = exc
		idx++
	}

	sort.Sort(exceptionsByName(unsorted))

	return unsorted
}

// Sorted list of enumerations by name.
func (i *Interface) EnumsSortedByName() []*Enum {
	unsorted := make([]*Enum, len(i.Enums))

	idx := 0
	for _, exc := range i.Enums {
		unsorted[idx] = exc
		idx++
	}

	sort.Sort(enumsByName(unsorted))

	return unsorted
}
