package email

import (
	"fmt"
	"net/smtp"
	"stockx-backend/conf"
)

func SendConfirmRegistrationEmail(recipient, username, token string) error {
	// Sender data.
	from := conf.Conf.Email.EmailAddress
	password := conf.Conf.Email.Password

	// Receiver email address.
	to := []string{
		recipient,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	msg := []byte("From: " + from + "\r\n" +
		"To: " + recipient + "\r\n" +
		"Subject: StockX - Confirm your account\r\n\r\n" +
		"Hey " + username + " - Welcome to StockX!\r\n\r\n" +
		"We hope you will have good time exploring Stock Market world, that you will learn a lot about it. We want to deliver a tool so you can do so without taking any risks.\r\n\r\n" +
		"Please confirm your account using following link: " +
		"http://stockx-lb-1734521826.eu-west-2.elb.amazonaws.com" + conf.Conf.Server.DevPort + "/user/confirm?t=" + token + "\r\n")

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	if err != nil {
		return err
	}

	fmt.Println("Confirmation email sent successfully!")

	return nil
}
