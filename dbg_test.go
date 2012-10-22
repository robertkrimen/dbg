package dbg

import (
	"bytes"
	. "github.com/robertkrimen/terst"
	"testing"
)

func Test(t *testing.T) {
	Terst(t)

	Is(toString("Xyzzy"), "Xyzzy")
	Is(inputToString(1, 2.1, "3"), "1 2.1 3")
	Is(inputToString(1, 2.0, "3"), "1 2 3")
	Is(inputToString("\\:%d %.1f %s", 1, 2.0, "3"), "1 2.0 3")
	Is(inputToString("Hello, World.\n"), "Hello, World.\n")

	var buffer bytes.Buffer
	Writer = &buffer

	buffer.Reset()
	Dbg("Hello, World.")
	Is(buffer.String(), "Hello, World.\n")

	Dbg("Nothing happens.")
	Is(buffer.String(), "Hello, World.\nNothing happens.\n")

	buffer.Reset()
	Dbg("Hello, World.")
	Is(buffer.String(), "Hello, World.\n")

	Dbg("\\:When it gets very %s, it gets very %s indeed.", "loud", "LOUD")
	Is(buffer.String(), "Hello, World.\nWhen it gets very loud, it gets very LOUD indeed.\n")
}
