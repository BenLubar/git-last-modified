package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func getModTime(path string) (time.Time, error) {
	pretty := "--pretty=format:%aI"
	if *flagCommit {
		pretty = "--pretty=format:%cI"
	}

	/* #nosec G204 */
	cmd := exec.Command("git", "log", "-n", "1", pretty, "--", path)
	cmd.Stderr = os.Stderr

	b, err := cmd.Output()
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "get last-modified time for file %q", path)
	}

	s := strings.TrimSpace(string(b))
	if s == "" {
		return time.Time{}, nil
	}

	t, err := time.Parse(time.RFC3339, s)
	return t, errors.Wrapf(err, "cannot parse time %q for file %q", s, path)
}

func setModTime(path string) error {
	t, err := getModTime(path)
	if err != nil {
		return err
	}

	if t.IsZero() {
		if *flagVerbose {
			fmt.Printf("%s: never committed in Git\n", path)
		}

		if !*flagQuiet {
			return errors.Errorf("file not found in Git history: %q", path)
		}

		// Never committed in Git. Leave it alone.
		return nil
	}

	if *flagVerbose {
		fmt.Printf("%s: last modified %v\n", path, t)
	}

	if *flagDryRun {
		return nil
	}

	return errors.Wrapf(os.Chtimes(path, t, t), "cannot change time for file %q", path)
}
