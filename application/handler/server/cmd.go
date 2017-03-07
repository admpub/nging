/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package server

import (
	"io"
	"os"
	"runtime"

	"os/exec"

	"strings"

	"time"

	"github.com/admpub/log"
	"github.com/admpub/nging/application/library/charset"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/sockjs-go/sockjs"
	"github.com/admpub/websocket"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

var (
	WebSocketLogger = log.GetLogger(`websocket`)
	IsWindows       bool
)

func init() {
	WebSocketLogger.SetLevel(`Info`)
	IsWindows = runtime.GOOS == `windows`
}

func Cmd(ctx echo.Context) error {
	var err error
	return ctx.Render(`manage/execmd`, err)
}

func CmdSendBySockJS(c sockjs.Session) error {
	send := make(chan string)
	//push(writer)
	go func() {
		for {
			message := <-send
			WebSocketLogger.Debug(`Push message: `, message)
			if err := c.Send(message); err != nil {
				WebSocketLogger.Error(`Push error: `, err.Error())
				return
			}
		}
	}()
	timeout := c.Request().URL.Query().Get(`timeout`)
	duration := config.ParseTimeDuration(timeout)
	//echo
	exec := func(session sockjs.Session) error {
		var (
			w   io.WriteCloser
			cmd *exec.Cmd
		)
		for {
			command, err := session.Recv()
			if err != nil {
				return err
			}
			if len(command) == 0 {
				continue
			}
			if w == nil {
				w, cmd, err = cmdRunner(command, send, func() {
					w.Close()
					w = nil
				}, duration)
				if err != nil {
					return err
				}
				continue
			}
			err = cmdContinue(command, w, cmd)
			if err != nil {
				return err
			}
		}
	}
	err := exec(c)
	if err != nil {
		WebSocketLogger.Error(err)
	}
	return nil
}

func cmdRunner(command string, send chan string, onEnd func(), timeout time.Duration) (w io.WriteCloser, cmd *exec.Cmd, err error) {
	cmd = com.CreateCmdStr(command, func(b []byte) (e error) {
		if IsWindows {
			b, e = charset.Convert(`gbk`, `utf-8`, b)
			if e != nil {
				return e
			}
		}
		send <- string(b)
		return nil
	})
	w, err = cmd.StdinPipe()
	if err != nil {
		return
	}
	done := make(chan struct{})
	go func() {
		log.Info(`[command] running: `, command)
		if e := cmd.Run(); e != nil {
			cmd.Stderr.Write([]byte(e.Error()))
		}
		done <- struct{}{}
		onEnd()
	}()
	cmdMaxTimeout := config.DefaultConfig.Sys.CmdTimeoutDuration
	if timeout <= 0 {
		timeout = time.Minute * 5
	}
	if timeout > cmdMaxTimeout {
		timeout = cmdMaxTimeout
	}
	go func() {
		ticker := time.NewTicker(timeout)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				log.Info(`[command] exited: `, command)
				return
			case <-ticker.C:
				if cmd == nil {
					return
				}
				if cmd.Process == nil {
					return
				}
				cmd.Stderr.Write([]byte(`timeout`))
				log.Info(`[command] timeout: `, command)
				err := cmd.Process.Kill()
				if err != nil {
					log.Error(err)
				}
				return
			}
		}
	}()
	return
}

func cmdContinue(command string, w io.WriteCloser, cmd *exec.Cmd) (err error) {
	switch command {
	case `^C`:
		err = cmd.Process.Signal(os.Interrupt)
		if err != nil {
			if !strings.HasPrefix(err.Error(), `not supported by `) {
				WebSocketLogger.Error(err)
			}
			err = cmd.Process.Kill()
			if err != nil {
				WebSocketLogger.Error(err)
			}
		}
	default:
		w.Write([]byte(command + "\n"))
	}
	return nil
}

func CmdSendByWebsocket(c *websocket.Conn, ctx echo.Context) error {
	send := make(chan string)
	//push(writer)
	go func() {
		for {
			message := <-send
			WebSocketLogger.Debug(`Push message: `, message)
			if err := c.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				WebSocketLogger.Error(`Push error: `, err.Error())
				return
			}
		}
	}()

	timeout := ctx.Query(`timeout`)
	duration := config.ParseTimeDuration(timeout)
	//echo
	exec := func(conn *websocket.Conn) error {
		var (
			w   io.WriteCloser
			cmd *exec.Cmd
		)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return err
			}
			command := string(message)
			if len(command) == 0 {
				continue
			}
			if w == nil {
				w, cmd, err = cmdRunner(command, send, func() {
					w.Close()
					w = nil
				}, duration)
				if err != nil {
					return err
				}
				continue
			}
			err = cmdContinue(command, w, cmd)
			if err != nil {
				return err
			}
		}
	}
	err := exec(c)
	if err != nil {
		WebSocketLogger.Error(err)
	}
	return nil
}
