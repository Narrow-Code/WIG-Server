package Models

type User struct {
	Username 	string 	'json:"username"'
	UserEmail 	string 	'json:"user_email"'
	UserSalt 	string	'json:"user_salt"'
	UserHash 	string 	'json:"user_hash"'
	UserConfirmed 	string 	'json:"user_confirmed"'
	UserTier 	string 	'json:"user_tier"'
}
