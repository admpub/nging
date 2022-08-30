package interfaces

import (
	"context"

	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager/domain/dnsdomain"
	"github.com/webx-top/echo"
)

type Updater interface {
	Name() string
	Description() string
	SignUpURL() string
	LineTypeURL() string
	Support() dnsdomain.Support
	Init(providerSettings echo.H, domains []*dnsdomain.Domain) error
	Update(ctx context.Context, recordType string, ip string) error
	ConfigItems() echo.KVList
}
