package utils

import (
	"runtime"
)

func CallerFunctionName() string {
	pc, _, _, _ := runtime.Caller(2)
	callerFunction := runtime.FuncForPC(pc)
	if callerFunction != nil {
		return callerFunction.Name()
	}
	return "unknown"
}

