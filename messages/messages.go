/*
* The messages package holds constant error and success messages for easy access.
*/
package messages

const (
	ErrorParsingRequest = "There was an error parsing JSON request." // Error message for when the JSON request could not be parsed.
	ErrorWithConnection = "Connection error" // Error message for whenever there is a connectivity issue.
	RecordNotFound = "Record was not found" // Error message for when a record is not found within a database.
	UIDEmpty = "No UID found" // Error message for when the UID in a request is empty but required.
	QRMissing = "No QR found" // Error message for when the QR in a request is empty but required.
	ItemNotFound = "Item not found in database" // Error message for when the item barcode was not found in the database.

	Location = "LOCATION" // Success message for when the QR code is stored as a location.
	Ownership = "OWNERSHIP" // Success message for when the QR code is stored as an ownership.
	New = "NEW" // Success message for when the QR code has not been stored in the database.

	ConversionError = "There was an error converting from string to int" // Error message for when there was an error message converting string to int
	NegativeError = "Cannot pass a negative argument" // Error message for when a negative argument is passed as a parameter.
	BarcodeIntError = "Barcode must be of int value" // Error message for when a non-int value is passed to the barcode argument.
	ErrorDeletingOwnership = "There was an error deleting ownership" // Error message for when there was an error deleting an ownership from the database.
	DoesNotExist = " does not exist" // Error message, meant to be appended after a field type to show that it the field does not exist.
	RecordInUse = " is in use" // Error message, meant to be appended after a field type to show that a field is already in use.
)
