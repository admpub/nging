package cron

import "github.com/webx-top/echo/param"

var (
	// Senders 发信程序
	senders = []func(param.Store) error{}
)

// AddSender 添加发信程序
func AddSender(sender func(params param.Store) error) {
	senders = append(senders, sender)
}

// Send 发送通知/信件
func Send(params param.Store) (err error) {
	for _, sender := range senders {
		err = sender(params)
		if err != nil {
			return err
		}
	}
	return err
}
