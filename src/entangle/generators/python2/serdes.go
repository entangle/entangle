package python2

import (
	"entangle/declarations"
	"fmt"
)

var (
	deserializationSubTypeMapping = map[declarations.TypeClass]string{
		declarations.BoolClass:    "bool",
		declarations.StringClass:  "string",
		declarations.BinaryClass:  "binary",
		declarations.Float32Class: "float32",
		declarations.Float64Class: "float64",
		declarations.Int8Class:    "int8",
		declarations.Int16Class:   "int16",
		declarations.Int32Class:   "int32",
		declarations.Int64Class:   "int64",
		declarations.Uint8Class:   "uint8",
		declarations.Uint16Class:  "uint16",
		declarations.Uint32Class:  "uint32",
		declarations.Uint64Class:  "uint64",
	}
)

// Name of subtype for deserializer for type.
func nameOfDeserializerSubtype(typeDecl declarations.Type) string {
	nilable := ""
	if typeDecl.Nilable() {
		nilable = "nilable_"
	}

	if simpleName, ok := deserializationSubTypeMapping[typeDecl.Class()]; ok {
		return fmt.Sprintf("%s%s", nilable, simpleName)
	}

	switch typeDecl.Class() {
	case declarations.EnumClass:
		return fmt.Sprintf("%s%s", nilable, typeDecl.(*declarations.EnumType).Enum().Name)

	case declarations.StructClass:
		return fmt.Sprintf("%s%s", nilable, typeDecl.(*declarations.StructType).Struct().Name)

	case declarations.MapClass, declarations.ListClass:
		return fmt.Sprintf("%s%s", nilable, suffixOfSerDes(typeDecl))

	default:
		panic("Unimplemented type")
	}
}

// Prefix of deserializer for type.
//
// Returns an empty value if the type does not need a custom deserializer.
func suffixOfSerDes(typeDecl declarations.Type) string {
	switch typeDecl.Class() {
	case declarations.MapClass:
		mapTypeDecl := typeDecl.(*declarations.MapType)
		keyTypeDecl := mapTypeDecl.KeyType()
		valueTypeDecl := mapTypeDecl.ValueType()

		return fmt.Sprintf("map_of_%s_to_%s", nameOfDeserializerSubtype(keyTypeDecl), nameOfDeserializerSubtype(valueTypeDecl))

	case declarations.ListClass:
		listTypeDecl := typeDecl.(*declarations.ListType)
		elementTypeDecl := listTypeDecl.ElementType()

		return fmt.Sprintf("list_of_%s", nameOfDeserializerSubtype(elementTypeDecl))

	default:
		return ""
	}
}

// Name of deserializer for type.
//
// Returns an empty value if the type does not need a custom deserializer.
func nameOfDeserializer(typeDecl declarations.Type) string {
	suffix := suffixOfSerDes(typeDecl)

	if suffix == "" {
		return ""
	}

	return fmt.Sprintf("deserialize_%s", suffix)
}

// Name of packer for type.
//
// Returns an empty value if the type does not need a custom packer.
func nameOfPacker(typeDecl declarations.Type) string {
	suffix := suffixOfSerDes(typeDecl)

	if suffix == "" {
		return ""
	}

	return fmt.Sprintf("pack_%s", suffix)
}

// Map a type to a serialization/deserialization map.
func mapTypeToSerDesMap(typeDecl declarations.Type, m *map[string]declarations.Type) {
	suffix := suffixOfSerDes(typeDecl)
	if suffix == "" {
		return
	}

	if _, found := (*m)[suffix]; found {
		return
	}

	switch typeDecl.Class() {
	case declarations.MapClass:
		mapTypeDecl := typeDecl.(*declarations.MapType)
		valueTypeDecl := mapTypeDecl.ValueType()

		switch valueTypeDecl.Class() {
		case declarations.MapClass, declarations.ListClass:
			mapTypeToSerDesMap(valueTypeDecl, m)
		}

	case declarations.ListClass:
		listTypeDecl := typeDecl.(*declarations.ListType)
		elementTypeDecl := listTypeDecl.ElementType()

		switch elementTypeDecl.Class() {
		case declarations.MapClass, declarations.ListClass:
			mapTypeToSerDesMap(elementTypeDecl, m)
		}
	}

	(*m)[suffix] = typeDecl
	return
}

// Build a serialization/deserialization map for interface.
func buildSerDesMap(interfaceDecl *declarations.Interface) (m map[string]declarations.Type) {
	m = make(map[string]declarations.Type)

	// Iterate across all function arguments in services.
	for _, service := range interfaceDecl.Services {
		for _, function := range service.Functions {
			for _, argument := range function.Arguments {
				mapTypeToSerDesMap(argument.Type, &m)
			}

			if function.ReturnType != nil {
				mapTypeToSerDesMap(function.ReturnType, &m)
			}
		}
	}

	// Iterate across all struct fields.
	for _, structDecl := range interfaceDecl.Structs {
		for _, field := range structDecl.Fields {
			mapTypeToSerDesMap(field.Type, &m)
		}
	}

	return
}
