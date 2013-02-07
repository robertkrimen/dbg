# dbg-import
--
Command dbg-import will smuggol "github.com/robertkrimen/dbg" into a local Go package/repository

### Install

    go get github.com/robertkrimen/dbg-import

### Usage

    Usage: dbg-import [target]
     -quiet=false: Be absolutely quiet
     -update=false: Update (go get -u) package first
     -verbose=false: Be more verbose

        # Import "github.com/robertkrimen/dbg" into the current directory
        $ dbg-import

        # Import "github.com/robertkrimen/dbg" into another directory
        $ dbg-import ./xyzzy

--
**godocdown** http://github.com/robertkrimen/godocdown
