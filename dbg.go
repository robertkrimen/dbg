/*
Package dbg is a printf-debugging utility library.

	dbg.Dbg("Just a string message.")
	// Just a string message

	dbg.Dbg("Or a series of values:", 1, 2.1, true)
	// Or a series of values 1 2.1 true

	dbg.Dbg("\\:Or format a %s message", "special")
	// Or format a special message

	// Change the string format
	dbg.ToStringFormat = "[%v]"
	dbg.Dbg("Now this sentence will be wrapped in [].")
	// [Now this sentence will be wrapped in [].]

*/
package dbg

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"strings"
)

var (
	// The writer to send Dbg output to.
	Writer         io.Writer = os.Stderr

	// The fmt "verb" to use when formatting a value for output.
	ToStringFormat string    = "%v"

	// The pattern to look for when detecting whether to use the
	// input as arguments to fmt.Sprintf, with the first argument 
	// being the format string.
	// Assumed to be a string at the start of the first argument.
	FormatIf                 = regexp.MustCompile(`^\\:`)
)

// Dbg will output the given input to Writer, applying any special formatting if necessary.
func Dbg(input ...interface{}) {
	fmt.Fprintln(Writer, inputToString(input...))
}

// Dbg will output the given input to Writer via fmt.Sprintf, with the first argument being the format string.
func Dbgf(input ...interface{}) {
	fmt.Fprintln(Writer, Fmt(input...))
}

// Fmt will process the input via fmt.Sprintf, with the first argument being the format string and return the result.
func Fmt(input ...interface{}) string {
	var format interface{} = ""
	if len(input) > 0 {
		format = input[0]
		input = input[1:]
	}
	return fmt.Sprintf(fmt.Sprintf("%v", format), input...)
}

// HereBeDragons returns a string formatted like Dbg, but prefixed with the call location (function name).
func HereBeDragons(input ...interface{}) string {
	pc, _, _, _ := runtime.Caller(1)
	name := runtime.FuncForPC(pc).Name()
	message := fmt.Sprintf("Here be dragons @ %s", name)
	if len(input) > 0 {
		message += fmt.Sprintf(": %s", inputToString(input))
	}
	return message
}

func inputToString(input ...interface{}) string {
	if len(input) == 0 {
		return ""
	}
	if FormatIf != nil {
		switch value := input[0].(type) {
		case string:
			if match := FormatIf.FindStringIndex(value); match != nil {
				format := value[match[1]:]
				input = input[1:]
				return fmt.Sprintf(format, input...)
			}
		}
	}

	// No (Sprintf) formatting done
	output := []string{}
	for _, argument := range input {
		output = append(output, toString(argument))
	}
	return strings.Join(output, " ")
}

func toString(value interface{}) string {
	return fmt.Sprintf(ToStringFormat, value)
}

