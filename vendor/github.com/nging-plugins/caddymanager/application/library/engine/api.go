package engine

import (
	"fmt"

	"github.com/admpub/resty/v2"
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

func (a *APIClient) Post(url string, data interface{}) error {
	resp, err := a.client.R().SetBody(data).Post(url)
	if err != nil {
		return err
	}
	if resp.IsError() {
		err = fmt.Errorf(`%s`, resp.Body())
	}
	return err
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
