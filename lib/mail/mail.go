package mail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"huma-app/lib/config"
	"path/filepath"
	"strings"

	gomail "gopkg.in/mail.v2"
)

type mail struct {
	tpl    emailTemplate
	to     string
	params interface{}
}

type emailTemplate string

const (
	verifyTemplate   emailTemplate = "verify.tpl"
	passwordTemplate emailTemplate = "password.tpl"
)

type templateData struct {
	subj string
	text string
	html string
}

type IncorrectTemplateFormat struct {
	Message string
}

func (e IncorrectTemplateFormat) Error() string {
	return e.Message
}

func (m *mail) compile() (*templateData, error) {
	var result bytes.Buffer
	t, err := template.New(string(m.tpl)).ParseFiles(filepath.Join("templates", "mail", string(m.tpl)))
	if err != nil {
		return nil, err
	}
	err = t.Execute(&result, m.params)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(result.String(), "-- ")
	if len(parts) < 4 {
		return nil, IncorrectTemplateFormat{Message: "Incorrect template format"}
	}
	return &templateData{
		strings.TrimSpace(strings.TrimPrefix(parts[1], "SUBJ")),
		strings.TrimSpace(strings.TrimPrefix(parts[2], "TEXT")),
		strings.TrimSpace(strings.TrimPrefix(parts[3], "HTML")),
	}, nil
}

func (m *mail) print() error {
	data, err := m.compile()
	if err != nil {
		return err
	}
	fmt.Printf("TO: %s\nSUBJ: %s\n\nTEXT:\n%s\n\nHTML:\n%s\n", m.to, data.subj, data.text, data.html)
	return nil
}

func (m *mail) send() error {

	data, err := m.compile()
	if err != nil {
		return err
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", config.Get().Smtp.From)
	msg.SetHeader("To", m.to)
	msg.SetHeader("Subject", data.subj)
	msg.SetBody("text/plain", data.text)
	msg.SetBody("text/html", data.html)

	d := gomail.NewDialer(
		config.Get().Smtp.Host,
		config.Get().Smtp.Port,
		config.Get().Smtp.From,
		config.Get().Smtp.Password,
	)
	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(msg); err != nil {
		return err
	}
	return nil
}

func sendMail(tpl emailTemplate, to string, params interface{}) error {
	m := &mail{tpl, to, params}
	if config.Get().Smtp.Password != "" {
		err := m.send()
		if err != nil {
			return err
		}
	} else {
		err := m.print()
		if err != nil {
			return err
		}
	}
	return nil
}

type VerifyEmailParams struct {
	AppName string
	Link    string
}

type PasswordEmailParams struct {
	AppName string
	Link    string
}

func SendVerifyMail(to string, params VerifyEmailParams) error {
	return sendMail(verifyTemplate, to, params)
}

func SendPasswordMail(to string, params PasswordEmailParams) error {
	return sendMail(passwordTemplate, to, params)
}
