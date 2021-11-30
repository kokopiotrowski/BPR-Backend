package email

import (
	"fmt"
	"net/smtp"
	"stockx-backend/conf"
)

func SendConfirmRegistrationEmail(recipient string, emailConfig conf.EmailConfigurations) error {
	// Sender data.
	from := emailConfig.EmailAddress
	password := emailConfig.Password

	// Receiver email address.
	to := []string{
		recipient,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message.
	message := []byte("Hello bro, you registered your account, congrats to you")

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return err
	}

	fmt.Println("Email Sent Successfully!")

	return nil
}
