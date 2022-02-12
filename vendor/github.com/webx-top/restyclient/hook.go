package restyclient

import (
	"fmt"
	"strings"

	"github.com/admpub/log"
	"github.com/admpub/resty/v2"
)

func InitRestyHook(client *resty.Client) {
	client.OnBeforeRequest(func(_ *resty.Client, r *resty.Request) error {
		log.Debugf("[resty] %s %s", r.Method, r.URL)
		OutputMaps("Request header", r.Header)
		return nil
	})
	client.OnAfterResponse(func(_ *resty.Client, r *resty.Response) error {
		OutputMaps("Response header", r.Header())
		log.Debugf("[resty] %s %s", r.Proto(), r.Status())
		return nil
	})
	client.OnError(func(r *resty.Request, err error) {
		log.Errorf("[resty] %s %s: %v", r.Method, r.URL, err)
	})
}

// OutputMaps Just debug a map
func OutputMaps(info string, args map[string][]string) {
	if !log.IsEnabled(log.LevelDebug) {
		return
	}
	s := "\n"
	for k, v := range args {
		s = s + fmt.Sprintf("%-25s| %-6s\n", k, strings.Join(v, "||"))
	}
	log.Debugf("[resty] %s %s", info, s)
}
