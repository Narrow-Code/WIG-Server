/*
* The messages package holds constant error and success messages for easy access.
 */
package messages

const (
	RecordNotFound  = "record was not found"         // Error message for when a record is not found within a database.
	UIDEmpty        = "no UID found"                 // Error message for when the UID in a request is empty but required.
	QRMissing       = "no QR found"                  // Error message for when the QR in a request is empty but required.
	ItemNotFound    = "item not found in database"   // Error message for when the item barcode was not found in the database.
	Location        = "LOCATION"                     // Success message for when the QR code is stored as a location.
	Ownership       = "OWNERSHIP"                    // Success message for when the QR code is stored as an ownership.
	New             = "NEW"                          // Success message for when the QR code has not been stored in the database.
	BarcodeIntError = "barcode must be of int value" // Error message for when a non-int value is passed to the barcode argument.
	BarcodeMissing  = "barcode is required"          // Error message for when a barcode is passed null.
)
