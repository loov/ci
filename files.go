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
		destinationErr := context.Global.SafeGlob(step.Destination)
		if destinationErr != nil {
			return destinationErr
		}

		sourceErr := context.Global.SafeGlob(step.SourceGlob)
		if sourceErr != nil {
			return sourceErr
		}

		_, err := filepath.Glob(step.SourceGlob)
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
		err := context.Global.SafeGlob(step.Glob)
		if err != nil {
			return err
		}

		// TODO: verify glob is inside
		return nil
	}
}
