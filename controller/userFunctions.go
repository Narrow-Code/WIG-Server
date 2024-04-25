package controller

import (
	"WIG-Server/db"
	"WIG-Server/models"

	"github.com/google/uuid"
)

/*
* createUser creates a new user and adds it to the database
*
* @param data The data map with all of the user information
* @return models.User the create User model
*/
func createUser(data map[string]string) models.User {
	// Build user
	user := models.User{
		Username: data["username"],
		Email:    data["email"],
		Salt:     data["salt"],
		Hash:     data["hash"],
		UserUID:  uuid.New(),
	}

	// Create user and return
	db.DB.Create(&user)
	return user
}
