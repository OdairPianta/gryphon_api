package notifications

import (
	"context"
	"fmt"
	"net/smtp"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

func ToMail(toAddress string, subject string, message string) error {
	var mailHost = os.Getenv("MAIL_HOST")
	var mailPort = os.Getenv("MAIL_PORT")
	var mailUsername = os.Getenv("MAIL_USERNAME")
	var mailPassword = os.Getenv("MAIL_PASSWORD")
	// var mailEncryption = os.Getenv("MAIL_ENCRYPTION")
	var mailFromAddress = os.Getenv("MAIL_FROM_ADDRESS")
	// var mailFromName = os.Getenv("MAIL_FROM_NAME")
	// var mailCharset = os.Getenv("MAIL_CHARSET")
	// var mailAutoTls = os.Getenv("MAIL_AUTO_TLS")
	var appName = os.Getenv("APP_NAME")

	var completeMessage = []byte("From: " + appName + " <" + mailFromAddress + ">\r\n" +
		"To: " + toAddress + "\r\n" +
		"Subject: " + subject + "!\r\n" +
		"\r\n" +
		message + "\r\n")

	var bytesMessage = []byte(completeMessage)

	// Create authentication
	auth := smtp.PlainAuth("", mailUsername, mailPassword, mailHost)

	// Send actual message
	err := smtp.SendMail(mailHost+":"+mailPort, auth, mailFromAddress, []string{toAddress}, bytesMessage)
	if err != nil {
		return err
	}
	return nil
}

func ToFcm(title string, body string, toToken string) error {

	opt := option.WithCredentialsFile("service_account_key.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}
	fcmClient, err := app.Messaging(context.Background())

	if err != nil {
		return err
	}

	_, err = fcmClient.Send(context.Background(), &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Token: toToken,
	})

	if err != nil {
		return err
	}

	return nil
}
