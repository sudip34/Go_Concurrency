package main

import (
	"bytes"
	"fmt"
	"sync"
	"text/template"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
	Wait        *sync.WaitGroup
	MailerChan  chan Message
	ErrorChan   chan error
	DoneCahn    chan bool
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
	Template    string
}

// a function to listen for any mail
func (app *Config) listenForMail() {
	for {
		select {
		case msg := <-app.Mailer.MailerChan:
			go app.Mailer.sendMail(msg, app.Mailer.ErrorChan)
		case err := <-app.Mailer.ErrorChan:
			app.ErrorLog.Println(err)
		case <-app.Mailer.DoneCahn:
			return
		}
	}
}

// a function to listen for messages on the MailerChan

func (m *Mail) sendMail(msg Message, errorChan chan error) {
	defer m.Wait.Done()

	if msg.Template == "" {
		msg.Template = "mail"
	}

	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	//build html mail
	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		errorChan <- err
	}

	//build plain text mail
	palinMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		errorChan <- err
	}

	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		errorChan <- err
	}

	email := mail.NewMSG()
	email.SetFrom(msg.From).AddTo(msg.To).SetSubject(msg.Subject)

	email.SetBody(mail.TextPlain, palinMessage)
	email.AddAlternative(mail.TextHTML, formattedMessage)

	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}

	err = email.Send(smtpClient)
	if err != nil {
		errorChan <- err
	}
}

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	templateFile := fmt.Sprintf("./cmd/web/templates/%s.html.gohtml", msg.Template)

	htmlTemplate, err := template.New("email-htmp").ParseFiles(templateFile)
	if err != nil {
		return "", nil
	}

	// write the template out
	var bufferWriter bytes.Buffer
	if err := htmlTemplate.ExecuteTemplate(&bufferWriter, "body", msg.DataMap); err != nil {
		return "", err
	}

	// make the output of the bufferWriter as formatted String
	formattedMsg := bufferWriter.String()
	formattedMsg, err = m.inlineCSS(formattedMsg)
	if err != nil {
		return "", nil
	}
	return formattedMsg, nil

}

func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	templateFile := fmt.Sprintf("./cmd/web/templates/%s.html.gohtml", msg.Template)

	plainTemplate, err := template.New("email-htmp").ParseFiles(templateFile)
	if err != nil {
		return "", nil
	}

	// write the template out
	var bufferWriter bytes.Buffer
	if err := plainTemplate.ExecuteTemplate(&bufferWriter, "body", msg.DataMap); err != nil {
		return "", err
	}

	// make the output of the bufferWriter as formatted String
	plainMsg := bufferWriter.String()

	return plainMsg, nil
}

func (m *Mail) getEncryption(e string) mail.Encryption {
	switch e {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS

	}
}

func (m *Mail) inlineCSS(s string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	htmlPremailer, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", nil
	}

	htmlString, err := htmlPremailer.Transform()
	if err != nil {
		return "", nil
	}
	return htmlString, nil
}
