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
	"strconv"
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

// Dbg will output the given message to Writer, applying any special formatting if necessary.
func Dbg(message ...interface{}) func(string, ...interface{}) string {
	if len(message) == 0 {
		dbg := _dbg{}
		return func(ctl string, message ...interface{}) string {
			process := dbg.process(ctl, message...)
			{
				message := process.compose(1)
				if process.quiet {
					return message
				}
				fmt.Fprintln(Writer, message)
				return ""
			}
		}
	}
	fmt.Fprintln(Writer, inputToString(message...))
	return nil
}

// Dbg will output the given message to Writer via fmt.Sprintf, with the first argument being the format string.
func Dbgf(message ...interface{}) {
	fmt.Fprintln(Writer, Fmt(message...))
}

// Fmt will process the message via fmt.Sprintf, with the first argument being the format string and return the result.
func Fmt(message ...interface{}) string {
	var format interface{} = ""
	if len(message) > 0 {
		format = message[0]
		message = message[1:]
	}
	return fmt.Sprintf(fmt.Sprintf("%v", format), message...)
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

func compose(format string, message ...interface{}) string {
	if format == "" {
		if len(message) == 0 {
			return ""
		}

		if FormatIf != nil {
			switch value := message[0].(type) {
			case string:
				if match := FormatIf.FindStringIndex(value); match != nil {
					format := value[match[1]:]
					message = message[1:]
					return fmt.Sprintf(format, message...)
				}
			}
		}

		output := []string{}
		for _, argument := range message {
			output = append(output, toString(argument))
		}
		return strings.Join(output, " ")
	}
	return fmt.Sprintf(format, message...)
}

type _dbg struct {
}

var (
	caller_re = regexp.MustCompile(`/@(?::(\d+))?`)
	format_re = regexp.MustCompile(`/%`)
	quiet_re = regexp.MustCompile(`/<`)
)

type _process struct {
	quiet bool
	caller int
	format bool
	message []interface{}
}

func (self _process) compose(skip int) string {
	output := []string{}

	{
		caller := self.caller
		if caller > 0 {
			caller += skip
			pc, _, _, _ := runtime.Caller(caller)
			name := runtime.FuncForPC(pc).Name()
			output = append(output, fmt.Sprintf("@ %s:", name))
		}
	}

	message := ""
	if self.format {
		if len(self.message) > 0 {
			switch format := self.message[0].(type) {
			case string:
				message = compose(format, self.message[1:]...)
			default:
				message = compose("", self.message...)
			}
		}
	} else {
		message = compose("", self.message...)
	}

	if message != "" {
		output = append(output, message)
	}
	return strings.Join(output, " ")
}

func (self _dbg) process(ctl string, message ...interface{}) _process {

	quiet := false
	caller := -1
	format := false

	if len(ctl) > 0 {
		if ctl[0] == '<' {
			quiet = true
			ctl = ctl[1:]
		}
		if match := caller_re.FindStringSubmatch(ctl); match != nil {
			caller = 1
			if match[1] != "" {
				tmp, _ := strconv.ParseInt(match[1], 10, 32)
				caller = int(tmp)
			}
		}
		if format_re.MatchString(ctl) {
			format = true
		}
		if quiet_re.MatchString(ctl) {
			quiet = true
		}
	}

	return _process{
		quiet: quiet,
		caller: caller,
		format: format,
		message: message,
	}
}

// dbg()("</@", ...)
// dbg()("/%", ...)
