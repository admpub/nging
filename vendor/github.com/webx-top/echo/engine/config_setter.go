package engine

import (
	"crypto/tls"
	"net"
	"time"
)

type ConfigSetter func(*Config)

// Address TCP address to listen on.
func Address(v string) ConfigSetter {
	return func(c *Config) {
		c.Address = v
	}
}

//Listener Custom `net.Listener`. If set, server accepts connections on it.
func Listener(v net.Listener) ConfigSetter {
	return func(c *Config) {
		c.Listener = v
	}
}

func ReusePort(v bool) ConfigSetter {
	return func(c *Config) {
		c.ReusePort = v
	}
}

func TLSAuto(v bool) ConfigSetter {
	return func(c *Config) {
		c.TLSAuto = v
	}
}

func TLSHosts(v []string) ConfigSetter {
	return func(c *Config) {
		c.TLSHosts = v
	}
}

func TLSEmail(v string) ConfigSetter {
	return func(c *Config) {
		c.TLSEmail = v
	}
}

func TLSCacheDir(v string) ConfigSetter {
	return func(c *Config) {
		c.TLSCacheDir = v
	}
}

func TLSConfig(v *tls.Config) ConfigSetter {
	return func(c *Config) {
		c.TLSConfig = v
	}
}

// TLSCertFile TLS certificate file path.
func TLSCertFile(v string) ConfigSetter {
	return func(c *Config) {
		c.TLSCertFile = v
	}
}

// TLSKeyFile TLS key file path.
func TLSKeyFile(v string) ConfigSetter {
	return func(c *Config) {
		c.TLSKeyFile = v
	}
}

// DisableHTTP2 Disables HTTP/2.
func DisableHTTP2(v bool) ConfigSetter {
	return func(c *Config) {
		c.DisableHTTP2 = v
	}
}

// ReadTimeout Maximum duration before timing out read of the request.
func ReadTimeout(v time.Duration) ConfigSetter {
	return func(c *Config) {
		c.ReadTimeout = v
	}
}

// WriteTimeout Maximum duration before timing out write of the response.
func WriteTimeout(v time.Duration) ConfigSetter {
	return func(c *Config) {
		c.WriteTimeout = v
	}
}

func MaxConnsPerIP(v int) ConfigSetter {
	return func(c *Config) {
		c.MaxConnsPerIP = v
	}
}

func MaxRequestsPerConn(v int) ConfigSetter {
	return func(c *Config) {
		c.MaxRequestsPerConn = v
	}
}

func MaxRequestBodySize(v int) ConfigSetter {
	return func(c *Config) {
		c.MaxRequestBodySize = v
	}
}
