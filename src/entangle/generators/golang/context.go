package golang

import (
	"entangle/declarations"
)

// Template context.
type context struct {
	// Interface definition.
	Interface *declarations.Interface

	// Serialization/deserialization mapping.
	SerDesMap map[string]declarations.Type

	// Package name.
	PackageName string
}
