package messages

import (
	"io/ioutil"
	"strings"
)

type EmailMessage struct {
	Title     string
	Body      string
	FromName string
	FromEmail string
	ToEmail   string
	ReplyTo   string
	sender    IEmailSender
}

// NewEmailMessage
func NewEmailMessage(sender IEmailSender) *EmailMessage {
	return &EmailMessage{
		sender: sender,
	}
}

// SendActivateLink
func (e *EmailMessage) SendActivateLink(activateURL, appName string) error {

	e.Title = "Authorization on the service " + appName
	e.ReplyTo = ""
	content, err := ioutil.ReadFile("services/transport/messages/templates/registration.html")
	if err != nil {
		return err
	}

	// depend template for this message type
	tpl := string(content)
	tpl = strings.ReplaceAll(tpl, "{APPLICATION_NAME}", appName)
	tpl = strings.ReplaceAll(tpl, "{link}", activateURL)
	e.Body = tpl

	return e.sender.Send(e)
}
