package ci

import "strings"

// Pipelines defines a collection of pipelines
type Pipelines []*Pipeline

// Pipeline defines a single execution tree
type Pipeline struct {
	Name  string
	Desc  string
	Steps []Step
}

// Stage defines a set of steps to be executed
type Stage struct {
	Name     string
	Parallel bool
	Steps    []Step
}

// Step defines an operation that is done in the execution tree
type Step interface {
	// Setup creates the necessary subtasks
	Setup(parent *Task)
}

// Find finds a named pipeline.
func (pipelines Pipelines) Find(name string) (*Pipeline, bool) {
	for _, pipeline := range pipelines {
		if strings.EqualFold(pipeline.Name, name) {
			return pipeline, true
		}
	}
	return nil, false
}

// Task creates a root task from a pipeline
func (pipeline *Pipeline) Task() *Task {
	task := &Task{}
	task.Name = pipeline.Name
	task.Desc = pipeline.Desc
	task.AddSteps(pipeline.Steps)
	return task
}

// Setup sets up the step
func (stage *Stage) Setup(parent *Task) {
	task := parent.Subtask(stage.Name)
	task.Parallel = stage.Parallel
	task.AddSteps(stage.Steps)
}
