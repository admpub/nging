package engine

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/acme/autocert"
)

const DefaultMaxRequestBodySize = 32 << 20 // 32 MB

// Config defines engine configuration.
type Config struct {
	Address            string       // TCP address to listen on.
	Listener           net.Listener // Custom `net.Listener`. If set, server accepts connections on it.
	ReusePort          bool
	TLSAuto            bool
	TLSHosts           []string
	TLSEmail           string
	TLSCacheDir        string
	TLSConfig          *tls.Config
	TLSCertFile        string        // TLS certificate file path.
	TLSKeyFile         string        // TLS key file path.
	DisableHTTP2       bool          // Disables HTTP/2.
	ReadTimeout        time.Duration // Maximum duration before timing out read of the request.
	WriteTimeout       time.Duration // Maximum duration before timing out write of the response.
	MaxConnsPerIP      int
	MaxRequestsPerConn int
	MaxRequestBodySize int
}

//usage:
//c.InitTLSConfig(`cert.pem`,`key.pem`).AddTLSCert(`cert2.pem`,`key2.pem`).SupportAutoTLS(nil,`webx.top`,`coscms.com`)
//or c.AddTLSCert(`cert.pem`,`key.pem`).SupportAutoTLS(nil,`webx.top`,`coscms.com`)
//or c.SupportAutoTLS(nil,`webx.top`,`coscms.com`)

func (c *Config) InitTLSConfig(certAndKey ...string) *Config {
	switch len(certAndKey) {
	case 2:
		c.TLSKeyFile = certAndKey[1]
		fallthrough
	case 1:
		c.TLSCertFile = certAndKey[0]
	}
	c.TLSConfig = new(tls.Config)
	if len(c.TLSCertFile) > 0 && len(c.TLSKeyFile) > 0 {
		cert, err := tls.LoadX509KeyPair(c.TLSCertFile, c.TLSKeyFile)
		if err != nil {
			panic(err)
		}
		c.TLSConfig.Certificates = append(c.TLSConfig.Certificates, cert)
	}
	if !c.DisableHTTP2 {
		c.TLSConfig.NextProtos = append(c.TLSConfig.NextProtos, "h2")
	}
	return c
}

func (c *Config) AddTLSCert(certFile, keyFile string) *Config {
	if c.TLSConfig == nil {
		c.InitTLSConfig()
	}
	if len(certFile) > 0 && len(keyFile) > 0 {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			panic(err)
		}
		c.TLSConfig.Certificates = append(c.TLSConfig.Certificates, cert)
	}
	return c
}

func (c *Config) NewAutoTLSManager(hosts ...string) *autocert.Manager {
	autoTLSManager := &autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Email:  c.TLSEmail,
	}
	if len(hosts) > 0 {
		c.TLSHosts = append(c.TLSHosts, hosts...)
	}
	autoTLSManager.HostPolicy = autocert.HostWhitelist(c.TLSHosts...) // Added security
	if len(c.TLSCacheDir) == 0 {
		home, err := homedir.Dir()
		if err != nil {
			panic(err)
		}
		c.TLSCacheDir = filepath.Join(home, ".webx.top", "cache", "autocert")
	}
	if _, err := os.Stat(c.TLSCacheDir); os.IsNotExist(err) {
		err = os.MkdirAll(c.TLSCacheDir, 0666)
		if err != nil {
			panic(err)
		}
	}
	autoTLSManager.Cache = autocert.DirCache(c.TLSCacheDir)
	return autoTLSManager
}

func (c *Config) SupportAutoTLS(autoTLSManager *autocert.Manager, hosts ...string) *Config {
	if c.TLSConfig == nil {
		c.InitTLSConfig()
	}
	if autoTLSManager == nil {
		autoTLSManager = c.NewAutoTLSManager(hosts...)
	}
	if c.Listener == nil && AddressPort(c.Address) != 80 {
		go func() {
			log.Println(`Starting serve: ACME "http-01" challenge responses.`)
			err := http.ListenAndServe(":http", autoTLSManager.HTTPHandler(nil))
			if err != nil {
				log.Println(err)
			}
		}()
	}
	//c.TLSConfig.GetCertificate = autoTLSManager.GetCertificate
	c.TLSConfig.BuildNameToCertificate()
	c.TLSConfig.GetCertificate = func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		if cert, ok := c.TLSConfig.NameToCertificate[clientHello.ServerName]; ok {
			// Use provided certificate
			return cert, nil
		}
		if c.TLSAuto {
			return autoTLSManager.GetCertificate(clientHello)
		}
		return nil, nil // No certificate
	}
	return c
}

func (c *Config) InitTLSListener(before ...func() error) error {
	if c.TLSConfig == nil {
		c.InitTLSConfig()
		if c.TLSAuto {
			c.SupportAutoTLS(nil)
		}
	}
	if len(before) > 0 && before[0] != nil {
		if err := before[0](); err != nil {
			return err
		}
	}
	ln, err := NewListener(c.Address, c.ReusePort)
	if err != nil {
		return err
	}
	c.Listener = tls.NewListener(ln, c.TLSConfig)
	return nil
}

func (c *Config) InitListener(before ...func() error) error {
	if c.TLSAuto || (len(c.TLSCertFile) > 0 && len(c.TLSKeyFile) > 0) {
		return c.InitTLSListener(before...)
	}
	if len(before) > 0 && before[0] != nil {
		if err := before[0](); err != nil {
			return err
		}
	}
	ln, err := NewListener(c.Address, c.ReusePort)
	if err != nil {
		return err
	}
	c.Listener = ln
	return nil
}

func (c *Config) Print(engine string) {
	var s string
	if c.TLSConfig != nil {
		s = `s`
	}
	log.Printf("%s ⇛ http%s server started on %s\n", engine, s, c.Listener.Addr())
}

func (c *Config) SetListener(ln net.Listener) *Config {
	c.Listener = ln
	return c
}
