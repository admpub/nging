package handler

import (
	"time"

	"github.com/admpub/log"
	"github.com/admpub/sockjs-go/sockjs"
	"github.com/webx-top/echo"
	sockjsHandler "github.com/webx-top/echo/handler/sockjs"
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
			if counter >= 10 { //测试只推10条
				return
			}
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
	sockjsHandler.DefaultExecuter(c)
	return nil
}
