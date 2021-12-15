# https/certificates
A simple library to generate server certs and keys for HTTPS support directly within your Go program.

The code is modified from http://golang.org/src/crypto/tls/generate_cert.go.

Use this library for testing purposes only, e.g. to experiment with the built-in Go HTTPS server. Do NOT use in production!

PR for this library is https://github.com/kabukky/httpscerts and https://github.com/gerald1248/httpscerts.

# Usage

```go
package main
    
import (
    "fmt"
    "github.com/admpub/httpscerts"
    "log"
    "net/http"
)
    
func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there!")
}
    
func main() {
    // Check if the cert files are available.
    err := httpscerts.Check("cert.pem", "key.pem")
    // If they are not available, generate new ones.
    if err != nil {
        err = httpscerts.Generate("cert.pem", "key.pem", "127.0.0.1:8081")
        if err != nil {
            log.Fatal("Error: Couldn't create https certs.")
        }
    }
    http.HandleFunc("/", handler)
    http.ListenAndServeTLS(":8081", "cert.pem", "key.pem", nil)
}
```

# Alternative usage without disk access

The method `httpscerts.GenerateArrays()` has been added to enable use cases where writing to disk is not desirable. If the initial check fails, a `tls.Certificate` is populated and passed to a `http.Server` instance.

```go
package main

import (
	"crypto/tls"
	"fmt"
	"github.com/admpub/httpscerts"
	"log"
	"net/http"
	"time"
)

type testHandler struct {
}

func (t *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there!")
}

func main() {
	// Check if the cert files are available.
	certFile := "cert.pem"
	keyFile := "key.pem"
	err := httpscerts.Check(certFile, keyFile)

	var handler = &testHandler{}

	// If they are not available, generate new ones.
	if err != nil {
		cert, key, err := httpscerts.GenerateArrays("127.0.0.1:8081")
		if err != nil {
			log.Fatal("Error: Couldn't create https certs.")
		}

		keyPair, err := tls.X509KeyPair(cert, key)
		if err != nil {
			log.Fatal("Error: Couldn't create key pair")
		}

		var certificates []tls.Certificate
		certificates = append(certificates, keyPair)

		cfg := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
			Certificates:             certificates,
		}

		s := &http.Server{
			Addr: ":8081",
			Handler:        handler,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
			TLSConfig:      cfg,
		}
		log.Fatal(s.ListenAndServeTLS("", ""))
	}

	log.Fatal(http.ListenAndServeTLS(":8081", certFile, keyFile, handler))
}

```
