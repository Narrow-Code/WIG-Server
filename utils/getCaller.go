package utils

import (
	"runtime"
)

/*
 * CallerFunctionName uses the runtime traceback to return the name of a function that made the current function call.
 *
 * @param callback The callback number of how far the traceback should get a function name.
 * @return string The original function's name.
 */
func CallerFunctionName(callback int) string {
	pc, _, _,ok := runtime.Caller(0 + callback)
	if !ok {
		return "unkown"
	}
	
	callerFunction := runtime.FuncForPC(pc)
	if callerFunction == nil {
		return "unknown"
	}

	return callerFunction.Name()
}

