package ci

// SetEnv changes the environment variable
type SetEnv struct {
	Global bool
	Env    string
	Value  string
}

// Setup sets up the step
func (step *SetEnv) Setup(parent *Task) {
	task := parent.Subtask("%v := %q", step.Env, step.Value)
	task.Exec = func(context, _ *Context) error {
		value, err := context.ExpandEnv(step.Value)
		if err != nil {
			return err
		}

		if step.Global {
			context.Global.GEnv.Set(step.Env, value)
		} else {
			context.SetEnv(step.Env, value)
		}

		return nil
	}
}

// WhenEnv executes only when the given environment variable matches the value
type WhenEnv struct {
	Env   string
	Value string
	Steps []Step
}

// Setup sets up the step
func (step *WhenEnv) Setup(parent *Task) {
	task := parent.Subtask("when %v == %q", step.Env, step.Value)
	task.Exec = func(context, _ *Context) error {
		value, err := context.ExpandEnv(step.Value)
		if err != nil {
			return err
		}

		current, _ := context.GetEnv(step.Env)
		if current != value {
			return ErrSkip
		}
		return nil
	}
	task.AddSteps(step.Steps)
}

// WhenEnvSet executes only when the given environment variable is set
type WhenEnvSet struct {
	Env   string
	Steps []Step
}

// Setup sets up the step
func (step *WhenEnvSet) Setup(parent *Task) {
	task := parent.Subtask("when %v", step.Env)
	task.Exec = func(context, _ *Context) error {
		current, _ := context.GetEnv(step.Env)
		if current == "" {
			return ErrSkip
		}
		return nil
	}
	task.AddSteps(step.Steps)
}

// CreateTempDir creates a temporary directory and sets it to environment variable
type CreateTempDir struct {
	Global bool
	Env    string
}

// Setup sets up the step
func (step *CreateTempDir) Setup(parent *Task) {
	task := parent.Subtask("%v := tempdir", step.Env)
	task.Exec = func(context, _ *Context) error {
		dir := context.Global.CreateTempDir(step.Env)
		if step.Global {
			context.Global.GEnv.Set(step.Env, dir)
		} else {
			context.SetEnv(step.Env, dir)
		}
		return nil
	}
}

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
