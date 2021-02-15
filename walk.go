package main

import (
	"os"
	"path/filepath"
	"sync"
)

type walker struct {
	ch chan string
	wg sync.WaitGroup
}

func (w *walker) start(concurrency int) {
	w.wg.Add(concurrency)
	w.ch = make(chan string, concurrency*2)
	for i := 0; i < concurrency; i++ {
		go w.worker()
	}
}

func (w *walker) callback(path string, info os.FileInfo, err error) error {
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

	w.ch <- path
	return nil
}

func (w *walker) worker() {
	defer w.wg.Done()

	for path := range w.ch {
		checkError(setModTime(path))
	}
}

func (w *walker) finish() {
	close(w.ch)
	w.wg.Wait()
}
