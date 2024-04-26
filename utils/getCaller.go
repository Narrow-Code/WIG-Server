package utils

import (
	"WIG-Server/models"
	"log"
	"runtime"
	"strings"

	"github.com/gofiber/fiber/v2"
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

/*
* Log prints a log message with the CallerFunctionName appended to it
*
* @param message The message to print
*/
func Log(message string) {	
	log.Printf("%s: %s", CallerFunctionName(2), message)
}

/*
* UserLog prints a log message with the CallerFunctionName and Username if applicable
*
* @param c The fiber context
* @param message The message to print
*/
func UserLog(c *fiber.Ctx, message string) {
	// Initialize user variable
	user := c.Locals("user").(models.User)

	// Append username if applicable
	if user.Username != "" {
		log.Printf("%s#%s: %s", user.Username, CallerFunctionName(2), message)
	} else {
		log.Printf("%s: %s", CallerFunctionName(2), message)
	}		
}
