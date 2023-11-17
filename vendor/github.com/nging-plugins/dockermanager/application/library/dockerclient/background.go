package dockerclient

import (
	"context"
	"io"
	"log"

	"github.com/admpub/nging/v5/application/library/background"
	"github.com/admpub/nging/v5/application/library/notice"
	"github.com/nging-plugins/dockermanager/application/library/utils"
	"github.com/webx-top/echo"
)

func StartBackgroundRun(ctx echo.Context, username string, actionIdent string, bgKey string, runner func(context.Context) (io.ReadCloser, error)) error {
	bg := background.New(context.Background(), nil)
	group, err := background.Register(ctx, actionIdent, bgKey, bg)
	if err != nil {
		return err
	}
	go func() {
		defer group.Cancel(bgKey)
		noticer := notice.NewP(ctx, actionIdent, username, bg.Context()).Add(100).AutoComplete(true)
		reader, err := runner(bg.Context())
		if err != nil {
			noticer.Send(err.Error(), notice.StateFailure)
			return
		}
		SyncReaderToNotice(noticer, reader)
	}()
	return nil
}

func SyncReaderToNotice(noticer notice.NProgressor, reader io.ReadCloser) {
	var lastMessage string
	defer reader.Close()
	buf := utils.NewResultWriter(func(r *utils.Result) error {
		if r.Completed() {
			noticer.Complete()
		}
		lastMessage = r.Status
		return noticer.Send(r.Status, notice.StateSuccess)
	})
	w := io.MultiWriter(buf, log.Writer())
	_, err := io.Copy(w, reader)
	if err != nil {
		noticer.Complete().Send(err.Error(), notice.StateFailure)
		return
	}
	buf.Flush()
	noticer.Done(100).Complete().Send(lastMessage, notice.StateSuccess)
}
