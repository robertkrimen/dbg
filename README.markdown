# dbg
--
    import "github.com/robertkrimen/dbg"

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

## Usage

```go
var (
	// The writer to send Dbg output to.
	Writer io.Writer = os.Stderr

	// The fmt "verb" to use when formatting a value for output.
	ToStringFormat string = "%v"

	// The pattern to look for when detecting whether to use the
	// input as arguments to fmt.Sprintf, with the first argument 
	// being the format string.
	// Assumed to be a string at the start of the first argument.
	FormatIf = regexp.MustCompile(`^\\:`)
)
```

#### func  Dbg

```go
func Dbg(input ...interface{})
```
Dbg will output the given input to Writer, applying any special formatting if
necessary.

#### func  Dbgf

```go
func Dbgf(input ...interface{})
```
Dbg will output the given input to Writer via fmt.Sprintf, with the first
argument being the format string.

#### func  Fmt

```go
func Fmt(input ...interface{}) string
```
Fmt will process the input via fmt.Sprintf, with the first argument being the
format string and return the result.

#### func  HereBeDragons

```go
func HereBeDragons(input ...interface{}) string
```
HereBeDragons returns a string formatted like Dbg, but prefixed with the call
location (function name).

--
**godocdown** http://github.com/robertkrimen/godocdown
