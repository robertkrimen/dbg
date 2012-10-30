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
	FormatIf                 = regexp.MustCompile(`^\\[:%]`)
	CommandIf				 = regexp.MustCompile(`^\\/`)
)

var  _debugDebug *Debug
func _debug() *Debug {
	if _debugDebug == nil {
		_debugDebug = New()
	}
	return _debugDebug
}

// Dbg will output the given message to Writer, applying any special formatting if necessary.
func Dbg(message ...interface{}) Fn {
	return _debug().Dbg(message...)
}

type Debug struct {
	Writer io.Writer
	ToStringFormat string
	FormatIf *regexp.Regexp
	CommandIf *regexp.Regexp
}

func New() *Debug {
	return &Debug{
		Writer: Writer,
		ToStringFormat: ToStringFormat,
		FormatIf: FormatIf,
		CommandIf: CommandIf,
	}
}

type Fn func(string, ...interface{}) string

func (self *Debug) Dbg(message ...interface{}) Fn {
	if len(message) == 0 {
		return func(command string, message ...interface{}) string {
			cmd := self.command(command, message...)
			msg := self.run(cmd, 1)
			if cmd.quiet {
				return msg
			}
			fmt.Fprintln(self.Writer, msg)
			return ""
		}
	}
	// This is probably where CommandIf detection would happen
	fmt.Fprintln(self.Writer, self.compose("", message...))
	return nil
}

var (
	caller_re = regexp.MustCompile(`/@(?::(\d+))?`)
	format_re = regexp.MustCompile(`/%`)
	quiet_re = regexp.MustCompile(`/<`)
)

type _dbgCommand struct {
	quiet bool
	caller int
	format bool
	message []interface{}
}

func (self *Debug) run(cmd _dbgCommand, skip int) string {
	output := []string{}

	{
		caller := cmd.caller
		if caller > 0 {
			caller += skip
			pc, _, _, _ := runtime.Caller(caller)
			name := runtime.FuncForPC(pc).Name()
			output = append(output, fmt.Sprintf("@ %s:", name))
		}
	}

	msg := ""
	if cmd.format {
		if len(cmd.message) > 0 {
			switch format := cmd.message[0].(type) {
			case string:
				msg = self.compose(format, cmd.message[1:]...)
			default:
				msg = self.compose("", cmd.message...)
			}
		}
	} else {
		msg = self.compose("", cmd.message...)
	}

	if msg != "" {
		output = append(output, msg)
	}
	return strings.Join(output, " ")
}

func (self *Debug) command(command string, message ...interface{}) _dbgCommand {

	quiet := false
	caller := -1
	format := false

	if len(command) > 0 {
		if command[0] == '<' {
			quiet = true
			command = command[1:]
		}
		if match := caller_re.FindStringSubmatch(command); match != nil {
			caller = 1
			if match[1] != "" {
				tmp, _ := strconv.ParseInt(match[1], 10, 32)
				caller = int(tmp)
			}
		}
		if format_re.MatchString(command) {
			format = true
		}
		if quiet_re.MatchString(command) {
			quiet = true
		}
	}

	return _dbgCommand{
		quiet: quiet,
		caller: caller,
		format: format,
		message: message,
	}
}


func (self *Debug) compose(format string, message ...interface{}) string {
	if format == "" {
		if len(message) == 0 {
			return ""
		}

		if self.FormatIf != nil {
			switch value := message[0].(type) {
			case string:
				if match := self.FormatIf.FindStringIndex(value); match != nil {
					format := value[match[1]:]
					message = message[1:]
					return fmt.Sprintf(format, message...)
				}
			}
		}

		// No (Sprintf) formatting done
		output := []string{}
		for _, argument := range message {
			output = append(output, self.toString(argument))
		}
		return strings.Join(output, " ")
	}
	return fmt.Sprintf(format, message...)
}

func (self *Debug) toString(value interface{}) string {
	return fmt.Sprintf(self.ToStringFormat, value)
}

