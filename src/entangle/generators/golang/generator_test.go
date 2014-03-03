package golang

import (
	"testing"
)

func TestNewGenerator(t *testing.T) {
	generator, err := NewGenerator()

	if err != nil {
		t.Errorf("Initializing Go generator failed: %v", err)
	} else if generator == nil {
		t.Errorf("Initializing Go generator did not fail, but returned a nil generator")
	}
}

func TestgeneratorGenerate(t *testing.T) {
	generator, err := NewGenerator()
	if err != nil {
		t.Fatalf("Initializing Go generator failed: %v", err)
	}
}
