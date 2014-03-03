package generators

import (
	"entangle/declarations"
)

// Generator.
type Generator interface {
	// Generate.
	Generate(interfaceDecl *declarations.Interface, target string) error
}
