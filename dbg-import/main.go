/*
Command dbg-import will smuggol "github.com/robertkrimen/dbg" into a local Go package/repository

Install

    go get github.com/robertkrimen/dbg-import

Usage

    Usage: dbg-import [target]
     -quiet=false: Be absolutely quiet
     -update=false: Update (go get -u) package first
     -verbose=false: Be more verbose

        # Import "github.com/robertkrimen/dbg" into the current directory
        $ dbg-import

        # Import "github.com/robertkrimen/dbg" into another directory
        $ dbg-import ./xyzzy

*/
package main

import (
	_ "github.com/robertkrimen/dbg"
	"github.com/robertkrimen/smuggol"
)

func main() {
	smuggol.Main("dbg-import", "github.com/robertkrimen/dbg", map[string]string{
		"dbg.go": `
package {{ .HostPackage }}

import (
    Dbg "{{ .ImportPath }}"
)

var dbg, dbgf = Dbg.New()
`,
	})
}
