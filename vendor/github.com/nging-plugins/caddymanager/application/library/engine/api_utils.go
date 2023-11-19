package engine

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/admpub/resty/v2"
	"github.com/webx-top/echo"
	"github.com/webx-top/restyclient"
)

var defaultClient = newRestyClient(nil)

func newRestyClient(client *http.Client) *resty.Client {
	var c *resty.Client
	if client == nil {
		c = resty.New()
	} else {
		c = resty.NewWithClient(client)
	}
	c.SetRetryCount(restyclient.DefaultMaxRetryCount).
		SetTimeout(restyclient.DefaultTimeout).
		SetRedirectPolicy(restyclient.DefaultRedirectPolicy)
	restyclient.InitRestyHook(c)
	return c
}

// /var/run/docker.sock
func NewSocketClient(sockAddr string) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", sockAddr)
			},
		},
	}
}

func ParseSocketAddr(sockAddr string) string {
	sockAddr = strings.TrimPrefix(sockAddr, `unix:`)
	if !strings.HasPrefix(sockAddr, `/`) {
		return `/` + sockAddr
	}
	maxIndex := len(sockAddr) - 1
	for index, char := range sockAddr {
		if char == '/' && index+1 <= maxIndex && sockAddr[index+1] != '/' {
			return sockAddr[index:]
		}
	}
	return sockAddr
}

func newCertClient(certPEMBlock, keyPEMBlock []byte) (rclient *resty.Client, err error) {
	// create certificate
	cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(certPEMBlock)

	// Create a HTTPS client and supply the created CA pool and certificate
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
	}
	rclient = resty.NewWithClient(client)
	rclient.SetRetryCount(restyclient.DefaultMaxRetryCount).
		SetTimeout(restyclient.DefaultTimeout).
		SetRedirectPolicy(restyclient.DefaultRedirectPolicy)
	restyclient.InitRestyHook(rclient)
	return
}

type ContainerExec func(ctx context.Context, containerID string, cmd []string, env []string, outWriter io.Writer, errWriter io.Writer) error

func getContainerExec() ContainerExec {
	exec, _ := echo.Get(`DockerContainerExec`).(func(context.Context, string, []string, []string, io.Writer, io.Writer) error)
	return exec
}
