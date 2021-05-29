package subdomains

import (
	"strings"

	"github.com/admpub/log"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/engine/fasthttp"
	"github.com/webx-top/echo/engine/standard"
)

var Default = New()

func New() *Subdomains {
	s := &Subdomains{
		Hosts:    map[string]string{},
		Alias:    map[string]*Info{},
		Default:  ``,
		Protocol: `http`,
	}
	s.dispatcher = s.DefaultDispatcher
	return s
}

func SetBaseURL(name string, url string) {
	echo.Set(`subdomains.`+name+`.url`, strings.TrimRight(url, `/`))
}

type Info struct {
	Protocol string
	Name     string
	Host     string
	*echo.Echo
}

func (info *Info) URL(s *Subdomains, uri string) string {
	if domain := echo.String(`subdomains.` + info.Name + `.url`); len(domain) > 0 {
		return domain + info.Prefix() + uri
	}
	if len(info.Host) == 0 {
		if s.Default == info.Name {
			return info.Prefix() + uri
		}
		return info.Prefix() + `/` + info.Name + uri
	}
	protocol := info.Protocol
	if len(protocol) == 0 {
		protocol = s.Protocol
		if len(protocol) == 0 {
			protocol = `http`
		}
	}
	return protocol + `://` + info.Host + info.Prefix() + uri
}

func (info *Info) URLByName(s *Subdomains, name string, args ...interface{}) string {
	return info.URL(s, info.Echo.URI(name, args...))
}

type Dispatcher func(r engine.Request, w engine.Response) (*echo.Echo, bool)

type Subdomains struct {
	Hosts      map[string]string //{host:name}
	Alias      map[string]*Info
	Default    string //default name
	Protocol   string //http/https
	Boot       string
	dispatcher Dispatcher
}

func (s *Subdomains) SetDispatcher(dispatcher Dispatcher) *Subdomains {
	s.dispatcher = dispatcher
	return s
}

// Add 添加子域名，name的值支持以下三种格式：
// 1. 别名@域名 ———— 一个别名可以对应多个域名，每个域名之间用半角逗号“,”分隔
// 2. 域名 ———— 可以添加多个域名，每个域名之间用半角逗号“,”分隔。这里会自动将第一个域名中的首个点号“.”前面的部分作为别名，例如“blog.webx.top,news.webx.top”会自动将“blog”作为别名
// 3. 别名 ———— 在不指定域名的情况下将无法访问，除非“默认访问别名”（Subdomains.Default的值）与此相同
func (s *Subdomains) Add(name string, e *echo.Echo) *Subdomains {
	log.Debug(`[Subdomains] `, `Register subdomain: `, name)
	r := strings.SplitN(name, `@`, 2)
	var hosts []string
	if len(r) > 1 { //blog@1.webx.top,2.webx.top
		name = r[0]
		hosts = strings.Split(r[1], `,`)
	} else {
		p := strings.Index(name, `.`)
		if p > 0 { //blog.webx.top
			hosts = strings.Split(name, `,`)
			name = name[0:p]
		} else { //blog
			hosts = append(hosts, ``)
		}
	}
	for _, host := range hosts {
		s.Hosts[host] = name
	}
	info := &Info{
		Protocol: `http`,
		Name:     name,
		Host:     hosts[0],
		Echo:     e,
	}
	info2 := strings.SplitN(info.Host, `://`, 2)
	if len(info2) == 2 {
		info.Protocol = info2[0]
		info.Host = info2[1]
		if len(info.Protocol) == 0 {
			info.Protocol = "http"
		}
	}
	s.Alias[name] = info
	return s
}

func (s *Subdomains) Get(args ...string) *Info {
	name := s.Default
	if len(args) > 0 {
		name = args[0]
	}
	if e, ok := s.Alias[name]; ok {
		return e
	}
	return nil
}

func (s *Subdomains) SetDebug(on bool) *Subdomains {
	for _, info := range s.Alias {
		info.SetDebug(on)
	}
	return s
}

// URL 多域名场景下的网址生成功能
// URL(网址路径,域名别名)，如果这里不传递域名别名，将使用默认别名的域名
// 例如：URL("/list?cid=1","blog")
// 对于一个别名对应有多个域名的情况，将总是使用第一个域名
func (s *Subdomains) URL(uri string, args ...string) string {
	info := s.Get(args...)
	if info == nil {
		return uri
	}
	return info.URL(s, uri)
}

// URLByName 根据路由别名生成网址
// 可以在名称中采用 #backend#name 的方式来获取子域名别名为bakcend的网址
func (s *Subdomains) URLByName(name string, params ...interface{}) string {
	var args []string
	if strings.HasPrefix(name, `#`) {
		name = strings.TrimPrefix(name, `#`)
		arr := strings.SplitN(name, `#`, 2)
		if len(arr) == 2 {
			args = append(args, arr[0])
			name = arr[1]
		}
	}
	info := s.Get(args...)
	if info == nil {
		return `/not-found:` + name
	}
	return info.URLByName(s, name, params...)
}

func (s *Subdomains) FindByDomain(host string) (*echo.Echo, bool) {
	name, exists := s.Hosts[host]
	if !exists {
		if p := strings.LastIndexByte(host, ':'); p > -1 {
			name, exists = s.Hosts[host[0:p]]
		}
		if !exists {
			name = s.Default
		}
	}
	var info *Info
	info, exists = s.Alias[name]
	if exists {
		return info.Echo, true
	}
	return nil, false
}

func (s *Subdomains) DefaultDispatcher(r engine.Request, w engine.Response) (*echo.Echo, bool) {
	return s.FindByDomain(r.Host())
}

func (s *Subdomains) ServeHTTP(r engine.Request, w engine.Response) {
	handler, exists := s.dispatcher(r, w)
	if exists && handler != nil {
		handler.ServeHTTP(r, w)
	} else {
		w.NotFound()
	}
}

func (s *Subdomains) Run(args ...interface{}) {
	if s.dispatcher == nil {
		s.dispatcher = s.DefaultDispatcher
	}
	var eng engine.Engine
	var arg interface{}
	size := len(args)
	if size > 0 {
		arg = args[0]
	}
	if size > 1 {
		if conf, ok := arg.(*engine.Config); ok {
			if v, ok := args[1].(string); ok {
				if v == `fast` {
					eng = fasthttp.NewWithConfig(conf)
				} else {
					eng = standard.NewWithConfig(conf)
				}
			} else {
				eng = fasthttp.NewWithConfig(conf)
			}
		} else {
			addr := `:80`
			if v, ok := arg.(string); ok && len(v) > 0 {
				addr = v
			}
			if v, ok := args[1].(string); ok {
				if v == `fast` {
					eng = fasthttp.New(addr)
				} else {
					eng = standard.New(addr)
				}
			} else {
				eng = fasthttp.New(addr)
			}
		}
	} else {
		switch v := arg.(type) {
		case string:
			eng = fasthttp.New(v)
		case engine.Engine:
			eng = v
		default:
			eng = fasthttp.New(`:80`)
		}
	}
	e := s.Get(s.Boot)
	if e == nil {
		for _, info := range s.Alias {
			e = info
			break
		}
	}
	for _, info := range s.Alias {
		if e == info {
			continue
		}
		info.Commit()
	}
	e.Logger().Info(`Server has been launched.`)
	e.Run(eng, s)
	e.Logger().Info(`Server has been closed.`)
}
