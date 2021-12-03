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

	msg := []byte("From: " + emailConfig.EmailAddress + "\r\n" +
		"To: " + recipient + "\r\n" +
		"Subject: STOCKX - Confirm your account\r\n\r\n" +
		"No słuchaj byku, gratuluję, zarejestrowałeś się. STONKS!\r\n")

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	if err != nil {
		return err
	}

	fmt.Println("Email Sent Successfully!")

	return nil
}
