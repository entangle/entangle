package golang

import (
	"entangle/declarations"
	"fmt"
)

var (
	deserializationSubTypeMapping = map[declarations.TypeClass]string{
		declarations.BoolClass:    "Bool",
		declarations.StringClass:  "String",
		declarations.BinaryClass:  "Binary",
		declarations.Float32Class: "Float32",
		declarations.Float64Class: "Float64",
		declarations.Int8Class:    "Int8",
		declarations.Int16Class:   "Int16",
		declarations.Int32Class:   "Int32",
		declarations.Int64Class:   "Int64",
		declarations.Uint8Class:   "Uint8",
		declarations.Uint16Class:  "Uint16",
		declarations.Uint32Class:  "Uint32",
		declarations.Uint64Class:  "Uint64",
	}
)

// Name of subtype for deserializer for type.
func nameOfDeserializerSubtype(typeDecl declarations.Type) string {
	nilable := ""
	if typeDecl.Nilable() {
		nilable = "Nilable"
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

		return fmt.Sprintf("MapOf%sTo%s", nameOfDeserializerSubtype(keyTypeDecl), nameOfDeserializerSubtype(valueTypeDecl))

	case declarations.ListClass:
		listTypeDecl := typeDecl.(*declarations.ListType)
		elementTypeDecl := listTypeDecl.ElementType()

		return fmt.Sprintf("ListOf%s", nameOfDeserializerSubtype(elementTypeDecl))

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

	return fmt.Sprintf("deserialize%s", suffix)
}

// Name of serializer for type.
//
// Returns an empty value if the type does not need a custom serializer.
func nameOfSerializer(typeDecl declarations.Type) string {
	suffix := suffixOfSerDes(typeDecl)

	if suffix == "" {
		return ""
	}

	return fmt.Sprintf("serialize%s", suffix)
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
