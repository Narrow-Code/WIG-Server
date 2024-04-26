package utils

import (
	"log"
	"runtime"
	"strings"
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
	functionName := callerFunction.Name()
    	functionName = functionName[strings.LastIndex(functionName, ".")+1:]
    	return functionName
}

func Log(message string) {	
	log.Printf("%s: %s", CallerFunctionName(2), message)
}
