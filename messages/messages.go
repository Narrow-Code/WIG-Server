package messages

const (
	// Common errors
	ErrorParsingRequest = "There was an error parsing JSON request."
	ErrorWithConnection = "Connection error"
	RecordNotFound = "Record was not found"
	UIDEmpty = "No UID found"
	TokenEmpty = "No token found"
	
	// Login messages
	UsernameDoesNotExist = "Username does not exist"
	UsernamePasswordDoNotMatch = "Username and password do not match"
	UserLoginSuccess = "User log in success"
	SaltReturned = "Salt was returned"
	ErrorToken = "Access denied by token"
	TokenPass = "Token authenticated"

	// Signup messages
	UsernameEmpty = "Username is required"
	UsernameInUse = "Username is already in use"
	EmailInUse = "Email associated with another account"
	EmailEmpty = "Email is required"
	SaltMissing = "Salt is missing"
	HashMissing = "Hash is missing"
	SignupSuccess = "User added successfully"
	ErrorUsernameRequirements = "Username does not match criteria"
	ErrorEmailRequirements = "Not a valid email address"
)
