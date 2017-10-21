package Util

import (
	"crypto/tls"
	"log"
	"strconv"

	"gopkg.in/gomail.v2"

	settings "github.com/horae-app/api/Settings"
)

func Invite(email string, name string, token int) (bool, string) {
	m := gomail.NewMessage()
	m.SetHeader("From", settings.MAIL_FROM)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Welcome to Horae app - Your token")
	m.SetBody("text/html", "Hello <b>"+name+"</b><br><br>Your access token is <b>"+strconv.Itoa(token)+"</b>")

	d := gomail.NewDialer(settings.MAIL_SMTP, settings.MAIL_PORT, settings.MAIL_USER, settings.MAIL_PWD)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		log.Println("[Error] Could not send invite to " + email)
		return false, err.Error()
	}

	return true, ""
}
