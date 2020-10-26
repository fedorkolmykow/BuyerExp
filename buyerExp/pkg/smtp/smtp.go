package smtp

import (
	"crypto/tls"
	"encoding/base64"
	"gopkg.in/gomail.v2"
	"os"

	"github.com/fedorkolmykow/avitoexp/pkg/api"

	log "github.com/sirupsen/logrus"
)

const confirmationMessage = "Follow this link to confirm your email:\n "
const newPriceMessage = "There is a new price for this notice:\n "
const schema = "http://"

type smtpClient struct{
	dialer *gomail.Dialer
}

type SmtpClient interface {
	SendConfirmationMail(user *api.User) (err error)
	SendMailsWithNewPrices(subs []api.Subscription)
}

func (s *smtpClient) SendMailsWithNewPrices(subs []api.Subscription) {
	for _, sub := range subs{
		m := gomail.NewMessage()
		m.SetHeader("From", "exp@avito.com")
		m.SetHeader("To", sub.User.Mail)
		m.SetHeader("Subject", "Notice's new price")
		mes := newPriceMessage + sub.Notice.URL
		m.SetBody("text/html", mes)

		err := s.dialer.DialAndSend(m)
		if err != nil{
			log.Warn(err)
		}
	}
	return
}


func (s *smtpClient) SendConfirmationMail(user *api.User) (err error){
	m := gomail.NewMessage()
	m.SetHeader("From", "exp@avito.com")
	m.SetHeader("To", user.Mail)
	m.SetHeader("Subject", "Email confirmation")
	mes := confirmationMessage +
			schema +
			os.Getenv("HOST") +
			os.Getenv("HTTP_PORT") +
			api.Сonfirmation +
			"?hash=" +
			base64.URLEncoding.EncodeToString(user.Hash)
	m.SetBody("text/html", mes)

	err = s.dialer.DialAndSend(m)
	return
}


func NewSMTP() SmtpClient{
	d := gomail.NewDialer(
		"exppostfix",
		25,
		os.Getenv("MAIL_USER"),
		os.Getenv("MAIL_PASSWORD"))
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	s := &smtpClient{
		dialer: d,
	}
	return s
}