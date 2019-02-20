package ci

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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
	return context, context.init()
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
