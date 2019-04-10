package ci

// TempGopath executes steps in a temporary gopath directory
type TempGopath struct {
	Steps []Step
}

// Setup sets up the step
func (step *TempGopath) Setup(parent *Task) {
	task := parent.Subtask("temp gopath")
	task.Exec = func(_, context *Context) error {
		dir := context.Global.CreateTempDir("gopath")
		context.SetEnv("GOPKG", "$GOPATH/pkg")
		context.SetEnv("GOPATH", dir)
		return nil
	}
	task.AddSteps(step.Steps)
}

// ChangeDir changes working directory
type ChangeDir struct {
	Target string
}

// Setup sets up the step
func (step *ChangeDir) Setup(parent *Task) {
	task := parent.Subtask("cd %q", step.Target)
	task.Exec = func(context, _ *Context) error {
		dir, err := context.ExpandEnv(step.Target)
		if err != nil {
			return nil
		}
		context.WorkingDir = dir
		return nil
	}
}
