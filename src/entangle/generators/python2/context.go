package python2

import (
	"entangle/declarations"
)

// Template context.
type context struct {
	// Interface definition.
	Interface *declarations.Interface

	// Serialization/deserialization mapping.
	SerDesMap map[string]declarations.Type
}
