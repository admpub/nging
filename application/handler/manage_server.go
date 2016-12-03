package handler

import (
	"os/exec"
	"strings"

	"github.com/admpub/log"
	"github.com/admpub/sockjs-go/sockjs"
	"github.com/admpub/websocket"
	"github.com/webx-top/echo"
)

type CMDResultCapturer struct {
	Do func([]byte) error
}

func (this CMDResultCapturer) Write(p []byte) (n int, err error) {
	err = this.Do(p)
	n = len(p)
	return
}

func ManageExeCMD(ctx echo.Context) error {
	var err error
	return ctx.Render(`manage/execmd`, err)
}

func runCMD(command string, recvResult func([]byte) error) {
	params := strings.Split(command, ` `)
	length := len(params)
	if len(params[0]) == 0 {
		return
	}
	var cmd *exec.Cmd
	if length > 1 {
		cmd = exec.Command(params[0], params[1:]...)
	} else {
		cmd = exec.Command(params[0])
	}
	out := CMDResultCapturer{Do: recvResult}
	cmd.Stdout = out
	cmd.Stderr = out

	go func() {
		err := cmd.Run()
		if err != nil {
			recvResult([]byte(err.Error()))
		}
	}()
}

func SockJSManageExeCMDSend(c sockjs.Session) error {
	send := make(chan func() string)
	//push(writer)
	go func() {
		for {
			fn := <-send
			message := fn()
			log.Info(`Push message: `, message)
			if err := c.Send(message); err != nil {
				log.Error(`Push error: `, err.Error())
				return
			}
		}
	}()

	//echo
	var execute = func(session sockjs.Session) error {
		for {
			command, err := session.Recv()
			if err != nil {
				return err
			}
			if len(command) > 0 {
				runCMD(command, func(b []byte) error {
					send <- func() string {
						return string(b)
					}
					return nil
				})
			}
			err = session.Send(command)
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
	send := make(chan func() string)
	//push(writer)
	go func() {
		for {
			fn := <-send
			message := fn()
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
			command := string(message)
			if len(command) > 0 {
				runCMD(command, func(b []byte) error {
					send <- func() string {
						return string(b)
					}
					return nil
				})
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
