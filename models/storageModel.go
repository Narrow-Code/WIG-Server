package Models

type Storage struct {
	StorageUID 	string 	'json:"storage_uid"'
	StorageOwner 	string 	'json:"storage_owner"'
	StorageName 	string	'json:"storage_name"'
	StorageType 	string 	'json:"storage_type"'
	StorageLocation string 	'json:"storage_location"'
	StorageQR 	string 	'json:"storage_qr"'
	StorageTags	string 	'json:"storage_tags"'
}
