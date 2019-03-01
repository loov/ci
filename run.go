package ci

import (
	"os"
	"os/exec"
)

// Run implements a step for executing a command
type Run struct {
	Command string
	Args    []string
}

// Setup sets up the step
func (run *Run) Setup(parent *Task) {
	task := parent.Subtask("run %q", run.Command)
	task.Exec = func(_, subcontext *Context) error {
		subcontext.Logger.Printf("run %q %q\n", run.Command, run.Args)

		cmd := exec.Command(run.Command, run.Args...)
		cmd.Dir = subcontext.WorkingDir
		cmd.Env = subcontext.Env
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		return cmd.Run()
	}
}
