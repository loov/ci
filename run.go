package ci

import (
	"time"
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
		// TODO:
		time.Sleep(time.Second)
		return nil
	}
}
