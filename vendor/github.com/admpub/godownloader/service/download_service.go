package service

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
	sockjsHandler "github.com/webx-top/echo/handler/sockjs"

	"github.com/admpub/godownloader/httpclient"
	"github.com/admpub/sockjs-go/sockjs"
)

var savePath string

func GetDownloadPath() string {
	if len(savePath) > 0 {
		return savePath
	}
	var homeDir string
	usr, err := user.Current()
	if err != nil {
		log.Println(err)
		homeDir = os.Getenv("HOME")
	} else if usr != nil {
		homeDir = usr.HomeDir
	}
	if len(homeDir) == 0 {
		homeDir = filepath.Dir(os.Args[0])
	}
	st := strconv.QuoteRune(os.PathSeparator)
	st = st[1 : len(st)-1]
	sv := homeDir + st + "Downloads" + st + "GoDownloader" + st
	fi, err := os.Stat(sv)
	if err != nil || !fi.IsDir() {
		os.MkdirAll(sv, 0755)
	}
	savePath = sv
	return sv
}

type DJob struct {
	Id         int
	FileName   string
	Size       int64
	Downloaded int64
	Progress   int64
	Speed      int64
	State      string
}

type NewJob struct {
	Url       string
	PartCount int64
	FilePath  string
	Pipes     []string
}

type DServ struct {
	dls      []*httpclient.Downloader
	oplock   sync.Mutex
	tmpl     string
	savePath func() string
}

func (srv *DServ) SetTmpl(tmpl string) *DServ {
	srv.tmpl = tmpl
	return srv
}

func (srv *DServ) Tmpl() string {
	return srv.tmpl
}

func (srv *DServ) SetSavePath(savePath func() string) *DServ {
	srv.savePath = savePath
	return srv
}

func (srv *DServ) SavePath() func() string {
	if srv.savePath == nil {
		return GetDownloadPath
	}
	return srv.savePath
}

func (srv *DServ) Register(r echo.RouteRegister, enableSockJS bool) {
	if len(srv.tmpl) == 0 {
		srv.tmpl = "index"
	}
	r.Route("GET", "/", srv.index)
	r.Route("GET", "/index.html", srv.index)
	r.Route("GET,POST", "/progress.json", srv.progressJson)
	r.Route("GET,POST", "/add_task", srv.addTask)
	r.Route("GET,POST", "/remove_task", srv.removeTask)
	r.Route("GET,POST", "/start_task", srv.startTask)
	r.Route("GET,POST", "/stop_task", srv.stopTask)
	r.Route("GET,POST", "/start_all_task", srv.startAllTask)
	r.Route("GET,POST", "/stop_all_task", srv.stopAllTask)
	if !enableSockJS {
		return
	}
	sockjsOpts := sockjsHandler.Options{
		Handle: srv.ProgressSockJS,
		Prefix: "/progress",
	}
	sockjsOpts.Wrapper(r)
}

func (srv *DServ) SaveSettings(sf string) error {
	var ss ServiceSettings
	for _, i := range srv.dls {

		ss.Ds = append(ss.Ds, DownloadSettings{
			FI: i.Fi,
			Dp: i.GetProgress(),
		})
	}

	return ss.SaveToFile(sf)
}

func (srv *DServ) LoadSettings(sf string) error {
	ss, err := LoadFromFile(sf)
	if err != nil {
		log.Println("error: when try load settings", err)
		return err
	}
	log.Println(`settings:`, ss)
	for _, r := range ss.Ds {
		dl, err := httpclient.RestoreDownloader(r.FI.Url, r.FI.FileName, r.Dp, srv.SavePath(), PipeGetList(r.FI.Pipes...)...)
		if err != nil {
			return err
		}
		srv.dls = append(srv.dls, dl)
	}
	return nil
}

func (srv *DServ) index(ctx echo.Context) error {
	ctx.Set(`pipes`, PipeList())
	return ctx.Render(srv.tmpl, nil)
}

