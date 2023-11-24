/*
* The messages package holds constant error and success messages for easy access.
 */
package messages

const (
	TokenEmpty   = "no token found"         // Error message for when the token in a request is empty but required.
	ErrorToken   = "access denied by token" // Error message for when the users stored token does not match their UID in a request.
	TokenPass    = "Token authenticated"    // Success message for when the users stored token matches their UID in a request.
	AccessDenied = "Unauthorized"           // Error message for when the access was denied to make the API request.

	UsernameDoesNotExist       = "Username does not exist"            // Error message for when an invalid username is being passed in a request.
	UsernameEmpty              = "Username is required"               // Error message for when the username is missing in a request.
	UsernameInUse              = "Username is already in use"         // Error message for when the username in the request is already in use.
	ErrorUsernameRequirements  = "Username does not match criteria"   // Error message for when the username in the request does not match the criteria.
	UsernamePasswordDoNotMatch = "Username and password do not match" // Error message for when the username does not match the hash in the request.

	EmailInUse             = "Email associated with another account" // Error message for when the email in the request is already in use.
	EmailEmpty             = "Email is required"                     // Error message for when the email is missing in a request.
	ErrorEmailRequirements = "Not a valid email address"             // Error message for when the email in the request does not match the criteria.

	SaltMissing = "Salt is missing" // Error message for when the salt is missing in a request.
	HashMissing = "Hash is missing" // Error message for when the hash is missing in a request.

	UserLoginSuccess = "User log in success"     // Success message for when the username and hash match in the requests.
	SignupSuccess    = "User added successfully" // Success message for when the user is successesfully added to the database.
	SaltReturned     = "Salt was returned"       // Success message for when the users matching salt was successfully found in the database and returned.

)
