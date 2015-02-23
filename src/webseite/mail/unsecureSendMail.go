package mail

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"
)

type SMTPServer struct {
	Addr string
	Auth smtp.Auth
}

type Mail struct {
	Server  SMTPServer
	From    mail.Address
	To      mail.Address
	Subject string
	Message string

	header map[string]string
}

func (mail *Mail) prepareHeader() {
	mail.header = make(map[string]string)
	mail.header["From"] = mail.From.String()
	mail.header["To"] = mail.To.String()
	mail.header["Subject"] = encodeRFC2047(mail.Subject)
	mail.header["MIME-Version"] = "1.0"
	mail.header["Content-Type"] = "text/plain; charset=\"utf-8\""
	mail.header["Content-Transfer-Encoding"] = "base64"
}

func (mail *Mail) prepareMessageBody() []byte {
	mail.prepareHeader()

	message := ""
	for k, v := range mail.header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(mail.Message))
	return []byte(message)
}

func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{String, ""}
	return strings.Trim(addr.String(), " <>")
}

func (mail *Mail) Send() error {
	// TCP Connection
	c, err := smtp.Dial(mail.Server.Addr)
	if err != nil {
		return err
	}
	defer c.Close()

	// Say HELO
	if err = c.Hello("localhost"); err != nil {
		return err
	}

	// Check if we need TLS
	if ok, _ := c.Extension("STARTTLS"); ok {
		config := &tls.Config{
			ServerName:         mail.Server.Addr,
			InsecureSkipVerify: true,
		}
		if err = c.StartTLS(config); err != nil {
			return err
		}
	}

	// Auth if needed
	if err = c.Auth(mail.Server.Auth); err != nil {
		return err
	}

	// Send from mail
	if err = c.Mail(mail.From.Address); err != nil {
		return err
	}

	// Send receipt
	if err = c.Rcpt(mail.To.Address); err != nil {
		return err
	}

	// Get and write data
	w, err := c.Data()
	if err != nil {
		return err
	}

	// Write data to the pipe
	_, err = w.Write(mail.prepareMessageBody())
	if err != nil {
		return err
	}

	// Close the writer
	err = w.Close()
	if err != nil {
		return err
	}

	// Say bye to the SMTP Server
	return c.Quit()
}
