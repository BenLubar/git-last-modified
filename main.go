// Command git-last-modified sets files in a Git repository to their last-modified date.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var flagSet = flag.NewFlagSet("git-last-modified", flag.ContinueOnError)
var flagQuiet = flagSet.Bool("q", false, "Quiet. Don't warn about files specified on the command line that are not in Git.")
var flagVerbose = flagSet.Bool("v", false, "Verbose. Print each filename and modification time as they are processed.")
var flagDryRun = flagSet.Bool("n", false, "Dry run. Implies -v. Don't modify any file modification times.")
var flagCommit = flagSet.Bool("commit-date", false, "Use the commit date for the last commit this file was involved in instead of the author date.")

var exit = os.Exit

func usage() {
	_, _ = fmt.Fprintln(os.Stderr, "usage: git last-modified [<options>] [[--] <path>...]")
	flagSet.PrintDefaults()
	exit(2)
}

func checkError(err error) {
	if err != nil {
		if *flagVerbose {
			_, _ = fmt.Fprintf(os.Stderr, "git-last-modified: %+v\n", err)
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "git-last-modified: %v\n", err)
		}
		exit(1)
	}
}

func main() {
	flagSet.Usage = usage
	if err := flagSet.Parse(os.Args[1:]); err != nil {
		usage()
	}
	if *flagDryRun {
		*flagVerbose = true
	}

	files := flag.Args()
	if len(files) == 0 {
		*flagQuiet = true

		var w walker
		w.start(runtime.GOMAXPROCS(0))
		checkError(filepath.Walk(".", w.callback))
		w.finish()
	} else {
		for _, path := range files {
			checkError(setModTime(path))
		}
	}
}
