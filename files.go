package ci

import (
	"fmt"
	"io"
	"io/ioutil"
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
			rel, err := filepath.Rel(sourcePrefix, match)
			if err != nil {
				return err
			}

			err = copyAny(sourcePrefix, rel, destination)
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

func copyAny(sourceDir, sourceDirOrFile, destinationDir string) (err error) {
	os.Mkdir(destinationDir, 0755)

	sourcePath := filepath.Join(sourceDir, sourceDirOrFile)
	destinationPath := filepath.Join(destinationDir, sourceDirOrFile)

	stat, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		paths, err := ioutil.ReadDir(sourcePath)
		if err != nil {
			return err
		}

		for _, path := range paths {
			name := path.Name()
			err := copyAny(sourcePath, name, destinationPath)
			if err != nil {
				return err
			}
		}
		return nil
	}

	source, err := os.Open(sourcePath)
	if err != nil {
		return err
	}

	defer source.Close()

	destination, err := os.Create(destinationPath)
	if err != nil {
		return err
	}

	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err == nil {
		sourceinfo, err := source.Stat()
		if err == nil {
			err = destination.Chmod(sourceinfo.Mode())
		}
	}

	return err
}
