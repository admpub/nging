/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package handler

import (
	"github.com/webx-top/echo"
	ws "github.com/webx-top/echo/handler/websocket"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/route"

	"github.com/nging-plugins/dockermanager/application/library/utils"
	"github.com/nging-plugins/dockermanager/application/request"

	"github.com/nging-plugins/dockermanager/application/handler/docker"
	"github.com/nging-plugins/dockermanager/application/handler/docker/compose"
	"github.com/nging-plugins/dockermanager/application/handler/docker/container"
	"github.com/nging-plugins/dockermanager/application/handler/docker/image"

	//"github.com/nging-plugins/dockermanager/application/handler/kubernetes"
	"github.com/nging-plugins/dockermanager/application/handler/swarm"
)

func RegisterRoute(r *route.Collection) {
	r.Backend.RegisterToGroup(`/docker`, registerRoute, func(c echo.Context) error {
		c.SetFunc(`ShortenID`, utils.ShortenID)
		return nil
	})
}

func registerRoute(g echo.RouteRegister) {

	dockerG := g.Group(`/base`)
	dockerG.Get(`/index`, docker.Index)

	registryG := dockerG.Group(`/registry`)
	registryG.Route(`GET,POST`, `/login`, handler.WithRequest(docker.Login, request.Login{}, `POST`))

	containerG := dockerG.Group(`/container`)
	containerG.Route(`GET`, `/index`, container.Index)
	containerG.Route(`GET,POST`, `/add`, handler.WithRequest(container.Add, request.ContainerAdd{}, `POST`))
	containerG.Get(`/detail/:id`, container.Detail)
	containerG.Route(`GET,POST`, `/edit/:id`, handler.WithRequest(container.Edit, request.ContainerEdit{}, `POST`))
	// containerG.Route(`POST`, `/resize/:id`, handler.WithRequest(container.Resize, request.ContainerResize{}, `POST`))
	containerG.Route(`GET,POST`, `/delete/:id`, container.Delete)
	containerG.Route(`GET,POST`, `/prune`, container.Prune)
	containerG.Route(`GET,POST`, `/kill/:id`, container.Kill)
	containerG.Route(`GET,POST`, `/start/:id`, container.Start)
	containerG.Route(`GET,POST`, `/stop/:id`, container.Stop)
	containerG.Route(`GET,POST`, `/restart/:id`, container.Restart)
	containerG.Route(`GET,POST`, `/pause/:id`, container.Pause)
	containerG.Route(`GET,POST`, `/top/:id`, container.Top)
	// containerG.Route(`GET,POST`, `/statsOneShot/:id`, container.StatsOneShot)
	containerG.Route(`GET,POST`, `/statsPage/:id`, container.StatsPage)
	containerG.Route(`GET,POST`, `/fileImport/:id`, container.FileImport)
	containerG.Route(`GET,POST`, `/fileExport/:id`, container.FileExport)
	containerG.Route(`GET,POST`, `/file/:id`, container.File)
	containerG.Route(`GET`, `/export/:id`, container.Export)
	ws.New("/pty/:id", container.Pty).Wrapper(containerG)
	ws.New("/logs/:id", container.Logs).Wrapper(containerG)
	ws.New("/stats/:id", container.Stats).Wrapper(containerG)

	imageG := dockerG.Group(`/image`)
	imageG.Route(`GET`, `/index`, image.Index)
	imageG.Route(`GET,POST`, `/add`, handler.WithRequest(image.Add, request.ImageAdd{}, `POST`))
	// imageG.Route(`POST`, `/edit`, handler.WithRequest(image.Edit, request.ImageTag{}, `POST`)) // TODO
	imageG.Route(`GET`, `/detail/:id`, image.Detail)
	imageG.Route(`GET,POST`, `/delete/:id`, image.Delete)
	imageG.Route(`GET,POST`, `/pull`, image.Pull)
	imageG.Route(`GET,POST`, `/prune`, image.Prune)
	imageG.Route(`GET`, `/download/:id`, image.Download)
	imageG.Route(`GET,POST`, `/load`, image.Load)
	imageG.Route(`GET,POST`, `/import`, handler.WithRequest(image.Import, request.ImageImport{}, `POST`))
	// imageG.Route(`GET,POST`, `/build`, handler.WithRequest(image.Build, request.ImageBuild{}, `POST`)) // TODO

	networkG := dockerG.Group(`/network`)
	networkG.Route(`GET`, `/index`, docker.NetworkIndex)
	networkG.Route(`GET,POST`, `/add`, handler.WithRequest(docker.NetworkAdd, request.NetworkAdd{}, `POST`))
	// networkG.Route(`GET,POST`, `/connect`, handler.WithRequest(docker.NetworkConnect, request.NetworkConnect{}, `POST`)) // TODO
	// networkG.Route(`GET,POST`, `/disconnect`, handler.WithRequest(docker.NetworkDisconnect, request.NetworkDisconnect{}, `POST`)) // TODO
	networkG.Route(`GET`, `/detail/:id`, docker.NetworkDetail)
	networkG.Route(`GET,POST`, `/delete/:id`, docker.NetworkDelete)

	volumeG := dockerG.Group(`/volume`)
	volumeG.Route(`GET`, `/index`, docker.VolumeIndex)
	volumeG.Route(`GET,POST`, `/add`, handler.WithRequest(docker.VolumeAdd, request.VolumeAdd{}, `POST`))
	volumeG.Route(`GET`, `/detail/:id`, docker.VolumeDetail)
	volumeG.Route(`GET,POST`, `/delete/:id`, docker.VolumeDelete)
	volumeG.Route(`GET,POST`, `/prune`, docker.VolumePrune)

	composeG := dockerG.Group(`/compose`)
	composeG.Route(`GET`, `/index`, compose.Index)
	composeG.Route(`GET,POST`, `/add`, handler.WithRequest(compose.Add, request.ComposeAdd{}, `POST`))
	composeG.Route(`GET,POST`, `/edit/:id`, handler.WithRequest(compose.Edit, request.ComposeEdit{}, `POST`))
	composeG.Route(`GET`, `/detail/:id`, compose.Detail)
	composeG.Route(`GET`, `/listContainers/:id`, compose.ListContainers)
	composeG.Route(`GET,POST`, `/delete/:id`, compose.Delete)
	composeG.Route(`GET,POST`, `/stop/:id`, compose.Stop)
	composeG.Route(`GET,POST`, `/start/:id`, compose.Start)
	composeG.Route(`GET,POST`, `/scale/:id/:service`, compose.Scale)

	swarmG := g.Group(`/swarm`)
	swarmG.Route(`GET`, `/index`, swarm.Index)
	swarmG.Route(`GET,POST`, `/init`, swarm.SwarmInit)
	swarmG.Route(`GET,POST`, `/join`, swarm.SwarmJoin)
	swarmG.Route(`GET,POST`, `/leave`, swarm.SwarmLeave)
	// - config -
	swarmG.Route(`GET`, `/config/index`, swarm.ConfigIndex)
	swarmG.Route(`GET,POST`, `/config/add`, handler.WithRequest(swarm.ConfigAdd, request.SwarmConfigEdit{}, `POST`))
	swarmG.Route(`GET,POST`, `/config/edit/:id`, handler.WithRequest(swarm.ConfigEdit, request.SwarmConfigEdit{}, `POST`))
	swarmG.Route(`GET`, `/config/detail/:id`, swarm.ConfigDetail)
	swarmG.Route(`GET,POST`, `/config/delete/:id`, swarm.ConfigDelete)
	// - node -
	swarmG.Route(`GET`, `/node/index`, swarm.NodeIndex)
	swarmG.Route(`GET,POST`, `/node/add`, swarm.NodeAdd)
	swarmG.Route(`GET,POST`, `/node/edit/:id`, handler.WithRequest(swarm.NodeEdit, request.SwarmNodeEdit{}, `POST`))
	swarmG.Route(`GET`, `/node/detail/:id`, swarm.NodeDetail)
	swarmG.Route(`GET,POST`, `/node/delete/:id`, swarm.NodeDelete)
	// - secret -
	swarmG.Route(`GET`, `/secret/index`, swarm.SecretIndex)
	swarmG.Route(`GET,POST`, `/secret/add`, handler.WithRequest(swarm.SecretAdd, request.SwarmSecretEdit{}, `POST`))
	swarmG.Route(`GET,POST`, `/secret/edit/:id`, handler.WithRequest(swarm.SecretEdit, request.SwarmSecretEdit{}, `POST`))
	swarmG.Route(`GET`, `/secret/detail/:id`, swarm.SecretDetail)
	swarmG.Route(`GET,POST`, `/secret/delete/:id`, swarm.SecretDelete)
	// - service -
	swarmG.Route(`GET`, `/service/index`, swarm.ServiceIndex)
	swarmG.Route(`GET,POST`, `/service/add`, handler.WithRequest(swarm.ServiceAdd, request.SwarmServiceAdd{}, `POST`))
	swarmG.Route(`GET,POST`, `/service/edit/:id`, handler.WithRequest(swarm.ServiceEdit, request.SwarmServiceEdit{}, `POST`))
	swarmG.Route(`GET`, `/service/detail/:id`, swarm.ServiceDetail)
	swarmG.Route(`GET,POST`, `/service/delete/:id`, swarm.ServiceDelete)
	swarmG.Route(`POST`, `/service/rollback/:id`, swarm.ServiceRollback)
	ws.New("/service/logs/:id", swarm.ServiceLogs).Wrapper(swarmG)
	// - task -
	swarmG.Route(`GET`, `/task/index`, swarm.TaskIndex)
	swarmG.Route(`GET`, `/task/detail/:id`, swarm.TaskDetail)
	ws.New("/task/logs/:id", swarm.TaskLogs).Wrapper(swarmG)
	// - stack -
	stackG := swarmG.Group(`/stack`)
	stackG.Route(`GET`, `/index`, swarm.StackIndex)
	stackG.Route(`GET,POST`, `/add`, handler.WithRequest(swarm.StackAdd, request.StackAdd{}, `POST`))
	stackG.Route(`GET,POST`, `/edit/:id`, handler.WithRequest(swarm.StackEdit, request.StackEdit{}, `POST`))
	stackG.Route(`GET`, `/detail/:id`, swarm.StackDetail)
	stackG.Route(`GET`, `/listTasks/:id`, swarm.StackListTasks)
	stackG.Route(`GET`, `/listServices/:id`, swarm.StackListServices)
	stackG.Route(`GET,POST`, `/delete/:id`, swarm.StackDelete)
	stackG.Route(`GET,POST`, `/stop/:id`, swarm.StackStop)
	stackG.Route(`GET,POST`, `/start/:id`, swarm.StackStart)
}
