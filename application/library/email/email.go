package email

import (
	"crypto/tls"
	"errors"
	"net/smtp"
	"strings"
	"time"

	"github.com/admpub/email"
	"github.com/admpub/log"
	"github.com/admpub/mail"
	"github.com/admpub/nging/application/library/config"
)

type queueItem struct {
	Email  *email.Email
	Config Config
}

func (q *queueItem) send1() error {
	if q.Config.SMTP.Secure == "SSL" || q.Config.SMTP.Secure == "TLS" {
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         q.Config.SMTP.Host,
		}
		return q.Email.SendWithTLS(q.Config.SMTP.Address(), q.Config.Auth, tlsconfig)
	}
	return q.Email.Send(q.Config.SMTP.Address(), q.Config.Auth)
}

func (q *queueItem) send2() error {
	return mail.SendMail(q.Config.Subject, string(q.Config.Content), q.Config.ToAddress, q.Config.From,
		q.Config.CcAddress, q.Config.SMTP, nil)
}

func (q *queueItem) Send() (err error) {
	if q.Email == nil {
		return q.send2()
	}
	if config.DefaultConfig.Email.Timeout <= 0 {
		return q.send1()
	}
	done := make(chan bool)
	go func() {
		err = q.send1()
		done <- true
	}()
	select {
	case <-done:
		return
	case <-time.After(time.Second * time.Duration(config.DefaultConfig.Email.Timeout)):
		log.Error("发送邮件超时，采用备用方案发送")
		close(done)
	}
	return q.send2()
}

type Config struct {
	SMTP       *mail.SMTPConfig
	From       string
	ToAddress  string
	ToUsername string
	Subject    string
	Content    []byte
	CcAddress  []string
	Auth       smtp.Auth
}

var (
	sendCh                chan *queueItem
	ErrSMTPNoSet          = errors.New(`SMTP is not set`)
	ErrSendChannelTimeout = errors.New(`SendMail: The sending channel timed out`)
	smtpClient            *mail.SMTPClient
)

func SMTPClient(conf *mail.SMTPConfig) *mail.SMTPClient {
	if smtpClient == nil {
		c := mail.NewSMTPClient(conf)
		smtpClient = &c
	}
	return smtpClient
}

func Initial(queueSizes ...int) {
	var queueSize int
	if len(queueSizes) > 0 {
		queueSize = queueSizes[0]
	} else {
		queueSize = config.DefaultConfig.Email.QueueSize
	}
	if sendCh != nil {
		close(sendCh)
	}
	if queueSize <= 0 {
		queueSize = 1
	}
	sendCh = make(chan *queueItem, queueSize)
	go func() {
		for {
			select {
			case m, ok := <-sendCh:
				if !ok {
					return
				}
				log.Info("<SendMail> Sending: ", m.Config.ToAddress)
				err := m.Send()
				if err != nil {
					log.Error("<SendMail> Error: ", err.Error())
				} else {
					log.Info("<SendMail> Result: ", m.Config.ToAddress, " [OK]")
				}
			}
		}
	}()
}

func SendMail(conf *Config) error {
	if sendCh == nil {
		Initial()
	}
	if conf.SMTP == nil {
		return ErrSMTPNoSet
	}
	if conf.Auth == nil {
		conf.Auth = conf.SMTP.Auth()
	}
	var mail *email.Email
	if config.DefaultConfig.Email.Engine == `email` {
		mail = email.NewEmail()
		mail.From = conf.From
		if len(mail.From) == 0 {
			mail.From = conf.SMTP.Username
			if !strings.Contains(mail.From, `@`) {
				mail.From += `@` + conf.SMTP.Host
			}
		}
		mail.To = []string{conf.ToAddress}
		mail.Subject = conf.Subject
		mail.HTML = conf.Content
		if len(conf.CcAddress) > 0 {
			mail.Cc = conf.CcAddress
		}
	}
	item := &queueItem{Email: mail, Config: *conf}
	select {
	case sendCh <- item:
		return nil
	case <-time.After(time.Second * 3):
		return ErrSendChannelTimeout
	}
}
