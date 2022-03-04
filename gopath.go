package ci

import (
	"os/exec"
	"strings"
)

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

func getGOPATH() (string, error) {
	out, err := exec.Command("go", "env", "GOPATH").CombinedOutput()
	return strings.TrimSpace(string(out)), err
}
