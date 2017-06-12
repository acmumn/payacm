package main

import (
	"fmt"
	html "html/template"
	"net/smtp"
	text "text/template"

	"github.com/domodwyer/mailyak"
)

const MAIL_TEMPLATE_NAME = "mail_template"
const MAIL_TEMPLATE_HTML = `
<p>This receipt confirms payment of {{.Amount | asMoney}} from <a href="mailto:{{.Email}}">{{.Email}}</a> for "{{.Reason}}".</p>

<p>If this information is incorrect, please email <a href="mailto:acm@umn.edu">acm@umn.edu</a> immediately.</p>
`
const MAIL_TEMPLATE_TEXT = `
This receipt confirms payment of {{.Amount | asMoney}} from {{.Email}} for "{{.Reason}}".

If this information is incorrect, please email acm@umn.edu immediately.
`

var mail_template_html = html.Must(html.New(MAIL_TEMPLATE_NAME).Funcs(html.FuncMap{
	"asMoney": func(cents uint64) string {
		return fmt.Sprintf("$%.2f", float64(cents)/100)
	},
}).Parse(MAIL_TEMPLATE_HTML))
var mail_template_text = text.Must(text.New(MAIL_TEMPLATE_NAME).Funcs(text.FuncMap{
	"asMoney": func(cents uint64) string {
		return fmt.Sprintf("$%.2f", float64(cents)/100)
	},
}).Parse(MAIL_TEMPLATE_TEXT))

func mail(payment Payment) error {
	// Open a connection to the server.
	server := fmt.Sprintf("%s:%s", getenv("SMTP_HOST"), getenv("SMTP_PORT"))
	mail := mailyak.New(server, smtp.PlainAuth("",
		getenv("SMTP_USER"),
		getenv("SMTP_PASS"),
		getenv("SMTP_HOST")))

	// Set up the mail headers.
	mail.From(getenv("SMTP_FROM"))
	mail.FromName("payacmumn")
	mail.To(payment.Email)
	mail.Bcc("acm@umn.edu")
	mail.ReplyTo("acm@umn.edu")
	mail.Subject("Receipt from payacmumn")

	// Render the mail body and return.
	if err := mail_template_html.Execute(mail.HTML(), payment); err != nil {
		return err
	}
	if err := mail_template_text.Execute(mail.Plain(), payment); err != nil {
		return err
	}
	return mail.Send()
}
