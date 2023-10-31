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
	TokenEmpty = "No token found" // TokenEmpty is a constant string error message for when the token in a request is empty but required.
	AccessDenied = "Unauthorized" // AccessDenied is a constant string error message for when the access was denied to make the API request.
	ItemNotFound = "Item not found in database" // ItemNotFound is a constant string error message for when the item barcode was not found in the database.

	UsernameDoesNotExist = "Username does not exist" // UsernameDoesNotExist is a constant string error message for when an invalid username is being passed in a request.
	UsernamePasswordDoNotMatch = "Username and password do not match" // UsernamePasswordDoNotMatch is a constant string error message for when the username does not match the hash in the request.
	UserLoginSuccess = "User log in success" // UserLoginSuccess is a constant string success message for when the username and hash match in the requests.
	SaltReturned = "Salt was returned" // SaltReturned is a constant string success message for when the users matching salt was successfully found in the database and returned.
	ErrorToken = "Access denied by token" // ErrorToken is a constant string error message for when the users stored token does not match their UID in a request.
	TokenPass = "Token authenticated" // TokenPass is a constant string success message for when the users stored token matches their UID in a request.

	UsernameEmpty = "Username is required" // UsernameEmpty is a constant string error message for when the username is missing in a request.
	UsernameInUse = "Username is already in use" // UsernameInUse is a constant string error message for when the username in the request is already in use.
	EmailInUse = "Email associated with another account" // EmailInUse is a constant string error message for when the email in the request is already in use.
	EmailEmpty = "Email is required" // EmailEmpty is a constant string error message for when the email is missing in a request.
	SaltMissing = "Salt is missing" // SaltMissing is a constant string error message for when the salt is missing in a request.
	HashMissing = "Hash is missing" // HashMissing is a constant string error message for when the hash is missing in a request.
	SignupSuccess = "User added successfully" // SignupSuccess is a constant string success message for when the user is successesfully added to the database.
	ErrorUsernameRequirements = "Username does not match criteria" // ErrorUsernameRequirements is a constant string error message for when the username in the request does not match the criteria.
	ErrorEmailRequirements = "Not a valid email address" // ErrorEmailRequirements is a constant string error message for when the email in the request does not match the criteria.
	Location = "LOCATION" // Location is a constant string success message for when the QR code is stored as a location.
	Ownership = "OWNERSHIP" // Ownership is a constant string success message for when the QR code is stored as an ownership.
	New = "NEW" // New is a constant string success message for when the QR code has not been stored in the database.

	ConversionError = "There was an error converting from string to int" // UIDConversionError is a constant string error message for when there was an error message converting string to int
)
