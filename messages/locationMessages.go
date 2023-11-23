package messages

const (
	LocationTypeInvalid  = "Lcoation type is invalid"       // Error for when an invalid location type was passed.
	LocationQRRequired   = "Location QR is required"        // Error for when a QR code is not passed for the lcoation.
	LocationNameRequired = "Location name is required"      // Error for when a location name is not passed for the lcoation.
	LocationAdded        = "Location added successfully"    // Success message for when locaton is added to database.
	LocationSelfError    = "Locatoin cannot be set as self" // Error message for when a lcoatoin is attempting to set itself as its location
	LocationUpdated      = "Location updated successfully"  // Success message for when location is upated.
	LocationNotFound     = "Location not found"             // Error message for when location is not found in database
)
