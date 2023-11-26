package utils

import (
	"runtime"
)

/*
* Uses the runtime traceback to return the name of a function that made the current function call.
*
* @param callback The callback number of how far the traceback should get a function name.
*
* @return string The original functions name.
*/
func CallerFunctionName(callback int) string {
	pc, _, _, _ := runtime.Caller(0 + callback)
	callerFunction := runtime.FuncForPC(pc)
	if callerFunction != nil {
		return callerFunction.Name()
	}
	return "unknown"
}

