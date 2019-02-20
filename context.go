package ci

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

type Logger interface {
	Named(name string) Logger
	Output() (stdout, stderr io.Writer)

	Print(v ...interface{})
	Printf(v ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})
}

// GlobalContext defines the global execution context
type GlobalContext struct {
	// ScriptDir is the script location
	ScriptDir string
	Context

	// TempDir defines the temporary working directory
	temp struct {
		init  sync.Once
		root  string
		def   string
		index int32
	}
}

// Context defines task execution context and environment variable management
type Context struct {
	Global     *GlobalContext
	WorkingDir string
	Env        []string

	Logger
}

// NewGlobalContext creates
func NewGlobalContext(logger Logger) (*GlobalContext, error) {
	context := &GlobalContext{}
	context.Global = context
	context.Logger = logger
	context.Env = os.Environ()

	err := context.init()
	if err != nil {
		return nil, err
	}

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
	if err := os.Mkdir(context.temp.root, 0777); err != nil {
		return err
	}

	// create default temporary directory for commands
	context.temp.def = filepath.Join(context.temp.root, "temp")
	if err := os.Mkdir(context.temp.def, 0777); err != nil {
		return err
	}

	return nil
}

// CreateTempDir creates a temporary directory
func (context *GlobalContext) CreateTempDir(prefix string) string {
	index := atomic.AddInt32(&context.temp.index, 1)
	dir := filepath.Join(context.temp.root, prefix+"-"+strconv.Itoa(int(index)))
	if err := os.Mkdir(dir, 0777); err != nil {
		context.Errorf("failed to create nested temporary directory: %v", err)
	}
	return ""
}

// Cleanup deletes all temporary data.
func (context *GlobalContext) Cleanup() error {
	return os.RemoveAll(context.temp.root)
}

// Clone creates a clone of the context
func (context *Context) Clone() *Context {
	return &Context{
		Global:     context.Global,
		WorkingDir: context.WorkingDir,
		Env:        append([]string{}, context.Env...),
	}
}

// SetEnv changes environment variable value
func (context *Context) SetEnv(env, value string) {
	_ = context.UnsetEnv(env)
	context.Env = append(context.Env, env+"="+value)
}

// UnsetEnv removes an existing environment value
func (context *Context) UnsetEnv(env string) bool {
	for i, env := range context.Env {
		eq := strings.Index(env, "=")
		if eq < 0 {
			continue
		}
		if strings.EqualFold(env[:eq], env) {
			context.Env = append(context.Env[:i], context.Env[i+1:]...)
			return true
		}
	}
	return false
}

// GetEnv finds the value of an environment variable
func (context *Context) GetEnv(env string) (string, bool) {
	for _, env := range context.Env {
		eq := strings.Index(env, "=")
		if eq < 0 {
			continue
		}
		if strings.EqualFold(env[:eq], env) {
			return env[eq+1:], true
		}
	}
	return "", false
}

// ExpandEnv replaces enviroment values in value,
// returns an error when it is missing
func (context *Context) ExpandEnv(value string) (string, error) {
	var missing []string
	expanded := os.Expand(value, func(env string) string {
		value, ok := context.GetEnv(env)
		if !ok {
			missing = append(missing, value)
		}
		return value
	})
	if len(missing) > 0 {
		return expanded, fmt.Errorf("missing variables: %v", missing)
	}
	return expanded, nil
}
