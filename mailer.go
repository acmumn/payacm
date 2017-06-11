package main

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
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

	// (Possibly) start encryption.
	if ok, params := conn.Extension("STARTTLS"); ok {
		log.Println("STARTTLS supported with params", params)
		conn.StartTLS(&tls.Config{
			ServerName: getenv("SMTP_HOST"),
		})
	}

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
	conn.Rcpt(payment.Email)
	conn.Rcpt("acm@umn.edu")

	// Create the request body.
	w, err := conn.Data()
	if err != nil {
		return err
	}
	defer w.Close()

	// Write the mail headers.
	ioutil.WriteString(w, fmt.Sprintf("To: %s\n", payment.Email))
	ioutil.WriteString(w, "Receipt from payacmumn")

	// Render the request body and return.
	return mail_template.Execute(w, payment)
}
