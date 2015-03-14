package mail

import (
	"bytes"
	"github.com/astaxie/beego"
	netMail "net/mail"
	"net/smtp"
	"os"
	"regexp"
	"text/template"
)

func SendTemplateMail(to string, templateName string, data map[string]string) *error {
	// Send mail
	info, _ := os.Stat("views/email/" + templateName + ".txt")
	file, err := os.OpenFile("views/email/"+templateName+".txt", 0, 0666)
	if err != nil {
		return &err
	}

	buffer := make([]byte, info.Size())
	_, errRead := file.Read(buffer)
	if errRead != nil {
		return &errRead
	}

	errClose := file.Close()
	if errClose != nil {
		return &errClose
	}

	contentString := string(buffer)
	templateEngine := template.New("views/email/" + templateName + ".txt")
	template, errTemplate := templateEngine.Parse(contentString)
	if errTemplate != nil {
		return &errTemplate
	}

	newbytes := bytes.NewBufferString("")
	errTemplateExecute := template.Execute(newbytes, data)
	if errTemplateExecute != nil {
		return &errTemplateExecute
	}

	// Set up authentication information.
	server := SMTPServer{
		Addr: beego.AppConfig.String("MailAddr"),
		Auth: smtp.PlainAuth(
			"",
			beego.AppConfig.String("MailUser"),
			beego.AppConfig.String("MailPassword"),
			beego.AppConfig.String("MailHost"),
		),
	}

	// Split the template
	tempOutput := newbytes.String()
	split := regexp.MustCompile("\n").Split(tempOutput, 3)

	mail := Mail{
		Server:  server,
		From:    netMail.Address{beego.AppConfig.String("MailFrom"), beego.AppConfig.String("MailFromAddress")},
		To:      netMail.Address{"", to},
		Subject: split[0],
		Message: split[2],
	}

	errMailSend := mail.Send()
	if errMailSend != nil {
		return &errMailSend
	}

	return nil
}
