package handler

import (
	"time"

	"github.com/admpub/log"
	"github.com/admpub/sockjs-go/sockjs"
	"github.com/admpub/websocket"
	"github.com/webx-top/echo"
)

func ManageExeCMD(ctx echo.Context) error {
	var err error
	return ctx.Render(`manage/execmd`, err)
}

func SockJSManageExeCMDSend(c sockjs.Session) error {
	//push(writer)
	go func() {
		var counter int
		for {
			time.Sleep(5 * time.Second)
			message := time.Now().String()
			log.Info(`Push message: `, message)
			if err := c.Send(message); err != nil {
				log.Error(`Push error: `, err.Error())
				return
			}
			counter++
		}
	}()

	//echo
	var execute = func(session sockjs.Session) error {
		for {
			msg, err := session.Recv()
			log.Info(`Recv message: `, msg)
			if err != nil {
				return err
			}
			err = session.Send(msg)
			if err != nil {
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

func WSManageExeCMDSend(c *websocket.Conn, ctx echo.Context) error {
	//push(writer)
	go func() {
		var counter int
		for {
			time.Sleep(5 * time.Second)
			message := time.Now().String()
			log.Info(`Push message: `, message)
			if err := c.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				log.Error(`Push error: `, err.Error())
				return
			}
			counter++
		}
	}()

	//echo
	var execute = func(conn *websocket.Conn) error {
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				return err
			}
			log.Infof("Websocket recv: %s", message)

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
