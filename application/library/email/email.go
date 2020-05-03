/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

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
	"github.com/admpub/nging/application/library/notice"
)

type queueItem struct {
	Email  *email.Email
	Config Config
}

func (q *queueItem) send1() error {
	log.Info(`<SendMail> Using: send1`)
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
	log.Info(`<SendMail> Using: send2`)
	return mail.SendMail(
		q.Config.Subject,
		string(q.Config.Content),
		q.Config.ToAddress,
		q.Config.From,
		q.Config.CcAddress,
		q.Config.SMTP,
		nil,
	)
}

func (q *queueItem) Send() (err error) {
	if q.Config.Timeout <= 0 {
		if q.Email == nil {
			return q.send2()
		}
		return q.send1()
	}
	done := make(chan bool)
	go func() {
		err = q.send1()
		done <- true
	}()
	t := time.NewTicker(time.Second * time.Duration(q.Config.Timeout))
	defer t.Stop()
	select {
	case <-done:
		return
	case <-t.C:
		log.Error("发送邮件超时，采用备用方案发送")
		close(done)
	}
	return q.send2()
}

var Callbacks = []func(*Config, error){}

func AddCallback(cb func(*Config, error)) {
	Callbacks = append(Callbacks, cb)
}

type Config struct {
	ID         uint64 //RequestID
	Engine     string
	SMTP       *mail.SMTPConfig
	From       string
	ToAddress  string
	ToUsername string
	Subject    string
	Content    []byte
	CcAddress  []string
	Auth       smtp.Auth
	Timeout    int64
	Noticer    notice.Noticer
}

var (
	sendCh                chan *queueItem
	ErrSMTPNoSet          = errors.New(`SMTP is not set`)
	ErrSenderNoSet        = errors.New(`The sender address is not set`)
	ErrRecipientNoSet     = errors.New(`The recipient address is not set`)
	ErrSendChannelTimeout = errors.New(`SendMail: The sending channel timed out`)
	smtpClient            *mail.SMTPClient
	//QueueSize 每批容量
	QueueSize = 50
	//QueueWait 队列等待时间
	QueueWait = time.Second * 30
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
	}
	if sendCh != nil {
		close(sendCh)
	}
	if queueSize <= 0 {
		queueSize = QueueSize
	}
	sendCh = make(chan *queueItem, queueSize)
	go func() {
		for {
			select {
			case m, ok := <-sendCh:
				if !ok {
					return
				}
				noticer := m.Config.Noticer
				if noticer == nil {
					noticer = notice.DefaultNoticer
				}
				noticer("<SendMail> Sending: "+m.Config.ToAddress, 1)
				err := m.Send()
				if err != nil {
					noticer("<SendMail> Result: "+m.Config.ToAddress+" Error: "+err.Error(), 0)
				} else {
					noticer("<SendMail> Result: "+m.Config.ToAddress+" [OK]", 1)
				}
				for _, callback := range Callbacks {
					callback(&m.Config, err)
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
	if len(conf.SMTP.Host) == 0 || len(conf.SMTP.Username) == 0 {
		return ErrSMTPNoSet
	}
	if len(conf.From) == 0 {
		return ErrSenderNoSet
	}
	if len(conf.ToAddress) == 0 {
		return ErrRecipientNoSet
	}
	if conf.Auth == nil {
		conf.Auth = conf.SMTP.Auth()
	}
	var mail *email.Email
	if conf.Engine == `email` || conf.Engine == `send1` {
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
	t := time.NewTicker(QueueWait)
	defer t.Stop()
	select {
	case sendCh <- item:
		return nil
	case <-t.C:
		return ErrSendChannelTimeout
	}
}
