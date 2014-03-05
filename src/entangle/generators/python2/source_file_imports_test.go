package python2

import (
	"math"
	"testing"
)

func TestModuleRelativity(t *testing.T) {
	for _, input := range []struct {
		Name string
		Expected uint32
	} {
		{ ".", 0 },
		{ ".some", 0 },
		{ ".some.submodule", 0 },
		{ "..", 1 },
		{ "..some", 1 },
		{ "..some.submodule", 1 },
		{ "socket", math.MaxUint32 },
	} {
		actual := moduleRelativity(input.Name)
		if actual != input.Expected {
			t.Errorf("expected module relativity of '%s' to be %d but it is %d", input.Name, input.Expected, actual)
		}
	}
}
