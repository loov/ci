package ci

import (
	"os"
	"os/exec"
	"strings"
)

// Run implements a step for executing a command
type Run struct {
	Command string
	Args    []string
}

// Setup sets up the step
func (run *Run) Setup(parent *Task) {
	task := parent.Subtask("run %q", run)
	task.Exec = func(_, subcontext *Context) error {
		subcontext.Logger.Printf("run %q\n", run)
		cmd := exec.Command(run.Command, run.Args...)
		cmd.Dir = subcontext.WorkingDir
		cmd.Env = subcontext.Env
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		return cmd.Run()
	}
}

// String returns string representation.
func (run *Run) String() string {
	var args []string
	args = append(args, run.Command)
	args = append(args, run.Args...)
	return strings.Join(args, " ")
}
