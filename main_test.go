package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

var stderr io.Writer = os.Stderr

func pipeStderr(t *testing.T, r *os.File, wait chan<- struct{}) {
	defer func() {
		close(wait)
		if err := r.Close(); err != nil {
			t.Error(err)
		}
	}()

	var buf [4096]byte
	for {
		n, err := r.Read(buf[:])
		if err != nil && err != io.EOF {
			t.Error(err)
			return
		}

		n1, err1 := stderr.Write(buf[:n])
		if err1 != nil {
			t.Error(err1)
			return
		}

		if n1 != n {
			t.Error(io.ErrShortWrite)
			return
		}

		if err == io.EOF {
			return
		}
	}
}

func testStderr(t *testing.T, f func(), expected string) {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	realStderr := os.Stderr
	os.Stderr = w
	wait := make(chan struct{})
	go pipeStderr(t, r, wait)

	var buf bytes.Buffer
	oldStderr := stderr
	stderr = &buf
	defer func() {
		if err := w.Close(); err != nil {
			t.Error(err)
		}
		<-wait

		os.Stderr = realStderr
		stderr = oldStderr
		if actual := strings.Replace(buf.String(), "\r\n", "\n", -1); !strings.EqualFold(actual, expected) {
			t.Error("Expected stderr does not match actual stderr.")
			t.Errorf("expected: %q", expected)
			t.Errorf("actual:   %q", actual)
		}
	}()

	f()
}

func testExit(t *testing.T, f func(), expectedCode int) {
	type exited int
	oldExit := exit
	exit = func(code int) {
		panic(exited(code))
	}
	defer func() {
		exit = oldExit
		if r := recover(); r == nil {
			t.Errorf("Expected call to exit(%d), but exit was not called.", expectedCode)
		} else if code, ok := r.(exited); ok {
			if int(code) != expectedCode {
				t.Errorf("Expected call to exit(%d), but exit(%d) was called.", expectedCode, int(code))
			}
		} else {
			panic(r)
		}
	}()

	f()
}

func TestUsage(t *testing.T) {
	testStderr(t, func() {
		testExit(t, usage, 2)
	}, "usage: git last-modified [<options>] [[--] <path>...]\n"+
		"  -commit-date\n"+
		"    \tUse the commit date for the last commit this file was involved in instead of the author date.\n"+
		"  -n\tDry run. Implies -v. Don't modify any file modification times.\n"+
		"  -q\tQuiet. Don't warn about files specified on the command line that are not in Git.\n"+
		"  -v\tVerbose. Print each filename and modification time as they are processed.\n")
}

func TestCheckError(t *testing.T) {
	testStderr(t, func() {
		checkError(nil) // should not exit
	}, "")

	err := errors.New("test")

	testStderr(t, func() {
		testExit(t, func() {
			checkError(err)
		}, 1)
	}, "git-last-modified: test\n")

	*flagVerbose = true
	defer func() {
		*flagVerbose = false
	}()

	testStderr(t, func() {
		testExit(t, func() {
			checkError(err)
		}, 1)
	}, fmt.Sprintf("git-last-modified: test%+v\n", err.(interface{ StackTrace() errors.StackTrace }).StackTrace()))
}
