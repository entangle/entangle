package golang

import (
	"testing"
)

func TestNewGenerator(t *testing.T) {
	generator, err := NewGenerator(&Options{})

	if err != nil {
		t.Errorf("Initializing Go generator failed: %v", err)
	} else if generator == nil {
		t.Errorf("Initializing Go generator did not fail, but returned a nil generator")
	}
}
