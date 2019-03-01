package dsl

import "github.com/loov/ci"

func Pipelines(pipelines ...*ci.Pipeline) ci.Pipelines {
	return pipelines
}

func Pipeline(name, desc string, steps ...ci.Step) *ci.Pipeline {
	return &ci.Pipeline{
		Name:  name,
		Desc:  desc,
		Steps: steps,
	}
}

func Stage(name string, steps ...ci.Step) *ci.Stage {
	return &ci.Stage{
		Name:     name,
		Parallel: false,
		Steps:    steps,
	}
}

func Parallel(name string, steps ...ci.Step) *ci.Stage {
	return &ci.Stage{
		Name:     name,
		Parallel: true,
		Steps:    steps,
	}
}

func Run(command string, args ...string) *ci.Run {
	return &ci.Run{Command: command, Args: args}
}

func SetEnv(name, value string) *ci.SetEnv {
	return &ci.SetEnv{
		Global: false,
		Env:    name,
		Value:  value,
	}
}

func SetGlobalEnv(name, value string) *ci.SetEnv {
	return &ci.SetEnv{
		Global: true,
		Env:    name,
		Value:  value,
	}
}

func WhenEnv(name, value string, steps ...ci.Step) *ci.WhenEnv {
	return &ci.WhenEnv{
		Env:   name,
		Value: value,
		Steps: steps,
	}
}

func WhenEnvSet(name string, steps ...ci.Step) *ci.WhenEnvSet {
	return &ci.WhenEnvSet{
		Env:   name,
		Steps: steps,
	}
}

func Copy(sourceGlob, destination string) *ci.Copy {
	return &ci.Copy{
		SourceGlob:  sourceGlob,
		Destination: destination,
	}
}

func Remove(glob string) *ci.Remove {
	return &ci.Remove{
		Glob: glob,
	}
}

func CD(target string) *ci.ChangeDir {
	return &ci.ChangeDir{Target: target}
}

func TempGopath(steps ...ci.Step) *ci.TempGopath {
	return &ci.TempGopath{Steps: steps}
}

func CreateGlobalTempDir(name string) *ci.CreateTempDir {
	return &ci.CreateTempDir{
		Global: true,
		Env:    name,
	}
}
