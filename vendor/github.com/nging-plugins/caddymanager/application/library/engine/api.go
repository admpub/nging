package engine

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/admpub/resty/v2"
	"github.com/webx-top/echo"
)

// NewAPIClient("cert.pem", "key.pem")
func NewAPIClient(certPEMBlock, keyPEMBlock []byte) (*APIClient, error) {
	var rclient *resty.Client
	if len(certPEMBlock) > 0 && len(keyPEMBlock) > 0 {
		var err error
		rclient, err = newCertClient(certPEMBlock, keyPEMBlock)
		if err != nil {
			return nil, err
		}
	} else {
		rclient = defaultClient
	}
	return &APIClient{client: rclient}, nil
}

type APIClient struct {
	client *resty.Client
}

func (a *APIClient) SetClient(client *http.Client) *APIClient {
	a.client = newRestyClient(client)
	return a
}

// Post url=/v1.43/containers/{id}/exec
func (a *APIClient) Post(url string, data interface{}) error {
	var idResp IDResponse
	resp, err := a.client.R().SetBody(data).SetResult(&idResp).Post(url)
	if err != nil {
		return err
	}
	if resp.IsError() {
		err = fmt.Errorf(`%s`, resp.Body())
		return err
	}
	if len(idResp.ID) == 0 {
		return err
	}
	parts := strings.SplitN(url, `/containers/`, 2)
	if len(parts) != 2 {
		return err
	}
	url = parts[0] + `/exec/` + idResp.ID + `/start`
	resp, err = a.client.R().SetBody(RequestDockerExecStart{}).Post(url)
	if err != nil {
		return err
	}
	if resp.IsError() {
		err = fmt.Errorf(`%s`, resp.Body())
		return err
	}
	return err
}

func PostDocker(containerID string, data RequestDockerExec) error {
	exec := getContainerExec()
	if exec == nil {
		return echo.ErrNotImplemented
	}
	return exec(context.Background(), containerID, data.Cmd, data.Env, nil, nil)
}

// documention: https://docs.docker.com/engine/api/v1.43/#tag/Exec/operation/ContainerExec
// API: /containers/{id}/exec
type RequestDockerExec struct {
	AttachStdin  bool     `json:",omitempty"`
	AttachStdout bool     `json:",omitempty"`
	AttachStderr bool     `json:",omitempty"`
	DetachKeys   string   `json:",omitempty"` //"ctrl-p,ctrl-q",
	Tty          bool     `json:",omitempty"`
	Cmd          []string `json:",omitempty"` // ["date"],
	Env          []string `json:",omitempty"` // ["FOO=bar","BAZ=quux"]
}
type IDResponse struct {
	// The id of the newly created object.
	// Required: true
	ID string `json:"Id"`
}

type RequestDockerExecStart struct {
	Detach bool
	Tty    bool
}
