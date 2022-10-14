package gopty

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"unicode/utf8"

	"github.com/admpub/gopty/interfaces"
)

// PTY2Websocket pty to websocket
func PTY2Websocket(ws WebsocketWriter, pty interfaces.Console) {
	buffer := make([]byte, 1024)
	var payload, overflow []byte
	for {
		n, err := pty.Read(buffer)
		if err != nil {
			fmt.Println("[PTY2Websocket] read from pty error: ", err)
			return
		}

		// Empty the overflow from the last read into the payload first.
		payload = append(payload, overflow...)
		overflow = nil
		// Then empty the new buf read into the payload.
		payload = append(payload, buffer[:n]...)

		// Strip out any incomplete utf-8 from current payload into overflow.
		for !utf8.Valid(payload) {
			overflow = append(overflow, payload[len(payload)-1])
			payload = payload[:len(payload)-1]
		}

		if len(payload) > 0 {
			err = ws.WriteMessage(BinaryMessage, payload[:])
			if err != nil {
				fmt.Println("[PTY2Websocket] write to ws error: ", err)
				return
			}
		}

		// Empty the payload.
		payload = nil
	}
}

// WebsocketWriter websocket writer
type WebsocketWriter interface {
	WriteMessage(int, []byte) error
}

// WebsocketReader websocket reader
type WebsocketReader interface {
	ReadMessage() (int, []byte, error)
}

// Websocketer websocket interface
// github.com/admpub/websocket
// github.com/gorilla/websocket
type Websocketer interface {
	WebsocketWriter
	WebsocketReader
}

// The message types are defined in RFC 6455, section 11.8.
const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
)

var (
	resizePrefix  = []byte("<RESIZE>")
	comma         = []byte(",")
	dangerCmdTips = []byte("Dangerous commands disabled: ")
	dangerRmrf    = regexp.MustCompile(`^[\s]*rm[\s]+-[\w]*r[\w]*[\s]+([\.]/?|/|(..[/]?)+)[*]*[\s]`)
)

func IsDangerCommand(b []byte) bool {
	return dangerRmrf.Match(b)
}

// Websocket2PTY websocket to pty
func Websocket2PTY(ws Websocketer, pty interfaces.Console) {
	for {
		mt, message, err := ws.ReadMessage()
		if mt == -1 || err != nil {
			log.Println("[Websocket2PTY] websocket read error: ", err)
			return
		}
		if bytes.HasPrefix(message, resizePrefix) {
			size := message[len(resizePrefix):]
			sizeArr := bytes.SplitN(size, comma, 2)
			if len(sizeArr) != 2 {
				err = ws.WriteMessage(BinaryMessage, message)
				if err != nil {
					log.Println("[Websocket2PTY] websocket write error: ", err)
				}
				continue
			}
			rows, _ := strconv.Atoi(string(sizeArr[0]))
			cols, _ := strconv.Atoi(string(sizeArr[1]))
			err = pty.SetSize(cols, rows)
			log.Printf("[Websocket2PTY] pty resize window to %d, %d\n", cols, rows)
			if err != nil {
				err = ws.WriteMessage(BinaryMessage, []byte(err.Error()))
				if err != nil {
					log.Println("[Websocket2PTY] websocket write error: ", err)
				}
			}
		} else if matches := dangerRmrf.FindAll(message, -1); len(matches) > 0 {
			tips := make([]byte, len(dangerCmdTips), len(dangerCmdTips)+len(matches[0]))
			copy(tips, dangerCmdTips)
			tips = append(tips, matches[0]...)
			err = ws.WriteMessage(BinaryMessage, tips)
			if err != nil {
				log.Println("[Websocket2PTY] websocket write error: ", err)
			}
		} else {
			_, err = pty.Write(message)
			if err != nil {
				log.Println("[Websocket2PTY] pty write error: ", err)
			}
		}
	}
}

var bash string
var flagVar string

func init() {
	if runtime.GOOS == "windows" {
		bash = "cmd.exe"
		flagVar = "/c"
	} else {
		shell := os.Getenv("SHELL")
		if len(shell) == 0 {
			shell = "/bin/bash"
			if _, err := os.Stat(shell); err != nil {
				shell = "/bin/sh"
			}
		}
		bash = shell
		flagVar = "-c"
	}
}

// GetBash get bash file
func GetBash() string {
	return bash
}

// GetFlagVar bash flag variable name
func GetFlagVar() string {
	return flagVar
}

// ServeWebsocket ServeWebsocket(wsc,120,60)
func ServeWebsocket(wsc Websocketer, cols, rows int) error {
	pty, err := New(cols, rows)
	if err != nil {
		return err
	}
	defer pty.Close()
	args := []string{bash}
	err = pty.Start(args)
	if err != nil {
		err = fmt.Errorf("[gopty] open terminal err: %w", err)
		return err
	}

	go PTY2Websocket(wsc, pty)
	// block from close
	Websocket2PTY(wsc, pty)
	return nil
}

// Execute execute command
func Execute(command string, resultWriter io.Writer) error {
	var cols, rows = 120, 60
	pty, err := New(cols, rows)
	if err != nil {
		return err
	}
	defer pty.Close()
	args := []string{GetBash(), GetFlagVar(), command}
	err = pty.Start(args)
	if err != nil {
		err = fmt.Errorf("[gopty] open terminal err: %w", err)
		return err
	}
	go func() {
		_, err = io.Copy(resultWriter, pty)
		if err != nil {
			log.Printf("[gopty] Error: %v\n", err)
		}
	}()
	_, err = pty.Wait()
	return err
}
