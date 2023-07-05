package email

import (
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"github.com/rmarasigan/bus-ticketing/internal/utility"
)

// Reference: https://www.w3.org/Protocols/rfc1341/7_3_Message.html
const MIME = "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"utf-8\";\r\nContent-Transfer-Encoding: quoted-printable\r\n"

// Content represents the content of an e-mail message, subject,
// and to whom it will be sent.
type Content struct {
	To      []string // Primary recipients of the e-mail
	CC      []string // Carbon Copy; E-mail addresses of the recipients who will receive a copy of the e-mail
	BCC     []string // Blind Carbon Copy; E-mail addresses of the recipients who will receive a copy of the e-mail
	Subject string   // Subject of the e-mail
	Message string   // The main body content of the e-mail
}

// Configuration contains the configuration settings for sending
// emails to the client.
type Configuration struct {
	ServerAddress   string  `json:"server_address"`    // It is the Outgoing Mail (SMTP) Server
	Username        string  `json:"username"`          // The e-mail address
	Password        string  `json:"password"`          // The e-mail application password
	Port            string  `json:"port"`              // The Outgoing Mail (SMTP) Port
	Content         Content `json:"content,omitempty"` // Content of the email, subject and recipients
	CustomerSupport string  `json:"customer_support"`  // The e-mail address of the customer support service
}

// authentication implements the smtp.PlainAuth that will return
// an Auth that has the PLAIN authentication.
func (email Configuration) authentication() smtp.Auth {
	return smtp.PlainAuth("", email.Username, email.Password, email.ServerAddress)
}

// message sets the mandatory headers and content of the message.
//  Mandatory headers
//   - Date: Time and date the message was sent
//   - From: Provides the sender's name
//   - To: To whom the message will be delivered
// Reference: https://www.oreilly.com/library/view/programming-internet-email/9780596802585/ch02s04.html
func (email Configuration) message() []byte {
	var msg string

	msg = fmt.Sprintf("From: %s\r\n", email.Username)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(email.Content.To, ","))
	msg += fmt.Sprintf("Date: %s\r\n", time.Now().UTC().Format(time.RFC1123Z))
	msg += fmt.Sprintf("Subject: %s\r\n", email.Content.Subject)
	msg += MIME + "\r\n"
	msg += email.Content.Message
	msg += "\r\n\r\n"

	return []byte(msg)
}

// Send authenticates and connects to the server (e.g. smtp.gmail.com:587) and
// sends an e-mail from address "from", to addresses "to", with message.
func (email Configuration) Send() error {
	host := fmt.Sprintf("%s:%s", email.ServerAddress, email.Port)
	return smtp.SendMail(host, email.authentication(), "", email.Content.To, email.message())
}

// Error sets the default key-value pair.
func (email Configuration) Error(err error, code, message string, kv ...utility.KVP) {
	kv = append(kv, utility.KVP{Key: "Integration", Value: "Bus Ticketing – Email"})
	utility.Error(err, code, message, kv...)
}
