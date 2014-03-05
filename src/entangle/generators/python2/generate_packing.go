package python2

import (
	"fmt"
	"entangle/declarations"
)

// Generate packing.py.
func generatePacking(ctx *context) (src *SourceFile, err error) {
	src = NewSourceFile("packing")

	// Write each needed deserializer.
	for suffix, typeDecl := range ctx.SerDesMap {
		funName := fmt.Sprintf("pack_%s", suffix)
		src.Export(funName)

		w := newCodeWriter()

		// Write the method definition.
		w.Linef("def %s(value, stream):", funName)
		w.Indent()

		src.ImportAs("entangle.exceptions", "PackingError", "PackingError_")
		src.ImportAs("entangle.packing", "packer", "packer_")

		switch typeDecl.Class() {
		case declarations.ListClass:
			listDecl := typeDecl.(*declarations.ListType)
			elemType := listDecl.ElementType()

			w.Line("if value is None or not isinstance(value, (list, tuple)):")
			w.Indent()
			w.RaiseException("PackingError_", fmt.Sprintf("cannot pack input as a list"))
			w.Unindent()
			w.BlankLine()

			w.Line("stream.write(packer_.pack_array_header(len(value)))")
			w.BlankLine()

			w.Line("for ser in value:")
			w.Indent()
			writeSingleInlinePacking("ser", "stream", "list element", elemType, w, src)
			w.Unindent()

		case declarations.MapClass:
			mapDecl := typeDecl.(*declarations.MapType)
			keyType := mapDecl.KeyType()
			valueType := mapDecl.ValueType()

			w.Line("if value is None or not isinstance(value, dict):")
			w.Indent()
			w.RaiseException("PackingError_", fmt.Sprintf("cannot pack input as a map"))
			w.Unindent()
			w.BlankLine()

			w.Line("result = {}")
			w.BlankLine()

			w.Line("for ser_key, ser_value in value:")
			w.Indent()

			w.Line("des_key, des_value = None")
			writeSingleInlineDeserialization("ser_key", "des_key", "map key", "", keyType, w, src)
			writeSingleInlineDeserialization("ser_value", "des_value", "map value", "", valueType, w, src)
			w.Line("result[des_key] = des_value")

			w.Unindent()

			w.BlankLine()
			w.Line("return result")

		default:
			panic("Cannot generate deserialization code for type")
		}

		w.Unindent()

		src.AddBlock(w.Bytes())
	}

	return
}
