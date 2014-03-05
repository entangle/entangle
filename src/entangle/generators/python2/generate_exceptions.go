package python2

import (
	"fmt"
)

// Generate exceptions.py.
func generateExceptions(ctx *context) (src *SourceFile, err error) {
	src = NewSourceFile("exceptions")

	// Write the individual exceptions.
	for _, exc := range ctx.Interface.Exceptions {
		src.Export(exc.Name)

		w := newCodeWriter()

		// Write the exception definition.
		src.ImportAs("entangle.exceptions", "EntangleException", "EntangleException_")
		w.Linef("class %s(EntangleException_):", exc.Name)
		w.Indent()
		w.Documentation(exc.Documentation)

		w.Linef("definition = '%s'", ctx.Interface.Name)
		w.Linef("name = '%s'", exc.Name)

		w.Unindent()

		src.AddBlock(w.Bytes())
	}

	// Write the exception mapping followed by the exception parser.
	w := newCodeWriter()
	src.ImportAs("entangle.exceptions", "parse_exception", "entangle_parse_exception")

	mapping := make(map[string]string, len(ctx.Interface.Exceptions))

	for _, exc := range ctx.Interface.Exceptions {
		mapping[fmt.Sprintf("'%s'", exc.Name)] = exc.Name
	}

	w.DictDefinition("exceptions", mapping)

	w.BlankLine()
	w.BlankLine()

	w.Line(`def parse_exception(definition, name, message):`)
	w.Indent()
	w.Line(`"""Parse an exception.`)
	w.BlankLine()
	w.Line(`:param definition: Definition.`)
	w.Line(`:param name: Name.`)
	w.Line(`:param message: Message.`)
	w.Line(`:returns:`)
	w.Line("    the parsed exception or :class:`entangle.exceptions.UnknownException`")
	w.Line(`    if the exception is not known.`)
	w.Line(`"""`)
	w.BlankLine()
	w.Linef(`if definition == '%s':`, ctx.Interface.Name)
	w.Line(`    try:`)
	w.Line(`        return exceptions[name](message)`)
	w.Line(`    except KeyError:`)
	w.Line(`        pass`)
	w.BlankLine()
	w.Line(`return entangle_parse_exception(definition, name, message)`)
	w.Unindent()

	src.Export("parse_exception")
	src.AddBlock(w.Bytes())

	return
}
