package main

import (
	"fmt"
	"github.com/golangit/dic/container"
	"github.com/golangit/dic/reference"
)

// the mailer objects
type Mailer interface {
	Send(string) bool
}

type mailer struct {
	Sender     MailSender
	MailPrefix string
}

func MailerNew(sender MailSender, mailPrefix string) Mailer {
	return &mailer{
		Sender:     sender,
		MailPrefix: mailPrefix,
	}
}

// the sender objects

type MailSender interface{}

type sendmailer struct {
	Str string
}

type SendmailSender interface {
	MailSender
}

func SendmailNew() SendmailSender {
	return &sendmailer{
		Str: "transport is sendmail",
	}
}

func (m *mailer) Send(str string) bool {
	fmt.Printf("Sent `%s%s` using '%v'\n", m.MailPrefix, str, m.Sender)
	return true
}

func main() {
	fmt.Println("Hello, High quality dev")

	fmt.Println("Registering services")
	cnt := container.New()
	cnt.Register("transport.sendmail", SendmailNew /*, lot of arguments*/)
	//cnt.Register("transport.postfix", PostfixNew, "different args than sendmail")
	//cnt.Alias("transport", "transport.sendmail")
	cnt.Register("mailer", MailerNew, reference.New("transport.sendmail"), "[golangit] ")

	fmt.Println("Getting Mailer")

	fmt.Println("Calling...")
	cnt.Get("mailer").(Mailer).Send("liuggio")
}
