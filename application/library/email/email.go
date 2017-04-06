package email

import (
	"errors"
	"net/smtp"
	"strconv"
	"time"

	"github.com/admpub/email"
	"github.com/admpub/log"
)

type queueItem struct {
	*email.Email
	Config Config
}

type SMTPConfig struct {
	Identity string
	Host     string
	Port     int
	Username string
	Password string
}

func (s *SMTPConfig) Address() string {
	if s.Port == 0 {
		s.Port = 25
	}
	return s.Host + `:` + strconv.Itoa(s.Port)
}

func (s *SMTPConfig) Auth() smtp.Auth {
	return smtp.PlainAuth(s.Identity, s.Username, s.Password, s.Host)
}

type Config struct {
	SMTP       *SMTPConfig
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

func Initial(queueSize int) {
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
		Initial(1)
	}
	if conf.SMTP == nil {
		return ErrSMTPNoSet
	}
	if conf.Auth == nil {
		conf.Auth = conf.SMTP.Auth()
	}
	mail := email.NewEmail()
	mail.From = conf.From
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
