/*
* The messages package holds constant error and success messages for easy access.
*/
package messages

const (
	ErrorParsingRequest = "There was an error parsing JSON request." // Error message for when the JSON request could not be parsed.
	ErrorWithConnection = "Connection error" // Error message for whenever there is a connectivity issue.
	ConversionError = "There was an error converting from string to int" // Error message for when there was an error message converting string to int
	NegativeError = "Cannot pass a negative argument" // Error message for when a negative argument is passed as a parameter.
	ErrorDeletingOwnership = "There was an error deleting ownership" // Error message for when there was an error deleting an ownership from the database.
	DoesNotExist = " does not exist" // Error message, meant to be appended after a field type to show that it the field does not exist.
	RecordInUse = " is in use" // Error message, meant to be appended after a field type to show that a field is already in use.
)
