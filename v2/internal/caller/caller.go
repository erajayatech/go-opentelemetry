package caller

import (
	"runtime"
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

	return fn.Name()
}

// WithSkip sets the number of stack frames to skip when identifying the caller.
func WithSkip(skip int) OptionFuncName {
	return func(o *OptFuncName) {
		o.Skip = skip + defaultSkip
	}
}
