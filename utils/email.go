package utils
                                
import (
    "fmt"
    "net/smtp"
    "os"
)

func SendVerificationEmail(receiver string, username string) {
    // TODO generate verification token & link
    verificationLink := "http://narrowcode.org"

    sender := os.Getenv("EMAIL")
    pass := os.Getenv("EMAIL_PASS")
    host := os.Getenv("EMAIL_HOST")
    port := os.Getenv("EMAIL_PORT")

    auth := smtp.PlainAuth("", sender, pass, host)

    // The message to send.
    to := []string{receiver}
    msg := []byte("To: " + receiver + "\r\n" +
        "Subject: WIG Verification\r\n" +
        "\r\n" +
        "Hello " + username + ",\r\n\r\n" +
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

