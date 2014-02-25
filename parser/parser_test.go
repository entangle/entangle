package parser

import (
	"../source"
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	src, err := source.FromString(`import "horse"
import base "base_stuff"

struct Hello {

}`, "test.etg")

	if err != nil {
		panic(err)
	}

	fmt.Println(Parse(src))
}
