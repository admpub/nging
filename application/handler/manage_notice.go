package handler

import (
	"errors"

	"fmt"

	"github.com/admpub/log"
	"github.com/admpub/websocket"
	"github.com/webx-top/echo"
)

type OnlineUser struct {
	Message chan string
	Clients uint
}

func NewOnlineUser() *OnlineUser {
	return &OnlineUser{
		Message: make(chan string),
		Clients: 1,
	}
}

func IsOfflineUser(user string) bool {
	chanUsers[user].Clients--
	if chanUsers[user].Clients <= 0 {
		delete(chanUsers, user)
		return true
	}
	return false
}

var chanUsers = map[string]*OnlineUser{}

func SendNotice(user string, message string) {
	_, exists := chanUsers[user]
	if !exists {
		fmt.Println(message)
		return
	}
	chanUsers[user].Message <- message
}

func ManageNotice(c *websocket.Conn, ctx echo.Context) error {
	user, ok := ctx.Get(`user`).(string)
	if !ok {
		return errors.New(ctx.T(`登录信息获取失败，请重新登录`))
	}
	_, exists := chanUsers[user]
	if !exists {
		chanUsers[user] = NewOnlineUser()
	} else {
		chanUsers[user].Clients++
	}
	defer IsOfflineUser(user)
	//push(writer)
	go func() {
		for {
			message := <-chanUsers[user].Message
			log.Info(`Push message: `, message)
			if err := c.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				log.Error(`Push error: `, err.Error())
				IsOfflineUser(user)
				return
			}
		}
	}()

	//echo
	var execute = func(conn *websocket.Conn) error {
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				return err
			}

			if err = conn.WriteMessage(mt, message); err != nil {
				return err
			}
		}
	}
	err := execute(c)
	if err != nil {
		log.Error(err)
	}
	return nil
}
