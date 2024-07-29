package utils

import (
	"WIG-Server/models"
	"fmt"
	"net/smtp"
	"os"
)

func SendVerificationEmail(user models.User) {
    // TODO generate verification token & link
    verificationLink := "http://narrowcode.org"

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
    // See if token exists for User`
    resetLink := "http://narrowcode.org"

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
