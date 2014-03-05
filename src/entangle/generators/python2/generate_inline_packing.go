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

		if decl.Type.Nilable() {
			w.Linef("if %s is not None:", decl.Source)
			w.Indent()
		}

		switch decl.Type.Class() {
		case declarations.BoolClass, declarations.StringClass, declarations.BinaryClass, declarations.Float32Class, declarations.Float64Class, declarations.Int8Class, declarations.Int16Class, declarations.Int32Class, declarations.Int64Class, declarations.Uint8Class, declarations.Uint16Class, declarations.Uint32Class, declarations.Uint64Class:
			packer := simplePackerMapping[decl.Type.Class()]
			src.ImportAs("entangle.packing", packer, fmt.Sprintf("%s_", packer))
			w.Linef("stream_.write(%s_(%s))", packer, decl.Source)

		case declarations.EnumClass:
			enumName := decl.Type.(*declarations.EnumType).Enum().Name

			src.ImportAs("entangle.exceptions", "PackingError", "PackingError_")
			src.Import(".types", enumName)

			w.Linef("if not isinstance(%s, %s):", decl.Source, enumName)
			w.Indent()
			w.RaiseException("PackingError_", fmt.Sprintf("%s is not an instance of %s", decl.Description, enumName))
			w.Unindent()
			w.Linef("stream_.write(%s.pack())", decl.Source)

		case declarations.StructClass:
			structName := decl.Type.(*declarations.StructType).Struct().Name

			src.ImportAs("entangle.exceptions", "PackingError", "PackingError_")
			src.Import(".types", structName)

			w.Linef("if not isinstance(%s, %s):", decl.Source, structName)
			w.Indent()
			w.RaiseException("PackingError_", fmt.Sprintf("%s is not an instance of %s", decl.Description, structName))
			w.Unindent()
			w.Linef("%s.pack(stream_)", decl.Source)

		case declarations.MapClass, declarations.ListClass:
			requiredType := "list"
			if decl.Type.Class() == declarations.MapClass {
				requiredType = "dict"
			}

			packer := nameOfPacker(decl.Type)
			src.ImportAs("entangle.exceptions", "PackingError", "PackingError_")
			src.ImportAs(".packing", packer, fmt.Sprintf("%s_", packer))

			w.Linef("if not isinstance(%s, %s):", decl.Source, requiredType)
			w.Indent()
			w.RaiseException("PackingError_", fmt.Sprintf("%s is not a list", decl.Description))
			w.Unindent()
			w.ParentherizedWithArguments(fmt.Sprintf("%s_", packer), "", decl.Source, "stream_")
		}

		if decl.Type.Nilable() {
			w.Unindent()
			w.Line("else:")
			w.Indent()
			w.Line("stream_.write('\\xc0')")
			w.Unindent()
		}
	}
}
