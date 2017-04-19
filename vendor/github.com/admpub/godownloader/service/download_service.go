package service

import (
	"errors"
	"log"
	"sync"

	"encoding/json"

	"github.com/admpub/godownloader/httpclient"
	"github.com/admpub/sockjs-go/sockjs"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
	sockjsHandler "github.com/webx-top/echo/handler/sockjs"
)

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
}

type DServ struct {
	dls    []*httpclient.Downloader
	oplock sync.Mutex
	tmpl   string
}

func (srv *DServ) SetTmpl(tmpl string) *DServ {
	srv.tmpl = tmpl
	return srv
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
	log.Println(ss)
	for _, r := range ss.Ds {
		dl, err := httpclient.RestoreDownloader(r.FI.Url, r.FI.FileName, r.Dp)
		if err != nil {
			return err
		}
		srv.dls = append(srv.dls, dl)
	}
	return nil
}

func (srv *DServ) index(ctx echo.Context) error {
	return ctx.Render(srv.tmpl, nil)
}

func (srv *DServ) addTask(ctx echo.Context) error {
	srv.oplock.Lock()
	defer srv.oplock.Unlock()
	var nj NewJob
	data := ctx.NewData()
	if err := ctx.MustBind(&nj); err != nil {
		return ctx.JSON(data.SetError(err))
	}
	dl, err := httpclient.CreateDownloader(nj.Url, nj.FilePath, nj.PartCount)
	if err != nil {
		return ctx.JSON(data.SetError(err))
	}
	srv.dls = append(srv.dls, dl)
	return ctx.JSON(data)
}

func (srv *DServ) startTask(ctx echo.Context) error {
	srv.oplock.Lock()
	defer srv.oplock.Unlock()
	data := ctx.NewData()
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
	data := ctx.NewData()
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
	data := ctx.NewData()
	srv.StartAllTask()
	return ctx.JSON(data)
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
	data := ctx.NewData()
	srv.StopAllTask()
	return ctx.JSON(data)
}

func (srv *DServ) removeTask(ctx echo.Context) error {
	srv.oplock.Lock()
	defer srv.oplock.Unlock()
	data := ctx.NewData()
	for _, id := range ctx.FormValues(`id[]`) {
		ind := ctx.Atop(id).Int()
		if !(len(srv.dls) > ind) {
			return ctx.JSON(data.SetError(errors.New("error: id is out of jobs list")))
		}

		log.Printf("try stop segment download %v", srv.dls[ind].StopAll())
		srv.dls = append(srv.dls[:ind], srv.dls[ind+1:]...)
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
		prs := i.GetProgress()
		var d int64
		var s int64
		for _, p := range prs {
			d = d + (p.Pos - p.From)
			s += p.Speed
		}
		j := DJob{
			Id:         ind,
			FileName:   i.Fi.FileName,
			Size:       i.Fi.Size,
			Progress:   (d * 100 / i.Fi.Size),
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
