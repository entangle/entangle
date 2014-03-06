package python2

import (
	"fmt"
	"entangle/declarations"
)

var simplePackerMapping = map[declarations.TypeClass]string {
	declarations.BoolClass: "pack_bool",
	declarations.StringClass: "pack_string",
	declarations.BinaryClass: "pack_binary",
	declarations.Float32Class: "pack_float32",
	declarations.Float64Class: "pack_float64",
	declarations.Int8Class: "pack_int8",
	declarations.Int16Class: "pack_int16",
	declarations.Int32Class: "pack_int32",
	declarations.Int64Class: "pack_int64",
	declarations.Uint8Class: "pack_uint8",
	declarations.Uint16Class: "pack_uint16",
	declarations.Uint32Class: "pack_uint32",
	declarations.Uint64Class: "pack_uint64",
}

// Inline packing declaration.
type inlinePackingDecl struct {
	// Source.
	Source string

	// Description.
	Description string

	// Type.
	//
	// If nil, a nil value is written instead of serialization.
	Type declarations.Type
}

// Write inline packing for a single type.
func writeSingleInlinePacking(source, stream, description string, typeDecl declarations.Type, w *codeWriter, src *SourceFile) {
	if typeDecl.Nilable() {
		w.Linef("if %s is not None:", source)
		w.Indent()
	}

	switch typeDecl.Class() {
	case declarations.BoolClass, declarations.StringClass, declarations.BinaryClass, declarations.Float32Class, declarations.Float64Class, declarations.Int8Class, declarations.Int16Class, declarations.Int32Class, declarations.Int64Class, declarations.Uint8Class, declarations.Uint16Class, declarations.Uint32Class, declarations.Uint64Class:
		packer := simplePackerMapping[typeDecl.Class()]
		src.ImportAs("entangle.packing", packer, fmt.Sprintf("%s_", packer))
		w.Linef("%s.write(%s_(%s))", stream, packer, source)

	case declarations.EnumClass, declarations.StructClass:
		var clsName string
		if typeDecl.Class() == declarations.StructClass {
			clsName = typeDecl.(*declarations.StructType).Struct().Name
		} else {
			clsName = typeDecl.(*declarations.EnumType).Enum().Name
		}

		if src.moduleName == "packing" {
			w.Linef("from .types import %s", clsName)
		} else {
			src.Import(".types", clsName)
		}

		src.ImportAs("entangle.exceptions", "PackingError", "PackingError_")

		w.Linef("if not isinstance(%s, %s):", source, clsName)
		w.Indent()
		w.RaiseException("PackingError_", fmt.Sprintf("%s is not an instance of %s", description, clsName))
		w.Unindent()
		w.Linef("%s.pack(%s)", source, stream)

	case declarations.MapClass, declarations.ListClass:
		requiredType := "list"
		if typeDecl.Class() == declarations.MapClass {
			requiredType = "dict"
		}

		src.ImportAs("entangle.exceptions", "PackingError", "PackingError_")

		packer := nameOfPacker(typeDecl)

		w.Linef("if not isinstance(%s, %s):", source, requiredType)
		w.Indent()
		w.RaiseException("PackingError_", fmt.Sprintf("%s is not a list", description))
		w.Unindent()

		if src.moduleName == "packing" {
			w.ParentherizedWithArguments(fmt.Sprintf(packer), "", source, stream)
		} else {
			src.ImportAs(".packing", packer, fmt.Sprintf("%s_", packer))
			w.ParentherizedWithArguments(fmt.Sprintf("%s_", packer), "", source, stream)
		}
	}

	if typeDecl.Nilable() {
		w.Unindent()
		w.Line("else:")
		w.Indent()
		w.Linef("%s.write('\\xc0')", stream)
		w.Unindent()
	}
}

// Write inline packing.
//
// Expects the output to be written to a writable called "stream_" in the code.
func writeInlinePacking(decls []inlinePackingDecl, w *codeWriter, src *SourceFile) {
	// Write None checkers.
	anyNonNilable := false

	for _, decl := range decls {
		if decl.Type == nil || decl.Type.Nilable() {
			continue
		}

		src.ImportAs("entangle.exceptions", "PackingError", "PackingError_")
		w.Linef("if %s is None:", decl.Source)
		w.Linef("    raise PackingError_('%s cannot be None')", decl.Description)
		anyNonNilable = true
	}

	if anyNonNilable {
		w.BlankLine()
	}

	// Write the array header.
	src.ImportAs("entangle.packing", "packer", "packer_")
	w.Linef("stream_.write(packer_.pack_array_header(%d))\n", len(decls))

	// Write serialization for each part.
	for _, decl := range decls {
		if decl.Type == nil {
			w.Line("stream_.write('\\xc0')")
			continue
		}

		writeSingleInlinePacking(decl.Source, "stream_", decl.Description, decl.Type, w, src)
	}
}
