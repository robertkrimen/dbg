package dbg

import (
	"bytes"
	. "github.com/robertkrimen/terst"
	"testing"
)

func Test(t *testing.T) {
	Terst(t)

	//Is(toString("Xyzzy"), "Xyzzy")
	//Is(inputToString(1, 2.1, "3"), "1 2.1 3")
	//Is(inputToString(1, 2.0, "3"), "1 2 3")
	//Is(inputToString("\\:%d %.1f %s", 1, 2.0, "3"), "1 2.0 3")
	//Is(inputToString("Hello, World.\n"), "Hello, World.\n")

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

func TestDbgFn(t *testing.T) {
	Terst(t)

	Is(Dbg()("<", "Hello, World."), "Hello, World.")
	Is(Dbg()("</@", "Hello, World."), "@ github.com/robertkrimen/dbg.TestDbgFn: Hello, World.")
}

func Test_run(t *testing.T) {
	Terst(t)

	dbg := _debug()

	cmd := dbg.command("")
	Is(cmd.quiet, false)
	Is(cmd.caller, -1)
	Is(cmd.format, false)
	Is(dbg.run(cmd, 0), "")

	cmd = dbg.command("<")
	Is(cmd.quiet, true)
	Is(cmd.caller, -1)
	Is(cmd.format, false)
	Is(dbg.run(cmd, 0), "")

	cmd = dbg.command("</@")
	Is(cmd.quiet, true)
	Is(cmd.caller, 1)
	Is(cmd.format, false)
	Is(dbg.run(cmd, 0), "@ github.com/robertkrimen/dbg.Test_run:")

	cmd = dbg.command("</%")
	Is(cmd.quiet, true)
	Is(cmd.caller, -1)
	Is(cmd.format, true)
	Is(dbg.run(cmd, 0), "")

	cmd = dbg.command("/@:2 /%", "")
	Is(cmd.quiet, false)
	Is(cmd.caller, 2)
	Is(cmd.format, true)
	Is(len(cmd.message), 1)
	Is(dbg.run(cmd, 0), "@ testing.tRunner:")

	cmd = dbg.command("/% /< /@")
	Is(cmd.quiet, true)
	Is(cmd.caller, 1)
	Is(cmd.format, true)
	Is(dbg.run(cmd, 0), "@ github.com/robertkrimen/dbg.Test_run:")

	cmd = dbg.command("/% /< /@", 1, true, false, 1.1)
	Is(cmd.quiet, true)
	Is(cmd.caller, 1)
	Is(cmd.format, true)
	Is(dbg.run(cmd, 0), "@ github.com/robertkrimen/dbg.Test_run: 1 true false 1.1")

	cmd = dbg.command("</%", "%s: %02d", "Xyzzy", 2)
	Is(cmd.quiet, true)
	Is(cmd.caller, -1)
	Is(cmd.format, true)
	Is(dbg.run(cmd, 0), "Xyzzy: 02")

	cmd = dbg.command("<", "%s: %02d", "Xyzzy", 2)
	Is(cmd.quiet, true)
	Is(cmd.caller, -1)
	Is(cmd.format, false)
	Is(dbg.run(cmd, 0), "%s: %02d Xyzzy 2")
}
