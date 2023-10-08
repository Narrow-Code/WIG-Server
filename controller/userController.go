package controller

import (
	"github.com/gofiber/fiber/v2"
	models "WIG-Server/models"
	db "WIG-Server/config"
)

func Signup(c *fiber.Ctx) error {
	
	var data map[string]string
	
	err := c.BodyParser(&data)

	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"success":false,
				"message":"Invalid data",
			})
	}

	if data["username"] == "" {
		return c.Status(400).JSON(
			fiber.Map{
				"success":false,
				"message":"Username is required"})
	}

	if data["email"] == "" {
		return c.Status(400).JSON(
			fiber.Map{
				"success":false,
				"message":"Email is required"})
		}

	if data["salt"] == "" {
		return c.Status(400).JSON(
			fiber.Map{
				"success":false,
				"message":"Salt is missing"})
		}

	if data["hash"] == "" {
		return c.Status(400).JSON(
			fiber.Map{
				"success":false,
				"message":"Hash is missing"})
		}

	var user models.User
	userResult := db.DB.Where("username = ?", data["username"]).First(&user)
	if userResult.RowsAffected != 0 {
		return c.Status(400).JSON(
			fiber.Map{
				"success":false,
				"message":"Username already in use"})
		}

	emailResult := db.DB.Where("user_email = ?", data["email"]).First(&user)
	if emailResult.RowsAffected != 0 {
		return c.Status(400).JSON(
			fiber.Map{
				"success":false,
				"message":"Email associated with another account"})
		}

	// Check username requirements

	// Check email validity

	user = models.User{
		Username: data["username"],
		UserEmail: data["email"],
		UserSalt: data["salt"],
		UserHash: data["hash"],
	}

	db.DB.Create(&user)
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "User added successfully",
		"data": user })


	// Send verification email

}
