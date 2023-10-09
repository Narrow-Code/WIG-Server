package controller

import (
	"github.com/gofiber/fiber/v2"
	models "WIG-Server/models"
	db "WIG-Server/config"
)

func Signup(c *fiber.Ctx) error {
	
	var data map[string]string
	
	err := c.BodyParser(&data)

	// If theres an error in delivery?
	if err != nil {
		// TODO log here
		return c.Status(400).JSON(
			fiber.Map{
				"success":false,
				"message":"Invalid data",
			})
	}
	
	// If username is empty
	if data["username"] == "" {
		// TODO log here
		return c.Status(400).JSON(
			fiber.Map{
				"success":false,
				"message":"Username is required"})
	}

	// If email is empty
	if data["email"] == "" {
		// TODO log here
		return c.Status(400).JSON(
			fiber.Map{
				"success":false,
				"message":"Email is required"})
		}
	
	// If salt is empty
	if data["salt"] == "" {
		// TODO log here
		return c.Status(400).JSON(
			fiber.Map{
				"success":false,
				"message":"Salt is missing"})
		}

	// If hash is empty
	if data["hash"] == "" {
		return c.Status(400).JSON(
			// TODO log here
			fiber.Map{
				"success":false,
				"message":"Hash is missing"})
		}

	// Checks if username exists in database
	var user models.User
	userResult := db.DB.Where("username = ?", data["username"]).First(&user)
	if userResult.RowsAffected != 0 {
		// TODO log here
		return c.Status(400).JSON(
			fiber.Map{
				"success":false,
				"message":"Username already in use"})
		}

	// Chefks if email exists in database
	emailResult := db.DB.Where("user_email = ?", data["email"]).First(&user)
	if emailResult.RowsAffected != 0 {
		// TODO log here
		return c.Status(400).JSON(
			fiber.Map{
				"success":false,
				"message":"Email associated with another account"})
		}

	// TODO Check username requirements

	// TODO Check email validity

	user = models.User{
		Username: data["username"],
		UserEmail: data["email"],
		UserSalt: data["salt"],
		UserHash: data["hash"],
	}

	// TODO get error for failure?
	db.DB.Create(&user)

	// TODO log here a success log
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "User added successfully",
		"data": user })


	// TODO Send verification email

}
