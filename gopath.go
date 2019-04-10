package ci

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// TempGopath executes steps in a temporary gopath directory
type TempGopath struct {
	Steps []Step
}

// Setup sets up the step
func (step *TempGopath) Setup(parent *Task) {
	task := parent.Subtask("temp gopath")
	task.Exec = func(_, context *Context) error {
		dir := context.Global.CreateTempDir("gopath")
		context.Printf("temp GOPATH := %q", dir)

		// share gopkg directory to share cache between steps
		gopkg, err := context.ExpandEnv("$GOPATH/pkg")
		if err != nil {
			gopath, err := getGOPATH()
			if err != nil {
				return err
			}
			gopkg = filepath.Join(gopath, "pkg")
		}

		err = os.Symlink(gopkg, filepath.Join(dir, "pkg"))
		if err != nil {
			return err
		}

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

func getGOPATH() (string, error) {
	out, err := exec.Command("go", "env", "GOPATH").CombinedOutput()
	return strings.TrimSpace(string(out)), err
}
