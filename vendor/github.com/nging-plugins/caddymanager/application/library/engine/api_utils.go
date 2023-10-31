package engine

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"

	"github.com/admpub/resty/v2"
	"github.com/webx-top/restyclient"
)

var defaultClient = resty.New().SetRetryCount(restyclient.DefaultMaxRetryCount).
	SetTimeout(restyclient.DefaultTimeout).
	SetRedirectPolicy(restyclient.DefaultRedirectPolicy)

func init() {
	restyclient.InitRestyHook(defaultClient)
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
