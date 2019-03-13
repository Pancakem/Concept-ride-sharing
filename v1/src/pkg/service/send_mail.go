package service

import (
	"bytes"
	"html/template"
	"log"
	"net/smtp"
)

var authe smtp.Auth

// Request struct
type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

// NewRequest returns a new request
func NewRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}

// SendMail sends a mail
// success returns a nil or failure returns a err
func (r *Request) SendMail() error {
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(subject + mime + "\n" + r.body)
	addr := "smtp.gmail.com"
	

	if err := smtp.SendMail(addr, authe, "westcoastcustomskenya@gmail.com", r.to, msg); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// ParseTemplate gets the format of the email being sent
func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}
