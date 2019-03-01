package ci

import (
	"fmt"
	"os"
	"path/filepath"
)

// Copy copies from source directory to destination directory
type Copy struct {
	SourceGlob  string
	Destination string
}

// Setup sets up the step
func (step *Copy) Setup(parent *Task) {
	task := parent.Subtask("cp %q %q", step.SourceGlob, step.Destination)
	task.Exec = func(context, _ *Context) error {
		source, sourcePrefix, err := context.AbsGlob(step.SourceGlob)
		if err != nil {
			return err
		}

		destination, destinationPrefix, err := context.AbsGlob(step.Destination)
		if err != nil {
			return err
		}
		if destination != destinationPrefix {
			return fmt.Errorf("glob not allowed in destination %q [expanded %q]", step.Destination, destination)
		}

		if err := os.MkdirAll(destination, 0777); err != nil && !os.IsExist(err) {
			return err
		}

		matches, err := filepath.Glob(source)
		if err != nil {
			return err
		}

		for _, match := range matches {
			err := copyAny(sourcePrefix, match, destination)
			if err != nil {
				return err
			}
		}

		// TODO: verify source and destination is inside
		return nil
	}
}

// Remove deletes the files matching a glob
type Remove struct {
	Glob string
}

// Setup sets up the step
func (step *Remove) Setup(parent *Task) {
	task := parent.Subtask("rm %q", step.Glob)
	task.Exec = func(context, _ *Context) error {
		glob, _, err := context.AbsGlob(step.Glob)
		if err != nil {
			return err
		}

		matches, err := filepath.Glob(glob)
		for _, match := range matches {
			if err := safeRemove(match); err != nil {
				return err
			}
		}

		return nil
	}
}

func isDir(path string) bool {
	_, err := os.Stat(path)
	return os.IsExist(err)
}

func safeRemove(path string) error {
	if filepath.Dir(path) == path {
		return fmt.Errorf("tried to delete %q", path)
	}
	return os.RemoveAll(path)
}

func copyAny(sourceDir, sourceDirOrFile, destinationDir string) error {
	return nil
}
