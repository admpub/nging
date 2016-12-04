package handler

import (
	"errors"

	"github.com/admpub/log"
	"github.com/admpub/websocket"
	"github.com/webx-top/echo"
)

var chanUsers = map[string]chan string{}

func SendNotice(user string, message string) {
	_, exists := chanUsers[user]
	if !exists {
		chanUsers[user] = make(chan string)
	}
	chanUsers[user] <- message
}

func ManageNotice(c *websocket.Conn, ctx echo.Context) error {
	user, ok := ctx.Get(`user`).(string)
	if !ok {
		return errors.New(ctx.T(`登录信息获取失败，请重新登录`))
	}
	_, exists := chanUsers[user]
	if !exists {
		chanUsers[user] = make(chan string)
	}
	//push(writer)
	go func() {
		for {
			message := <-chanUsers[user]
			log.Info(`Push message: `, message)
			if err := c.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				log.Error(`Push error: `, err.Error())
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
