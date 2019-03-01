package ci

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// ErrSkip is used to skip a particular task, without terminating execution
var ErrSkip = errors.New("skip")

// TaskStatus defines the status for the Task.
type TaskStatus struct {
	Started  time.Time
	Finished time.Time

	Running   bool
	Skipped   bool
	Done      bool
	Errored   bool
	ExecError error

	Stderr bytes.Buffer
	Stdout bytes.Buffer
}

// Task defines the execution tree.
type Task struct {
	Name     string
	Desc     string
	Parallel bool

	// Exec is executed before Tasks,
	// where context is the callers context and
	// subcontext is the context used to execute sub tasks
	Exec  func(context, subcontext *Context) error
	Tasks []*Task

	mu     sync.Mutex
	status TaskStatus
}

// Start marks this task as started.
func (status *TaskStatus) Start() {
	status.Running = true
	status.Started = time.Now()
}

// Finish marks this status to be finished.
func (status *TaskStatus) Finish() {
	status.Running = false
	status.Done = true
	status.Finished = time.Now()
}

// Skip marks this task to be skipped.
func (status *TaskStatus) Skip() {
	status.Running = false
	status.Skipped = true
}

// Subtask creates a new subtask with a name
func (task *Task) Subtask(name string, args ...interface{}) *Task {
	subtask := &Task{
		Name: fmt.Sprintf(name, args...),
	}
	task.Tasks = append(task.Tasks, subtask)
	return subtask
}

// AddSteps sets up steps with this task as parent
func (task *Task) AddSteps(steps []Step) {
	for _, step := range steps {
		step.Setup(task)
	}
}

// Run executes the given task
func (task *Task) Run(context *Context) (err error) {
	task.updateStatus((*TaskStatus).Start)
	defer task.updateStatus((*TaskStatus).Finish)
	defer task.updateStatus(func(status *TaskStatus) { status.Errored = err != nil })

	subcontext := context.Sub(task.Name)
	if task.Exec != nil {
		err := task.Exec(context, subcontext)
		if err == ErrSkip {
			task.updateStatus((*TaskStatus).Skip)
			return nil
		}
		if err != nil {
			task.updateStatus(func(status *TaskStatus) {
				status.ExecError = err
			})
			return err
		}
	}

	if !task.Parallel {
		for _, subtask := range task.Tasks {
			err := subtask.Run(subcontext)
			if err != nil {
				return err
			}
		}
		return nil
	} else {
		var group errgroup.Group
		for _, subtask := range task.Tasks {
			subtask := subtask
			group.Go(func() error {
				return subtask.Run(subcontext)
			})
		}
		return group.Wait()
	}
}

func (task *Task) updateStatus(fn func(*TaskStatus)) {
	task.mu.Lock()
	defer task.mu.Unlock()
	fn(&task.status)
}

// Status reads the current task status.
func (task *Task) Status() TaskStatus {
	task.mu.Lock()
	defer task.mu.Unlock()
	return task.status
}

// PrintTo prints the execution tree
func (task *Task) PrintTo(w io.Writer, ident string) {
	status := task.Status()

	stat := " "
	duration := ""
	switch {
	case status.Running:
		stat = "R"
		duration = formatDuration(time.Since(status.Started))
	case status.Skipped:
		stat = "S"
		duration = formatDuration(status.Finished.Sub(status.Started))
	case status.Errored:
		stat = "E"
		duration = formatDuration(status.Finished.Sub(status.Started))
	case status.Done:
		stat = "+"
		duration = formatDuration(status.Finished.Sub(status.Started))
	}

	if len(task.Tasks) == 0 {
		fmt.Fprintf(w, "%5s [%s] %s%s\n", duration, stat, ident, task.Name)
		return
	}
	if task.Name != "" {
		var desc string
		if task.Desc != "" {
			desc = " " + task.Desc
		}

		if task.Parallel {
			fmt.Fprintf(w, "%5s [%s] %s%s:%s (parallel)\n", duration, stat, ident, task.Name, desc)
		} else {
			fmt.Fprintf(w, "%5s [%s] %s%s:%s\n", duration, stat, ident, task.Name, desc)
		}
	}
	for _, task := range task.Tasks {
		task.PrintTo(w, ident+"    ")
	}
}

func formatDuration(d time.Duration) string {
	return d.Truncate(time.Second).String()
}
