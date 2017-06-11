package main

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/smtp"
)

const MAIL_FILE = "mail_template.html"
const MAIL_TEMPLATE = `
<p>This receipt confirms payment of {{.Amount | asMoney}} from <a href="{{.Email}}">{{.Email}}</a> for {{.Reason}}.</p>

<p>If this information is incorrect, please email <a href="mailto:acm@umn.edu">acm@umn.edu</a> immediately.</p>
`

var mail_template = template.Must(template.New(MAIL_FILE).Funcs(template.FuncMap{
	"asMoney": func(cents uint64) string {
		return fmt.Sprintf("$%.2f", float64(cents)/100)
	},
}).Parse(MAIL_TEMPLATE))

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
	io.WriteString(w, fmt.Sprintf("To: %s\n", payment.Email))
	io.WriteString(w, "Receipt from payacmumn")

	// Render the request body and return.
	return mail_template.Execute(w, payment)
}
