package main

import (
	"html/template"
	"time"

	"github.com/admpub/websocket"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine/fasthttp"
	"github.com/webx-top/echo/engine/standard"
	ws "github.com/webx-top/echo/handler/websocket"
	mw "github.com/webx-top/echo/middleware"
)

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.echo}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            if (typeof(evt.data)!='undefined') {
                print("ERROR: " + evt.data);
            }else{
                console.dir(evt);
            }
        }
        return false;
    };

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

    var wsn = new WebSocket("{{.notice}}");
    var notice = document.getElementById("notice");
    wsn.onopen = function(evt) {
        notice.innerHTML = "[NOTICE] OPEN";
    }
    wsn.onclose = function(evt) {
        notice.innerHTML = "[NOTICE] CLOSE";
        wsn = null;
    }
    wsn.onmessage = function(evt) {
        notice.innerHTML = "[NOTICE] RESPONSE: " + evt.data;
    }
    wsn.onerror = function(evt) {
        notice.innerHTML = "[NOTICE] ERROR: " + (typeof(evt.data)!='undefined'?evt.data:evt);
    }
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="notice" style="color:red"></div>
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))

func main() {
	e := echo.New()
	e.Use(mw.Log())

	e.Get("/", func(c echo.Context) error {
		homeTemplate.Execute(c.Response(), map[string]string{
			"echo":   "ws://" + c.Request().Host() + "/websocket",
			"notice": "ws://" + c.Request().Host() + "/notice",
		})
		return nil
	})
	//e.Get("/websocket", ws.Websocket(nil))

	e.HandlerWrapper = ws.HanderWrapper

	e.Get("/websocket", func(c *websocket.Conn, ctx echo.Context) error {
		//push(writer)
		go func() {
			var counter int
			for {
				if counter >= 10 { //测试只推10条
					return
				}
				time.Sleep(5 * time.Second)
				message := time.Now().String()
				ctx.Logger().Info(`Push message: `, message)
				if err := c.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					ctx.Logger().Error(`Push error: `, err.Error())
					return
				}
				counter++
			}
		}()

		//echo
		ws.DefaultExecuter(c, ctx)
		return nil
	})

	e.Get("/notice", func(c *websocket.Conn, ctx echo.Context) error {
		for {
			time.Sleep(5 * time.Second)
			message := time.Now().String()
			ctx.Logger().Info(`Push message: `, message)
			if err := c.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				return err
			}
		}
		return nil
	})

	switch `fast` {
	case `fast`:

		// FastHTTP
		e.Run(fasthttp.New(":4444"))

	default:
		// Standard
		e.Run(standard.New(":4444"))
	}
}
