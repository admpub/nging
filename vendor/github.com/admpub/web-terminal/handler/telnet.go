package handler

import (
	"io"
	"net"
	"os"
	"runtime"
	"strings"

	"github.com/admpub/web-terminal/config"
	"github.com/admpub/web-terminal/library/telnet"
	"golang.org/x/net/websocket"
)

func TelnetShell(ws *websocket.Conn) {
	defer ws.Close()
	ctx := NewContext(ws)
	hostname := ParamGet(ctx, "hostname")
	port := ParamGet(ctx, "port")
	if 0 == len(port) {
		port = "23"
	}
	charset := ParamGet(ctx, "charset")
	if 0 == len(charset) {
		if "windows" == runtime.GOOS {
			charset = "GB18030"
		} else {
			charset = "UTF-8"
		}
	}
	//columns := toInt(ParamGet(ctx,"columns"), 80)
	//rows := toInt(ParamGet(ctx,"rows"), 40)

	var dumpOut io.WriteCloser
	var dumpIn io.WriteCloser

	client, err := net.Dial("tcp", hostname+":"+port)
	if nil != err {
		logString(ws, "Failed to dial: "+err.Error())
		return
	}
	defer func() {
		client.Close()
		if nil != dumpOut {
			dumpOut.Close()
		}
		if nil != dumpIn {
			dumpIn.Close()
		}
	}()

	debug := config.Default.Debug
	if "true" == strings.ToLower(ParamGet(ctx, "debug")) {
		debug = true
	}

	if debug {
		var err error
		dumpOut, err = os.OpenFile(config.Default.LogDir+hostname+".dump_telnet_out.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
		if nil != err {
			dumpOut = nil
		}
		dumpIn, err = os.OpenFile(config.Default.LogDir+hostname+".dump_telnet_in.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
		if nil != err {
			dumpIn = nil
		}
	}

	conn, e := telnet.NewConnWithRead(client, warp(client, dumpIn))
	if nil != e {
		logString(nil, "failed to create connection: "+e.Error())
		return
	}
	columns := toInt(ParamGet(ctx, "columns"), 80)
	rows := toInt(ParamGet(ctx, "rows"), 40)
	conn.SetWindowSize(byte(rows), byte(columns))

	go func() {
		_, err := io.Copy(decodeBy(charset, client), warp(ws, dumpOut))
		if nil != err {
			logString(nil, "copy of stdin failed:"+err.Error())
		}
	}()

	if _, err := io.Copy(decodeBy(charset, ws), conn); err != nil {
		logString(ws, "copy of stdout failed:"+err.Error())
		return
	}
}
