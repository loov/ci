package ci

import "regexp"

// GlobalContext defines the global execution context
type GlobalContext struct {
	// ScriptDir is the script location
	ScriptDir string
	// TempDir defines the temporary working directory
	TempDir string

	Context
}

// Context defines task execution context and environment variable management
type Context struct {
	Global     *GlobalContext
	WorkingDir string
	Env        []string
}

// CreateTempDir creates a temporary directory
func (context *GlobalContext) CreateTempDir() string {
	// TODO:
	return ""
}

// Clone creates a clone of the context
func (context *Context) Clone() *Context {
	return &Context{
		Global:     context.Global,
		WorkingDir: context.WorkingDir,
		Env:        append([]string{}, context.Env...),
	}
}

// rxEnv matches any environment variable
var rxEnv = regexp.MustCompile(`\$[a-zA-Z0-9_]+`)

// SetEnv changes environment variable value
func (context *Context) SetEnv(env, value string) {
	// TODO:
}

// GetEnv finds the value of an environment variable
func (context *Context) GetEnv(env string) string {
	// TODO:
	return ""
}

// EvalEnv replaces enviroment values in value,
// returns an error when it is missing
func (context *Context) EvalEnv(value string) (string, error) {
	// TODO:
	return value, nil
}
