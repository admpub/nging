package cron

import "github.com/admpub/nging/v4/application/registry/alert"

var (
	// Senders 发信程序
	senders = []func(alertData *alert.AlertData) error{}
)

// AddSender 添加发信程序
func AddSender(sender func(alertData *alert.AlertData) error) {
	senders = append(senders, sender)
}

// Send 发送通知/信件
func Send(alertData *alert.AlertData) (err error) {
	for _, sender := range senders {
		err = sender(alertData)
		if err != nil {
			return err
		}
	}
	return err
}
