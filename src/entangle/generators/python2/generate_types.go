package python2

import (
	"fmt"
)

// Generate types.py.
func generateTypes(ctx *context) (src *SourceFile, err error) {
	src = NewSourceFile("types")

	// Generate enumerations.
	for _, enum := range ctx.Interface.Enums {
		src.Export(enum.Name)
		buffer := new(safeBuffer)

		// Write the enum class definition.
		src.ImportAs("entangle.types", "Enum", "Enum_")
		buffer.WriteString(fmt.Sprintf("class %s(Enum_):\n", enum.Name))

		writeDocumentation(buffer, enum.Documentation, 1)

		// Write each value.
		for _, value := range enum.ValuesSortedByValue() {
			buffer.WriteString(fmt.Sprintf("    %s = %d\n", value.Name, value.Value))
			writeDocumentation(buffer, value.Documentation, 1)
		}

		if len(enum.Values) == 0 {
			buffer.WriteString("    pass")
		}

		src.AddBlock(buffer.Bytes())
	}

	// Generate structs.
	for _, strct := range ctx.Interface.Structs {
		src.Export(strct.Name)
		w := newCodeWriter()

		// Write the enum class definition.
		w.Linef("class %s(object):", strct.Name)
		w.Indent()
		w.Documentation(strct.Documentation)

		// Build names.
		pyNameMapping := make(map[string]string, len(strct.Fields))
		for _, field := range strct.Fields {
			pyNameMapping[field.Name] = snakeCaseString(field.Name)
		}

		// Write the slots.
		fieldNames := make([]string, len(strct.Fields))
		i := 0

		for _, field := range strct.Fields {
			fieldNames[i] = pyNameMapping[field.Name]
			i++
		}

		w.ParentherizedDefinition("__slots__", stringifyStrings(fieldNames)...)
		w.BlankLine()

		// Write the initializer.
		if len(fieldNames) > 0 {
			nonedArgs := make([]string, len(fieldNames) + 1)
			nonedArgs[0] = "self"
			for i, n := range fieldNames {
				nonedArgs[i + 1] = fmt.Sprintf("%s=None", n)
			}

			w.ParentherizedWithArguments("def __init__", ":", nonedArgs...)

			w.Indent()

			for _, name := range fieldNames {
				w.Linef("self.%s = %s", name, name)
			}

			w.Unindent()
			w.BlankLine()
		}

		// Write the packer.
		w.Line("def pack(self, stream_):")
		w.Indent()
		w.Line(`"""Pack.`)
		w.BlankLine()
		w.Line(`:param stream_: Stream to pack the type to.`)
		w.Line(`:raises entangle.PackingError:`)
		w.Line(`    if the data structure could not be packed.`)
		w.Line(`"""`)
		w.BlankLine()

		decls := make([]inlinePackingDecl, strct.SerializedLength())
		for _, field := range strct.Fields {
			decls[field.Index - 1] = inlinePackingDecl {
				Source: fmt.Sprintf("self.%s", pyNameMapping[field.Name]),
				Description: fmt.Sprintf("property %s", pyNameMapping[field.Name]),
				Type: field.Type,
			}
		}
		writeInlinePacking(decls, w, src)

		w.Unindent()

		if len(decls) > 0 {
			w.BlankLine()
		}

		// Write the deserializer.
		w.Line("@classmethod")
		w.Line("def deserialize(cls, ser):")
		w.Indent()
		w.Line(`"""Deserialize.`)
		w.BlankLine()
		w.Line(`:raises entangle.DeserializationError:`)
		w.Line(`    if the serialized input could not be deserialized.`)
		w.Line(`"""`)
		w.BlankLine()
		w.Line("des = cls()")
		w.BlankLine()

		desDecls := make([]inlineDeserializationDecl, strct.SerializedLength())
		for _, field := range strct.Fields {
			desDecls[field.Index - 1] = inlineDeserializationDecl {
				Target: fmt.Sprintf("des.%s", pyNameMapping[field.Name]),
				Description: fmt.Sprintf("property %s", pyNameMapping[field.Name]),
				Type: field.Type,
			}
		}
		writeInlineDeserialization(desDecls, strct.Name, w, src)

		w.BlankLine()
		w.Line("return des")

		w.Unindent()

		src.AddBlock(w.Bytes())
	}

	return
}
