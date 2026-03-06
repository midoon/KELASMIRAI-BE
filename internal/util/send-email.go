package util

import (
	"bytes"
	"html/template"

	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

func SendEmail(to string, subject string, templatePath string, data interface{}, emailConfig *viper.Viper) error {
	host := emailConfig.GetString("email.smtp.host")
	port := emailConfig.GetInt("email.smtp.port")
	username := emailConfig.GetString("email.smtp.username")
	password := emailConfig.GetString("email.smtp.password")

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(host, port, username, password)

	return d.DialAndSend(m)
}
