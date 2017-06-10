package main

import (
	"fmt"
	"html/template"
	"net/smtp"
)

var mail_template = template.Must(template.ParseFiles("templates/mail.html"))

func mail(payment Payment) error {
	// Open a connection to the server.
	conn, err := smtp.Dial(fmt.Sprintf("%s:%s",
		getenv("SMTP_HOST"),
		getenv("SMTP_PORT")))
	if err != nil {
		return err
	}
	defer conn.Close()

	// Authenticate to the server.
	err = conn.Auth(smtp.PlainAuth("",
		getenv("SMTP_USER"),
		getenv("SMTP_PASS"),
		getenv("SMTP_HOST")))
	if err != nil {
		return err
	}

	// Set up the mail headers.
	conn.Mail(getenv("SMTP_FROM"))
	conn.Rcpt(payment.email)
	conn.Rcpt("acm@umn.edu")

	// Create the request body.
	w, err := conn.Data()
	if err != nil {
		return err
	}
	defer w.Close()

	// Render the request body and return.
	return mail_template.Execute(w, payment)
}
