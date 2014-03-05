package python2

import (
	"fmt"
)

// Generate clients.py.
func generateClients(ctx *context) (src *SourceFile, err error) {
	src = NewSourceFile("clients")

	for _, srvc := range ctx.Interface.Services {
		clientName := fmt.Sprintf("%sClient", srvc.Name)
		src.Export(clientName)

		w := newCodeWriter()

		// Write the class definition.
		src.ImportAs("entangle.client", "Client", "Client_")
		w.Linef("class %s(Client_):", clientName)
		w.Indent()
		w.Documentation(srvc.Documentation)

		// Write each function.
		for _, fun := range srvc.Functions {
			// Write the function definition.
			args := make([]string, len(fun.Arguments) + 3)
			args[0] = "self_"
			args[len(args) - 2] = "trace=False"
			args[len(args) - 1] = "notify=False"

			for i, arg := range fun.ArgumentsSortedByIndex() {
				args[i + 1] = snakeCaseString(arg.Name)
			}

			w.ParentherizedWithArguments(fmt.Sprintf("def %s", snakeCaseString(fun.Name)), ":", args...)
			w.Indent()

			w.Documentation(fun.Documentation)

			// Write argument serialization.
			src.ImportAs("io", "BytesIO", "BytesIO_")
			w.Comment("Pack arguments.")

			w.Line("stream_ = BytesIO_()")
			w.BlankLine()

			packDecls := make([]inlinePackingDecl, fun.SerializedLength())
			for _, arg := range fun.ArgumentsSortedByIndex() {
				name := snakeCaseString(arg.Name)
				packDecls[arg.Index - 1] = inlinePackingDecl {
					Source: name,
					Description: fmt.Sprintf("argument %s", name),
					Type: arg.Type,
				}
			}
			writeInlinePacking(packDecls, w, src)

			if len(packDecls) > 0 {
				w.BlankLine()
			}
			w.Line("stream_.seek(0)")
			w.BlankLine()

			// Write the service calling.
			w.Comment("Call the service.")
			w.ParentherizedWithArguments("response = self_._call", "", fmt.Sprintf("'%s'", fun.Name), "stream_.getvalue()", "trace=trace", "notify=notify")
			w.BlankLine()

			// Write the exception response handling.
			w.Comment("Handle exceptions.")
			src.ImportAs("entangle.message", "ExceptionMessage", "ExceptionMessage_")
			src.ImportAs(".exceptions", "parse_exception", "parse_exception_")
			w.Line("if isinstance(response, ExceptionMessage_):")
			w.Indent()
			w.ParentherizedWithArguments("raise parse_exception_", "", "response.definition", "response.name", "response.description")
			w.Unindent()
			w.BlankLine()

			w.Line("if notify:")
			w.Line("    return (None, None)")
			w.BlankLine()

			// Write the result response handling.
			if fun.ReturnType != nil {
				w.Comment("Deserialize result.")
				if fun.ReturnType.Nilable() {
					w.Line("result = None")
				}
				writeSingleInlineDeserialization("response.result", "result", "result", "", fun.ReturnType, w, src)
				w.Line("return (result, response.trace)")
			} else {
				w.Line("return (None, response.trace)")
			}

			w.Unindent()
			w.BlankLine()
		}

		src.AddBlock(w.Bytes())
	}

	return
}
