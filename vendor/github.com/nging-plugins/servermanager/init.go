package servermanager

import (
	"github.com/nging-plugins/servermanager/pkg/handler"
	_ "github.com/nging-plugins/servermanager/pkg/library/cmder"
	_ "github.com/nging-plugins/servermanager/pkg/library/setup"
)

var LeftNavigate = handler.LeftNavigate

func init() {

}
