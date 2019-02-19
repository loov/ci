package ci

// Copy copies from source directory to destination directory
type Copy struct {
	SourceGlob  string
	Destination string
}

// Setup sets up the step
func (step *Copy) Setup(parent *Task) {
	task := parent.Subtask("cp %q %q", step.SourceGlob, step.Destination)
	task.Exec = func(context, _ *Context) error {
		// TODO: verify source and destination is inside context.Global.WorkDir
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
		// TODO: verify glob is inside context.Global.WorkDir
		return nil
	}
}
