package ci

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
)

type Logger interface {
	Named(name string) Logger
	// Output() (stdout, stderr io.Writer)

	Print(v ...interface{})
	Printf(format string, v ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})
}

// GlobalContext defines the global execution context
type GlobalContext struct {
	// ScriptDir is the script location
	ScriptDir string
	// GEnv is the global environment variables
	GEnv Env

	Context

	// TempDir defines the temporary working directory
	temp struct {
		root  string
		def   string
		index int32
	}
}

// Sub creates a sub context
func (context *Context) Sub(name string) *Context {
	return &Context{
		Global:     context.Global,
		WorkingDir: context.WorkingDir,
		Env:        context.Env.Clone(),
		Logger:     context.Logger.Named(name),
	}
}

// Context defines task execution context and environment variable management
type Context struct {
	Global     *GlobalContext
	WorkingDir string
	Env        Env

	Logger
}

// NewGlobalContext creates
func NewGlobalContext(scriptDir string, logger Logger) (*GlobalContext, error) {
	context := &GlobalContext{}
	context.Global = context

	context.Logger = logger
	if context.Logger == nil {
		context.Logger = NewStd()
	}

	context.Env = os.Environ()

	err := context.init()
	if err != nil {
		return nil, err
	}

	absScriptDir, err := filepath.Abs(scriptDir)
	if err != nil {
		return nil, err
	}

	context.ScriptDir = absScriptDir
	context.SetEnv("SCRIPTDIR", context.ScriptDir)

	if runtime.GOOS == "windows" {
		context.SetEnv("TEMP", context.temp.def)
		context.SetEnv("TMP", context.temp.def)
	} else {
		context.SetEnv("TMPDIR", context.temp.def)
	}

	return context, err
}

func (context *GlobalContext) init() error {
	var err error

	// create root temp directory
	context.temp.root, err = ioutil.TempDir("", "ci")
	if err != nil {
		return err
	}

	// create default temporary directory for commands
	context.temp.def = filepath.Join(context.temp.root, "temp")
	if err := os.Mkdir(context.temp.def, 0777); err != nil {
		return err
	}

	return nil
}

// SafePath checks whether glob can be changed
func (context *GlobalContext) SafePath(glob string) error {
	return nil
}

// CreateTempDir creates a temporary directory
func (context *GlobalContext) CreateTempDir(prefix string) string {
	index := atomic.AddInt32(&context.temp.index, 1)
	dir := filepath.Join(context.temp.root, prefix+"-"+strconv.Itoa(int(index)))
	if err := os.Mkdir(dir, 0777); err != nil {
		context.Errorf("failed to create nested temporary directory: %v", err)
	}
	return dir
}

// Cleanup deletes all temporary data.
func (context *GlobalContext) Cleanup() error {
	return os.RemoveAll(context.temp.root)
}

// SetEnv changes environment variable value
func (context *Context) SetEnv(key, value string) {
	context.Env.Set(key, value)
}

// UnsetEnv removes an existing environment value
func (context *Context) UnsetEnv(target string) bool {
	return context.Env.Unset(target)
}

// GetEnv finds the value of an environment variable
func (context *Context) GetEnv(target string) (string, bool) {
	v, ok := context.Env.Get(target)
	if !ok {
		return context.Global.GEnv.Get(target)
	}
	return v, ok
}

// ExpandEnv replaces enviroment values in value,
// returns an error when it is missing
func (context *Context) ExpandEnv(value string) (string, error) {
	var missing []string
	expanded := os.Expand(value, func(env string) string {
		value, ok := context.GetEnv(env)
		if !ok {
			missing = append(missing, env)
		}
		return value
	})
	if len(missing) > 0 {
		return expanded, fmt.Errorf("missing variables: %v", missing)
	}
	return expanded, nil
}

// AbsGlob converts path with environment variables to a absolute path
func (context *Context) AbsGlob(value string) (abs string, absprefix string, err error) {
	expanded, err := context.ExpandEnv(value)
	if err != nil {
		return "", "", err
	}

	abs, err = filepath.Abs(expanded)
	if err != nil {
		return "", "", err
	}

	absprefix = extractGlobPrefix(abs)

	// TODO: verify absprefix

	return abs, absprefix, nil
}

func extractGlobPrefix(glob string) string {
	p := strings.IndexAny(glob, "?*")
	if p < 0 {
		return glob
	}
	return glob[:p]
}

// Env defines a set of environment variables
type Env []string

// Clone creates a deep clone of the environment.
func (env Env) Clone() Env { return append(Env{}, env...) }

// Set changes environment variable value
func (env *Env) Set(key, value string) {
	_ = env.Unset(key)
	*env = append(*env, key+"="+value)
}

// Unset removes an existing environment value
func (env *Env) Unset(target string) bool {
	for i, key := range *env {
		eq := strings.Index(key, "=")
		if eq < 0 {
			continue
		}
		if strings.EqualFold(key[:eq], target) {
			*env = append((*env)[:i], (*env)[i+1:]...)
			return true
		}
	}
	return false
}

// Get finds the value of an environment variable
func (env *Env) Get(target string) (string, bool) {
	for _, key := range *env {
		eq := strings.Index(key, "=")
		if eq < 0 {
			continue
		}
		if strings.EqualFold(key[:eq], target) {
			return key[eq+1:], true
		}
	}
	return "", false
}
