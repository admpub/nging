package email

import (
	"errors"
	"net/smtp"
	"time"

	"strings"

	"github.com/admpub/email"
	"github.com/admpub/log"
	"github.com/admpub/nging/application/library/config"
)

type queueItem struct {
	*email.Email
	Config Config
}

type Config struct {
	SMTP       *config.SMTPConfig
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
)

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
				if err := m.Email.Send(m.Config.SMTP.Address(), m.Config.Auth); err != nil {
					log.Error("SendMail:", err.Error())
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
	mail := email.NewEmail()
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
	item := &queueItem{Email: mail, Config: *conf}
	select {
	case sendCh <- item:
		return nil
	case <-time.After(time.Second * 3):
		return ErrSendChannelTimeout
	}
}
