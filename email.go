package main

import (
	"net/smtp"
)

func ConfirmationEmail(domain, link, username, email string) error {
	auth := smtp.PlainAuth("", "", "", "localhost")
	msgString := "From: " + domain + " <noreply@" + domain + ">\n"
	msgString += "To: " + email + "\n"
	msgString += "Subject: Welcome, " + username + "\n"
	msgString += "\n"
	msgString += "Hi and welcome to " + domain + "!\n"
	msgString += "\n"
	//msgString += "Your username is: " + username + "\n"
	//msgString += "\n"
	msgString += "Confirm the registration by following this link:\n"
	msgString += link + "\n"
	msgString += "\n"
	msgString += "Thank you.\n"
	msgString += "\n"
	msgString += "Best regards,\n"
	msgString +="    The " + domain + " registration system\n"
	msg := []byte(msgString)
	from := "noreply@" + domain
	to := []string{email}
	host := "localhost:25"
	return smtp.SendMail(host, auth, from, to, msg)
}

