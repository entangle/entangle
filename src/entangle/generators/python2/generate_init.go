package python2

// Generate __init__.py.
func generateInit(ctx *context) (src *SourceFile, err error) {
	src = NewSourceFile("__init__")
	return
}
