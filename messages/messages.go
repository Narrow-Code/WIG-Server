/*
* The messages package holds constant error and success messages for easy access.
*/
package messages

const (
	ErrorParsingRequest = "There was an error parsing JSON request." // ErrorParsingRequest is a constant string error message for when the JSON request could not be parsed.
	ErrorWithConnection = "Connection error" // ErrorWithConnection is a constant string error message for whenever there is a connectivity issue.
	RecordNotFound = "Record was not found" // RecordNotFound is a constant string error message for when a record is not found within a database.
	UIDEmpty = "No UID found" // UIDEmpty is a constant string error message for when the UID in a request is empty but required.
	QRMissing = "No QR found" // QRMissing is a constant string error message for when the QR in a request is empty but required.
	ItemNotFound = "Item not found in database" // ItemNotFound is a constant string error message for when the item barcode was not found in the database.


	Location = "LOCATION" // Location is a constant string success message for when the QR code is stored as a location.
	Ownership = "OWNERSHIP" // Ownership is a constant string success message for when the QR code is stored as an ownership.
	New = "NEW" // New is a constant string success message for when the QR code has not been stored in the database.

	ConversionError = "There was an error converting from string to int" // UIDConversionError is a constant string error message for when there was an error message converting string to int
	NegativeError = "Cannot pass a negative argument" // NegativeError is a constant string error message for when a negative argument is passed as a parameter.
	BarcodeIntError = "Barcode must be of int value" // BarcodeIntError is a constant string error message for when a non-int value is passed to the barcode argument.
	ErrorDeletingOwnership = "There was an error deleting ownership" // ErrorDeletingOwnership is a constant string errormessage for when there was an error deleting an ownership from the database.
	DoesNotExist = " does not exist" // DoesNotExist is a constant string error message, meant to be appended after a field type to show that it the field does not exist.
	RecordInUse = " is in use" // RecordInUse is a constant string error message, meant to be appended after a field type to show that a field is already in use.
)
