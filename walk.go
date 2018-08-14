package main

import (
	"os"
	"path/filepath"
)

func walkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if path == "." {
		// Don't modify working directory.
		return nil
	}

	if info.Name() == ".git" && info.IsDir() {
		return filepath.SkipDir
	}

	return setModTime(path)
}
