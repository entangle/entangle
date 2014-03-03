package declarations

import (
	"sort"
)

// Service declaration.
type Service struct {
	// Service name.
	Name string

	// Parent service name.
	//
	// Empty if the service does not inherit from a parent.
	ParentName string

	// Documentation paragraphs.
	Documentation []string

	// Functions.
	//
	// Do not modify this slice directly. Always use AddFunction.
	Functions []*Function

	// Function name mapping.
	functionNameMapping map[string]*Function
}

// New service declaration.
func NewService(name string, documentation []string) *Service {
	return &Service{
		Name:                name,
		ParentName:          "",
		Documentation:       documentation,
		Functions:           []*Function{},
		functionNameMapping: map[string]*Function{},
	}
}

// Add a function to a service declaration.
//
// The caller is expected to have validated that neither the name nor index are
// in use before calling AddFunction.
func (s *Service) AddFunction(function *Function) {
	s.Functions = append(s.Functions, function)
	s.functionNameMapping[function.Name] = function
}

// Inherit from the current service to a new service.
func (s *Service) Inherit(name string, documentation []string) *Service {
	c := NewService(name, documentation)
	c.ParentName = s.Name

	return c
}

// Determine if a function name is in use.
func (s *Service) FunctionNameInUse(name string) bool {
	_, inUse := s.functionNameMapping[name]
	return inUse
}

// Sorted list of functions by name.
func (s *Service) FunctionsSortedByName() []*Function {
	unsorted := make([]*Function, len(s.Functions))

	idx := 0
	for _, exc := range s.Functions {
		unsorted[idx] = exc
		idx++
	}

	sort.Sort(functionsByName(unsorted))

	return unsorted
}
