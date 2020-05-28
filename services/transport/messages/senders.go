package messages

import (
	"crypto/tls"
	"fmt"
	"github.com/lexbond13/api_core/config"
	"github.com/jordan-wright/email"
	"log"
	"net"
	"net/mail"
	"net/smtp"
)

var emailSender IEmailSender

type IEmailSender interface {
	FromEmail() string
	FromName() string
	Send(message *EmailMessage) error
}

// NewEmailSender
func InitEmailSender(config *config.Email) {
	//emailSender = NewMailGunSender(config.MailGunSender)
	emailSender = NewSMTPSender(config.SMTPSender)
}

type MailGunSender struct {
	Host     string
	User     string
	Password string
}

// NewMailGunSender
func NewMailGunSender(config *config.MailGunSender) IEmailSender {
	mailgun := &MailGunSender{
		Host:     config.Host,
		User:     config.Username,
		Password: config.Password,
	}

	return mailgun
}

// Send
func (m *MailGunSender) Send(message *EmailMessage) error {
	e := email.NewEmail()
	e.To = []string{message.ToEmail}
	e.Subject = message.Title
	e.HTML = []byte(message.Body)

	if message.FromEmail == "" {
		message.FromEmail = m.FromEmail()
	}

	e.From = fmt.Sprintf("\"%s\" <%s>", message.FromName, message.FromEmail)
	e.Sender = message.FromEmail

	// set default sender data if it's not set
	if e.From == "" {
		e.From = m.FromName()
	}

	if e.Sender == "" {
		e.Sender = m.FromEmail()
	}

	err := e.Send(m.Host +":587", smtp.PlainAuth("", m.User, m.Password, m.Host))
	if err != nil {
		return err
	}

	return nil
}

func (m *MailGunSender) FromEmail() string {
	return m.User
}


func (m *MailGunSender) FromName() string {
	return m.User
}

type SMTPSender struct {
	Auth        smtp.Auth
	AddrString  string
	SenderEmail string
	SenderName  string
}

// NewSMTPSender
func NewSMTPSender(config *config.SMTPSender) IEmailSender {
	auth := smtp.PlainAuth(
		"",
		config.AuthEmail,
		config.Password,
		config.Host,
	)

	addrString := fmt.Sprintf("%s:%d", config.Host, config.Port)

	sender := &SMTPSender{
		Auth: auth,
		AddrString: addrString,
		SenderEmail: config.SenderEmail,
		SenderName: config.SenderName,
	}

	return sender
}

// Send
func (em *SMTPSender) Send(message *EmailMessage) error {

	from := mail.Address{em.SenderName, em.SenderEmail}
	to   := mail.Address{"", message.ToEmail}

	if message.FromName != "" {
		from.Name = message.FromName
	}

	if message.FromEmail != "" {
		from.Address = message.FromEmail
	}

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = message.Title
	headers["MIME-version"] = "1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	// Setup message
	newMessage := ""
	for k,v := range headers {
		newMessage += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	newMessage += "\r\n" + message.Body

	// Connect to the SMTP Server
	host, _, _ := net.SplitHostPort(em.AddrString)

	// TLS config
	tlsconfig := &tls.Config {
		InsecureSkipVerify: true,
		ServerName: host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", em.AddrString, tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	// Auth
	if err = c.Auth(em.Auth); err != nil {
		return err
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		return err
	}

	if err = c.Rcpt(to.Address); err != nil {
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(newMessage))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	err = c.Quit()
	if err != nil {
		return err
	}

	return nil
}

func (em *SMTPSender) FromEmail() string {
	return em.SenderEmail
}

func (em *SMTPSender) FromName() string {
	return em.SenderName
}

func GetEmailSender() IEmailSender {
	return emailSender
}
