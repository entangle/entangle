package declarations

import (
	"../utils"
)

// Interface declaration.
type Interface struct {
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
