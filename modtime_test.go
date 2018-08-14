package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGetModTime(t *testing.T) {
	mod, err := getModTime("README.md")
	if err != nil {
		t.Fatal(err)
	}

	if mod.IsZero() {
		t.Error("README.md was not found in Git.")
	}

	fi, err := os.Stat("README.md")
	if err != nil {
		t.Fatal(err)
	}

	if mod.After(fi.ModTime()) {
		t.Error("README.md was last committed before it was last modified.")
	}
}

func TestGetModTime_commit(t *testing.T) {
	*flagCommit = true
	defer func() {
		*flagCommit = false
	}()

	TestGetModTime(t)
}

func TestGetModTime_uncommitted(t *testing.T) {
	// The .git folder is never committed.
	mod, err := getModTime(".git")
	if err != nil {
		t.Fatal(err)
	}

	if !mod.IsZero() {
		t.Error(".git was found in Git.")
	}
}

func TestGetModTime_notgit(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err = os.RemoveAll(dir); err != nil {
			t.Error(err)
		}
	}()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(wd); err != nil {
			t.Error(err)
		}
	}()

	if err := os.Chdir(dir); err != nil {
		t.Error(err)
	}

	testStderr(t, func() {
		_, err := getModTime("somefile.txt")
		if err == nil {
			t.Error("Expected error, but no error was returned.")
		}
	}, "fatal: not a git repository (or any of the parent directories): .git\n")
}
