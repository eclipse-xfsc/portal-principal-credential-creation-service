package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strconv"
	"text/template"

	mail "github.com/xhit/go-simple-mail/v2"
)

func sendEmail(emailTemplate string, templateKeys map[string]string, qrCode string, emailRecipient string) error {
	// Build message body
	template, err := template.New("invitation").Parse(emailTemplate)
	if err != nil {
		Logger.Error(err)
		err = fmt.Errorf("Can't parse template")
		return err
	}
	var tpl bytes.Buffer
	err = template.Execute(&tpl, templateKeys)
	if err != nil {
		Logger.Error(err)
		err = fmt.Errorf("Can't parse template")
		return err
	}
	config, _ := getConfig()
	server := mail.NewSMTPClient()
	server.Host = config.mailSmtpHost
	server.Port, _ = strconv.Atoi(config.mailSmtpPort)
	server.Username = config.mailSmtpUsername
	server.Password = config.mailSmtpPassword
	server.Encryption = mail.EncryptionTLS

	smtpClient, err := server.Connect()
	if err != nil {
		Logger.Error(err)
	}

	// Create email
	email := mail.NewMSG()
	email.SetFrom("GXFS Principal Invitation Service<" + config.mailSupportAddress + ">")
	email.AddTo(emailRecipient)
	email.SetSubject("GXFS Federation Invitiation")

	email.SetBody(mail.TextHTML, tpl.String())

	qr, _ := base64.StdEncoding.DecodeString(qrCode)

	email.Attach(&mail.File{Data: qr, Name: "invitation.png", Inline: true})

	err = email.Send(smtpClient)
	if err != nil {
		Logger.Error(err)
	}

	return nil
}
