package interfaces

import (
	"github.com/admpub/nging/v3/application/library/ddnsmanager/domain"
	"github.com/webx-top/echo"
)

type Updater interface {
	Init(providerSettings echo.H, domains domain.Domains) error
	Update(recordType string) error
}
