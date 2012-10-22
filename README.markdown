# dbg
--
    import "github.com/robertkrimen/dbg"

Package dbg is a println/printf/log-debugging utility library.

    import (
        Dbg "github.com/robertkrimen/dbg"
    )

    dbg, dbgf := Dbg.New()

    dbg("Emit some debug stuff", []byte{120, 121, 122, 122, 121}, math.Pi)
    # "2013/01/28 16:50:03 Emit some debug stuff [120 121 122 122 121] 3.141592653589793"

    dbgf("With a %s formatting %.2f", "little", math.Pi)
    # "2013/01/28 16:51:55 With a little formatting (3.14)"

    dbgf("%/fatal//A fatal debug statement: should not be here")
    # "A fatal debug statement: should not be here"
    # ...and then, os.Exit(1)

    dbgf("%/panic//Can also panic %s", "this")
    # "Can also panic this"
    # ...and then, os.Exit(1)

    dbgf("Any %s arguments without a corresponding %%", "extra", "are treated like arguments to dbg()")
    # "2013/01/28 17:14:40 Any extra arguments (without a corresponding %) are treated like arguments to dbg()"

    dbgf("%d %d", 1, 2, 3, 4, 5)
    # "2013/01/28 17:16:32 Another example: 1 2 3 4 5"

    dbgf("%@: Include the function name for a little context (via %s)", "%@")
    # "2013/01/28 17:18:56 github.com/robertkrimen/dbg.TestSynopsis: Include the function name for a little context (via %@)"

By default, dbg uses log (log.Println, log.Printf, log.Panic, etc.) for output.
However, you can also provide your own output destination by supplying a invoking dbg.New
with a customization function:

    import (
        "bytes"
        Dbg "github.com/robertkrimen/dbg"
        "os"
    )

    # dbg to os.Stderr
    dbg, dbgf := Dbg.New(func(dbgr *Dbgr) {
        dbgr.SetOutput(os.Stderr)
    })

    # A slightly contrived example:
    var buffer bytes.Buffer
    dbg, dbgf := New(func(dbgr *Dbgr) {
        dbgr.SetOutput(&buffer)
    })

## Usage

#### type DbgFunction

```go
type DbgFunction func(values ...interface{})
```


#### func  New

```go
func New(options ...interface{}) (dbg DbgFunction, dbgf DbgFunction)
```
New will create and return a pair of debugging functions. You can customize
where they output to by passing in an (optional) customization function:

    import (
        Dbg "github.com/robertkrimen/dbg"
        "os"
    )
    # dbg to os.Stderr
    dbg, dbgf := Dbg.New(func(dbgr *Dbgr) {
        dbgr.SetOutput(os.Stderr)
    })

#### type Dbgr

```go
type Dbgr struct {
}
```


#### func  NewDbgr

```go
func NewDbgr() *Dbgr
```

#### func (Dbgr) Dbg

```go
func (self Dbgr) Dbg(values ...interface{})
```

#### func (Dbgr) DbgDbgf

```go
func (self Dbgr) DbgDbgf() (dbg DbgFunction, dbgf DbgFunction)
```

#### func (Dbgr) Dbgf

```go
func (self Dbgr) Dbgf(values ...interface{})
```

#### func (*Dbgr) SetOutput

```go
func (self *Dbgr) SetOutput(output interface{})
```
SetOutput will accept a *log.Logger or io.Writer as a destination for output.
Passing in nil will reset output to go to the standard logger from the log
package.

--
**godocdown** http://github.com/robertkrimen/godocdown
