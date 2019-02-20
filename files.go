package ci

import (
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
		source, err := context.ExpandEnv(step.SourceGlob)
		if err != nil {
			return err
		}

		destination, err := context.ExpandEnv(step.Destination)
		if err != nil {
			return err
		}

		destinationErr := context.Global.SafeGlob(destination)
		if destinationErr != nil {
			return destinationErr
		}

		sourceErr := context.Global.SafeGlob(source)
		if sourceErr != nil {
			return sourceErr
		}

		matches, err := filepath.Glob(source)
		if err != nil {
			return err
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
		glob, err := context.ExpandEnv(step.Glob)
		if err != nil {
			return err
		}

		err = context.Global.SafeGlob(glob)
		if err != nil {
			return err
		}

		// TODO: verify glob is inside
		return nil
	}
}
