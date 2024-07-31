package verification

import (
	"WIG-Server/db"
	"WIG-Server/models"
	"WIG-Server/utils"
	"encoding/base64"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"regexp"
	"time"
	"unicode"

	"github.com/gofiber/fiber/v2"
)

func SendVerificationEmail(user models.User) {
    // Generate token and verification link
    token, _ := GenerateVerificationToken(user)
    verificationLink := "http://ec2-18-209-15-108.compute-1.amazonaws.com:30001/verification/" + token

    sender := os.Getenv("EMAIL")
    pass := os.Getenv("EMAIL_PASS")
    host := os.Getenv("EMAIL_HOST")
    port := os.Getenv("EMAIL_PORT")

    auth := smtp.PlainAuth("", sender, pass, host)

    // The message to send.
    to := []string{user.Email}
    msg := []byte("To: " + user.Email + "\r\n" +
        "Subject: WIG Verification\r\n" +
        "\r\n" +
        "Hello " + user.Username + ",\r\n\r\n" +
        "We just need to verify your email address before you can access WIG.\r\n\r\n" +
        "Verify your email address here: " + verificationLink + "\r\n\r\n" +
        "Thanks, \r\n\r\n" +
        "Narrow Code")

    // Send the email.
    err := smtp.SendMail(host + ":" + port, auth, sender, to, msg)
    if err != nil {
        fmt.Println("Failed to send email:", err)
        return
    }

    fmt.Println("Email sent successfully!")
}

func SendResetPasswordEmail(user models.User) {
    // Generate token and verification link
    token, _ := GeneratePassowrdToken(user)
    resetLink := "http://ec2-18-209-15-108.compute-1.amazonaws.com:30001/resetpassword/" + token

    sender := os.Getenv("EMAIL")
    pass := os.Getenv("EMAIL_PASS")
    host := os.Getenv("EMAIL_HOST")
    port := os.Getenv("EMAIL_PORT")

    auth := smtp.PlainAuth("", sender, pass, host)

    // The message to send.
    to := []string{user.Email}
    msg := []byte("To: " + user.Email + "\r\n" +
        "Subject: WIG Verification\r\n" +
        "\r\n" +
        "Hello " + user.Username + ",\r\n\r\n" +
        "You have requested to reset your password for our WIG service.\r\n\r\n" +
        "If you would like to reset your password please click the following link: " + resetLink + "\r\n\r\n" +
        "Thanks, \r\n\r\n" +
        "Narrow Code")

    // Send the email.
    err := smtp.SendMail(host + ":" + port, auth, sender, to, msg)
    if err != nil {
        fmt.Println("Failed to send email:", err)
        return
    }

    fmt.Println("Email sent successfully!")
}

func ResetPasswordPage(c *fiber.Ctx) error {
	utils.Log("Started Reset Password")
    	uid := c.Params("uid")

	// Set regex statement
	regex := `^[A-Za-z\d\s!@#$%^&*()_+={}\[\]:;<>,.?~\\-]{8,}$`
	passwordRegex, err := regexp.Compile(regex)
	if err != nil {
        	utils.Log("Error compiling regex: " + err.Error())
        	return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error.")
    	}

	// Check if verification token is valid and get user
	var passwordChange models.PasswordChange
	var user models.User
	result := db.DB.Where("password_change_token = ?", uid).First(&passwordChange)
	if result.Error != nil || result.RowsAffected == 0 || passwordChange.PasswordExpiresAt.Before(time.Now()) {
		return c.SendString("The password change link you followed does not exist or has expired.")
	}
	utils.Log("token found")
	db.DB.Where("user_uid = ?", passwordChange.PasswordUserID).First(&user)

	// Start GET or POST
    	if c.Method() == fiber.MethodGet {
		utils.Log("Reset Password page began")
        	return c.Render("verification/reset.html", fiber.Map{
			"UID": uid, })
    	} else if c.Method() == fiber.MethodPost {
		utils.Log("Password submitted")
        	newPassword := c.FormValue("password")
		
		// Validate Password
		if passwordRegex.MatchString(newPassword) && containsDigit(newPassword) && containsUppercase(newPassword) {
			salt, err := GenerateSalt()
			if err != nil {
				utils.Log("error with salt")
				log.Fatalf("Error generating salt: %v", err)
			}
			utils.Log("Salt created")

			// Generate hash
			hash, err := GenerateHash(newPassword, salt)
			if err != nil {
				utils.Log("error with hash")
				log.Fatalf("Error generating hash: %v", err)
			}
			utils.Log("hash made")

			// Save Salt and Hash
			user.Hash = hash
			user.Salt = salt

			// Save User settings and delete token
			if err := db.DB.Save(&user).Error; err != nil {
            			utils.Log("error saving user: " + err.Error())
            			return c.Status(fiber.StatusInternalServerError).SendString("Error updating user.")
        		}

			db.DB.Delete(&passwordChange)
			utils.Log("user saved")

        		return c.SendString(fmt.Sprintf("Password changed successfully:"))
		}
		return c.Render("verification/reset.html", fiber.Map{
			"UID": uid,
			"Error": "Password does not meet criteria.",
		})
    	}
    	return nil
}

func containsUppercase(s string) bool {
    for _, ch := range s {
        if unicode.IsUpper(ch) {
            return true
        }
    }
    return false
}

func containsDigit(s string) bool {
    for _, ch := range s {
        if unicode.IsDigit(ch) {
            return true
        }
    }
    return false
}
