package python2

import (
	"bytes"
)

// Safe buffer.
//
// Easier to use as writing should never cause problems, in which case panic
// will be called.
type safeBuffer struct {
	b bytes.Buffer
}

// Write.
func (b *safeBuffer) Write(data []byte) {
	if _, err := b.b.Write(data); err != nil {
		panic(err)
	}
}

// Write string.
func (b *safeBuffer) WriteString(data string) {
	if _, err := b.b.WriteString(data); err != nil {
		panic(err)
	}
}

// Bytes.
func (b *safeBuffer) Bytes() []byte {
	return b.b.Bytes()
}
