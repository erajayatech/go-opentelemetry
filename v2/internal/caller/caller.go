package caller

import (
	"runtime"
	"strings"

	"github.com/erajayatech/go-opentelemetry/v2/internal/config"
)

// OptFuncName represents options for the FuncName function.
type OptFuncName struct {
	Skip int
}

// OptionFuncName represents an option for the FuncName function.
type OptionFuncName func(*OptFuncName)

const defaultSkip = 1

// FuncName return caller function name.
func FuncName(options ...OptionFuncName) string {
	option := &OptFuncName{
		Skip: defaultSkip,
	}
	for _, opt := range options {
		opt(option)
	}

	pc, _, _, ok := runtime.Caller(option.Skip)
	if !ok {
		return "?Caller?"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "?FuncForPC?"
	}

	isShort, _ := config.GetOtelSpanNameShort()
	if isShort {
		fnName := strings.Split(fn.Name(), "/")
		return fnName[len(fnName)-1]
	}

	return fn.Name()
}

// WithSkip sets the number of stack frames to skip when identifying the caller.
func WithSkip(skip int) OptionFuncName {
	return func(o *OptFuncName) {
		o.Skip = skip + defaultSkip
	}
}