func (srv *DServ) addTask(ctx echo.Context) error {
	srv.oplock.Lock()
	defer srv.oplock.Unlock()
	var nj NewJob
	data := ctx.Data()
	if err := ctx.MustBind(&nj); err != nil {
		return ctx.JSON(data.SetError(err))
	}
	nj.FilePath = strings.Replace(nj.FilePath, `..`, ``, -1)
	nj.FilePath = strings.TrimLeft(nj.FilePath, `/`)
	nj.FilePath = strings.TrimLeft(nj.FilePath, `\`)
	dl, err := httpclient.CreateDownloader(nj.Url, nj.FilePath, nj.PartCount, srv.SavePath(), PipeGetList(nj.Pipes...)...)
	if err != nil {
		return ctx.JSON(data.SetError(err))
	}
	srv.dls = append(srv.dls, dl)
	return ctx.JSON(data)
}

func (srv *DServ) startTask(ctx echo.Context) error {
	srv.oplock.Lock()
	defer srv.oplock.Unlock()
	data := ctx.Data()
	for _, id := range ctx.FormValues(`id[]`) {
		ind := ctx.Atop(id).Int()
		if !(len(srv.dls) > ind) {
			return ctx.JSON(data.SetError(errors.New("error: id is out of jobs list")))
		}

		if errs := srv.dls[ind].StartAll(); len(errs) > 0 {
			return ctx.JSON(data.SetError(errors.New("error: can't start all part")))
		}
	}
	return ctx.JSON(data)
}

func (srv *DServ) stopTask(ctx echo.Context) error {
	srv.oplock.Lock()
	defer srv.oplock.Unlock()
	data := ctx.Data()
	for _, id := range ctx.FormValues(`id[]`) {
		ind := ctx.Atop(id).Int()
		if !(len(srv.dls) > ind) {
			return ctx.JSON(data.SetError(errors.New("error: id is out of jobs list")))
		}

		srv.dls[ind].StopAll()
	}
	return ctx.JSON(data)
}

func (srv *DServ) startAllTask(ctx echo.Context) error {
	srv.StartAllTask()
	return ctx.JSON(ctx.Data())
}

func (srv *DServ) StopAllTask() {
	srv.oplock.Lock()
	defer srv.oplock.Unlock()
	for _, e := range srv.dls {
		log.Println("info stopall result:", e.StopAll())
	}
}

func (srv *DServ) StartAllTask() {
	srv.oplock.Lock()
	defer srv.oplock.Unlock()
	for _, e := range srv.dls {
		log.Println("info start all result:", e.StartAll())
	}
}
func (srv *DServ) stopAllTask(ctx echo.Context) error {
	srv.StopAllTask()
	return ctx.JSON(ctx.Data())
}

func (srv *DServ) removeTask(ctx echo.Context) error {
	srv.oplock.Lock()
	defer srv.oplock.Unlock()
	data := ctx.Data()
	var decr int
	var ids []int
	for _, id := range ctx.FormValues(`id[]`) {
		ind := ctx.Atop(id).Int()
		var exists bool
		for _, i := range ids {
			if i == ind {
				exists = true
				break
			}
		}
		if exists {
			continue
		}
		ids = append(ids, ind)
	}
	sort.Ints(ids)
	for _, ind := range ids {
		ind -= decr
		if !(len(srv.dls) > ind) {
			return ctx.JSON(data.SetError(errors.New("error: id is out of jobs list")))
		}

		log.Printf("try stop segment download %v", srv.dls[ind].StopAll())
		srv.dls = append(srv.dls[:ind], srv.dls[ind+1:]...)
		decr++
	}
	return ctx.JSON(data)
}

func (srv *DServ) progressJson(ctx echo.Context) error {
	srv.oplock.Lock()
	defer srv.oplock.Unlock()
	return ctx.JSON(srv.Progress())
}

func (srv *DServ) Progress() []DJob {
	jbs := make([]DJob, 0, len(srv.dls))
	for ind, i := range srv.dls {
		var d int64        // 已下载字节
		var s int64        // 速度
		var total int64    // 总尺寸
		var progress int64 // 进度
		if i.ProgressGetter() != nil {
			d, total, progress, s = i.ProgressGetter()()
		} else {
			prs := i.GetProgress()
			for _, p := range prs {
				d = d + (p.Pos - p.From)
				s += p.Speed
			}
			total = i.Fi.Size
			progress = (d * 100 / total)
		}
		j := DJob{
			Id:         ind,
			FileName:   i.Fi.FileName,
			Size:       total,
			Progress:   progress,
			Downloaded: d,
			Speed:      s,
			State:      i.State().String(),
		}
		jbs = append(jbs, j)
	}
	return jbs
}

func (srv *DServ) ProgressSockJS(c sockjs.Session) error {
	exec := func(session sockjs.Session) error {
		for {
			command, err := session.Recv()
			if err != nil {
				log.Println(`Recv error: `, err.Error())
				return err
			}
			if len(command) == 0 {
				continue
			}
			message, _ := json.Marshal(srv.Progress())
			if err := c.Send(engine.Bytes2str(message)); err != nil {
				log.Println(`Push error: `, err.Error())
				return err
			}
		}
	}
	err := exec(c)
	if err != nil {
		log.Println(err)
	}
	return nil
}
