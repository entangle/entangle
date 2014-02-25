package term

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type OutputFormatting uint64

const (
	NONE = iota
	BOLD OutputFormatting = 1 << (iota)
	BLACK
	RED
	GREEN
	YELLOW
	BLUE
	MAGENTA
	CYAN
	WHITE
	BACKGROUND_BLACK
	BACKGROUND_RED
	BACKGROUND_GREEN
	BACKGROUND_YELLOW
	BACKGROUND_BLUE
	BACKGROUND_MAGENTA
	BACKGROUND_CYAN
	BACKGROUND_WHITE
)


func Printf(formatting OutputFormatting, format string, a ...interface{}) (n int, err error) {
	// Determine the ANSI escape codes.
	flags := make([]string, 0, 3)
	
	if formatting & BOLD != 0 { flags = append(flags, "1") }
	
	for i := 0; i < 8; i++ {
		if formatting & (BLACK << uint(i)) == 0 { continue }
		flags = append(flags, strconv.Itoa(30 + i))
		break
	}
	
	for i := 0; i < 8; i++ {
		if formatting & (BACKGROUND_BLACK << uint(i)) == 0 { continue }
		flags = append(flags, strconv.Itoa(40 + i))
		break
	}
	
	// Print the formatted output.
	return fmt.Printf("\x1b[%sm%s\x1b[0m", strings.Join(flags, ";"), fmt.Sprintf(format, a...))
}


func Println(formatting OutputFormatting, a ...interface{}) (n int, err error) {
	return Printf(formatting, "%s", fmt.Sprintln(a...))
}


func Errorln(a ...interface {}) (n int, err error) {
	n, err = Printf(BOLD | RED, "Error: ")
	if err != nil { return }
	n, err = fmt.Println(a...)
	return
}


func Errorlnf(format string, a ...interface {}) (n int, err error) {
	return Errorln(fmt.Sprintf(format, a...))
}


func Fatalln(a ...interface {}) {
	Errorln(a...)
	os.Exit(1)
}


func Fatallnf(format string, a ...interface {}) {
	Errorlnf(format, a...)
	os.Exit(1)
}
