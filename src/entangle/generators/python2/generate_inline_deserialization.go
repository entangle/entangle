package python2

import (
	"entangle/declarations"
	"fmt"
)

var simpleDeserializerMapping = map[declarations.TypeClass]string {
	declarations.BoolClass: "deserialize_bool",
	declarations.StringClass: "deserialize_string",
	declarations.BinaryClass: "deserialize_binary",
	declarations.Float32Class: "deserialize_float32",
	declarations.Float64Class: "deserialize_float64",
	declarations.Int8Class: "deserialize_int8",
	declarations.Int16Class: "deserialize_int16",
	declarations.Int32Class: "deserialize_int32",
	declarations.Int64Class: "deserialize_int64",
	declarations.Uint8Class: "deserialize_uint8",
	declarations.Uint16Class: "deserialize_uint16",
	declarations.Uint32Class: "deserialize_uint32",
	declarations.Uint64Class: "deserialize_uint64",
}

// Inline deserialization declaration.
type inlineDeserializationDecl struct {
	// Target.
	Target string

	// Description.
	Description string

	// Type.
	//
	// If nil, a nil value is written instead of serialization.
	Type declarations.Type
}

// Write inline deserialization of a single variable.
//
// If a predicate is to be evaluated prior to accessing the source, the
// predicate must be provided.
func writeSingleInlineDeserialization(source, target, description, predicate string, typeDecl declarations.Type, w *codeWriter, src *SourceFile) {
	src.ImportAs("entangle.exceptions", "DeserializationError", "DeserializationError_")

	if !typeDecl.Nilable() {
		if predicate != "" {
			w.Linef("if %s:", predicate)
			w.Indent()
		}

		w.Linef("if %s is None:", source)
		w.Indent()
		w.RaiseException("DeserializationError_", fmt.Sprintf("%s is unexpectedly None", description))
		w.Unindent()
	} else {
		if predicate != "" {
			w.Linef("if %s and %s is not None:", predicate, source)
		} else {
			w.Linef("if %s is not None:", source)
		}

		w.Indent()
	}

	switch typeDecl.Class() {
	case declarations.BoolClass, declarations.StringClass, declarations.BinaryClass, declarations.Float32Class, declarations.Float64Class, declarations.Int8Class, declarations.Int16Class, declarations.Int32Class, declarations.Int64Class, declarations.Uint8Class, declarations.Uint16Class, declarations.Uint32Class, declarations.Uint64Class:
		deserializer := simpleDeserializerMapping[typeDecl.Class()]
		src.ImportAs("entangle.deserialization", deserializer, fmt.Sprintf("%s_", deserializer))
		w.Linef("%s = %s_(%s)", target, deserializer, source)

	case declarations.EnumClass, declarations.StructClass:
		var clsName string
		if typeDecl.Class() == declarations.StructClass {
			clsName = typeDecl.(*declarations.StructType).Struct().Name
		} else {
			clsName = typeDecl.(*declarations.EnumType).Enum().Name
		}

		if src.moduleName == "deserialization" {
			w.Linef("from .types import %s", clsName)
		} else {
			src.Import(".types", clsName)
		}

		w.Linef("%s = %s.deserialize(%s)", target, clsName, source)

	case declarations.MapClass, declarations.ListClass:
		deserializer := nameOfDeserializer(typeDecl)
		if src.moduleName == "deserialization" {
			w.ParentherizedWithArguments(fmt.Sprintf("%s = %s", target, deserializer), "", source)
		} else {
			src.ImportAs(".deserialization", deserializer, fmt.Sprintf("%s_", deserializer))
			w.ParentherizedWithArguments(fmt.Sprintf("%s = %s_", target, deserializer), "", source)
		}
	}

	if typeDecl.Nilable() || predicate != "" {
		w.Unindent()
	}
}

// Write inline deserialization.
//
// Expects the serialized input to be available in a variable named "ser" in
// the code.
func writeInlineDeserialization(decls []inlineDeserializationDecl, targetDesc string, w *codeWriter, src *SourceFile) {
	// Determine the minimum length of the deserialized array.
	minLength := 0

	for i, decl := range decls {
		if decl.Type != nil && !decl.Type.Nilable() {
			minLength = i + 1
		}
	}

	// Write type and length validation for serialized input.
	src.ImportAs("entangle.exceptions", "DeserializationError", "DeserializationError_")

	w.Line("if not isinstance(ser, (list, tuple)):")
	w.Indent()
	w.RaiseException("DeserializationError_", fmt.Sprintf("deserialization of %s requires a list or tuple as input", targetDesc))
	w.Unindent()

	if len(decls) == 0 {
		return
	}

	w.Linef("if not len(ser) < %d:", minLength)
	w.Indent()
	w.RaiseException("DeserializationError_", fmt.Sprintf("insufficient data to deserialize %s", targetDesc))
	w.Unindent()

	// Write deserialization for each part.
	for i, decl := range decls {
		if decl.Type == nil {
			continue
		}

		source := fmt.Sprintf("ser[%d]", i)
		predicate := ""
		if decl.Type.Nilable() && i >= minLength {
			predicate = fmt.Sprintf("len(ser) > %d", i)
		}

		writeSingleInlineDeserialization(source, decl.Target, decl.Description, predicate, decl.Type, w, src)
	}
}
