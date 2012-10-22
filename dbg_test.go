package dbg

import (
	"bytes"
	. "github.com/robertkrimen/terst"
	"math"
	"testing"
)

func Test_parseFormat(t *testing.T) {
	Terst(t)

	frmt := _frmt{}
	Is(frmt.ctl, "")
	Is(frmt.format, "")

	frmt = parseFormat("")
	Is(frmt.ctl, "")
	Is(frmt.format, "")

	frmt = parseFormat("%/panic")
	Is(frmt.ctl, "%/panic")
	Is(frmt.format, "")
	Is(frmt.panic, true)
	Is(frmt.fatal, false)

	frmt = parseFormat(" %/fatal")
	Is(frmt.ctl, "%/fatal")
	Is(frmt.format, "")
	Is(frmt.panic, false)
	Is(frmt.fatal, true)

	frmt = parseFormat(" %/fatal//Xyzzy")
	Is(frmt.ctl, "%/fatal")
	Is(frmt.format, "Xyzzy")
	Is(frmt.panic, false)
	Is(frmt.fatal, true)

	frmt = parseFormat(" %/fatal// Nothing happens.")
	Is(frmt.ctl, "%/fatal")
	Is(frmt.format, " Nothing happens.")
	Is(frmt.panic, false)
	Is(frmt.fatal, true)

	frmt = parseFormat("   %/panic %/fatal//")
	Is(frmt.ctl, "%/panic %/fatal")
	Is(frmt.format, "")
	Is(frmt.panic, true)
	Is(frmt.fatal, true)
}

func Test_operandCount(t *testing.T) {
	Terst(t)

	test := func(target string, count int) {
		Is(operandCount(target), count)
	}

	test("", 0)
	test("%%", 0)
	test("%@", 0)
	test("%@ %%", 0)
	test("%%%", 0)
	test("%%@", 0)
	test("%d %s %%", 2)
	test("%d %- s %%", 2)
	test("%.2f", 1)
	test("Something extra at %s end (with %%) %.2f:", 2)
}

func Test(t *testing.T) {
	Terst(t)

	var buffer bytes.Buffer

	dbg, dbgf := New(func(dbgr *Dbgr) {
		dbgr.SetOutput(&buffer)
	})

	buffer.Reset()
	dbg("Hello, World.")
	Like(buffer.String(), "Hello, World.\n")

	dbg("Nothing happens.")
	Like(buffer.String(), "Hello, World.\nNothing happens.\n")

	buffer.Reset()
	dbg("Hello, World.")
	Is(buffer.String(), "Hello, World.\n")

	dbgf("When it gets very %s, it gets very %s indeed.", "loud", "LOUD")
	Is(buffer.String(), "Hello, World.\nWhen it gets very loud, it gets very LOUD indeed.\n")

	buffer.Reset()
	dbgf("Xyzzy", "Nothing happens.")
	Is(buffer.String(), "Xyzzy Nothing happens.\n")

	buffer.Reset()
	dbgf("Xyzzy (%s)", "Nothing happens.", 1)
	Is(buffer.String(), "Xyzzy (Nothing happens.) 1\n")

	tmp := "github.com/robertkrimen/dbg.Test"
	buffer.Reset()
	dbgf("%@ Xyzzy (%s)", "Nothing happens.", 1)
	Is(buffer.String(), tmp+" Xyzzy (Nothing happens.) 1\n")

	buffer.Reset()
	dbgf("%%@")
	Is(buffer.String(), "%@\n")

	buffer.Reset()
	dbgf("%%%@")
	Is(buffer.String(), "%"+tmp+"\n")

	buffer.Reset()
	dbgf("%@: Nothing happens.")
	Is(buffer.String(), tmp+": Nothing happens.\n")

	buffer.Reset()
	dbgf("Something extra at %s end:", "the", "Here.")
	Is(buffer.String(), "Something extra at the end: Here.\n")

	buffer.Reset()
	dbgf("Something extra at %s end: %.2f", "the", math.Pi, "+", 1)
	Is(buffer.String(), "Something extra at the end: 3.14 + 1\n")

	buffer.Reset()
	dbgf("Something extra at %s end (with %%): %.2f", "the", math.Pi, "+", 1)
	Is(buffer.String(), "Something extra at the end (with %): 3.14 + 1\n")
}

func TestSynopsis(t *testing.T) {
	Terst(t)

	var buffer bytes.Buffer

	dbg, dbgf := New(func(dbgr *Dbgr) {
		dbgr.SetOutput(&buffer)
	})

	buffer.Reset()
	dbg("Emit some debug stuff", []byte{120, 121, 122, 122, 121}, math.Pi)
	Is(buffer.String(), "Emit some debug stuff [120 121 122 122 121] 3.141592653589793\n")

	buffer.Reset()
	dbgf("With a %s formatting (%.2f)", "little", math.Pi)
	Is(buffer.String(), "With a little formatting (3.14)\n")

	buffer.Reset()
	dbgf("Any %s arguments (without a corresponding %%)", "extra", "are treated like arguments to dbg()")
	Is(buffer.String(), "Any extra arguments (without a corresponding %) are treated like arguments to dbg()\n")

	buffer.Reset()
	dbgf("Another example: %d %d", 1, 2, 3, 4, 5)
	Is(buffer.String(), "Another example: 1 2 3 4 5\n")

	buffer.Reset()
	dbgf("%@: Include the function name for a little context (via %s)", "%@")
	Is(buffer.String(), "github.com/robertkrimen/dbg.TestSynopsis: Include the function name for a little context (via %@)\n")
}
