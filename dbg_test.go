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

func TestDbg_dbg(t *testing.T) {
	Terst(t)

	Is(Dbg()("<", "Hello, World."), "Hello, World.")
	Is(Dbg()("</@", "Hello, World."), "@ github.com/robertkrimen/dbg.TestDbg_dbg: Hello, World.")
}

func Test_process(t *testing.T) {
	Terst(t)

	dbg := _dbg{}

	process := dbg.process("")
	Is(process.quiet, false)
	Is(process.caller, -1)
	Is(process.format, false)
	Is(process.compose(0), "")

	process = dbg.process("<")
	Is(process.quiet, true)
	Is(process.caller, -1)
	Is(process.format, false)
	Is(process.compose(0), "")

	process = dbg.process("</@")
	Is(process.quiet, true)
	Is(process.caller, 1)
	Is(process.format, false)
	Is(process.compose(0), "@ github.com/robertkrimen/dbg.Test_process:")

	process = dbg.process("</%")
	Is(process.quiet, true)
	Is(process.caller, -1)
	Is(process.format, true)
	Is(process.compose(0), "")

	process = dbg.process("/@:2 /%", "")
	Is(process.quiet, false)
	Is(process.caller, 2)
	Is(process.format, true)
	Is(len(process.message), 1)
	Is(process.compose(0), "@ testing.tRunner:")

	process = dbg.process("/% /< /@")
	Is(process.quiet, true)
	Is(process.caller, 1)
	Is(process.format, true)
	Is(process.compose(0), "@ github.com/robertkrimen/dbg.Test_process:")

	process = dbg.process("/% /< /@", 1, true, false, 1.1)
	Is(process.quiet, true)
	Is(process.caller, 1)
	Is(process.format, true)
	Is(process.compose(0), "@ github.com/robertkrimen/dbg.Test_process: 1 true false 1.1")

	process = dbg.process("</%", "%s: %02d", "Xyzzy", 2)
	Is(process.quiet, true)
	Is(process.caller, -1)
	Is(process.format, true)
	Is(process.compose(0), "Xyzzy: 02")

	process = dbg.process("<", "%s: %02d", "Xyzzy", 2)
	Is(process.quiet, true)
	Is(process.caller, -1)
	Is(process.format, false)
	Is(process.compose(0), "%s: %02d Xyzzy 2")
}
