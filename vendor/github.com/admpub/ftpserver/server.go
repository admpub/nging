package server

import (
	"bufio"
	"crypto/tls"
	"net"
	"strconv"
)

func Version() string {
	return "0.2.1104"
}

// serverOpts contains parameters for server.NewServer()
type ServerOpts struct {
	// The factory that will be used to create a new FTPDriver instance for
	// each client connection. This is a mandatory option.
	Factory DriverFactory `json:"-"`

	Auth Auth `json:"-"`

	// Server Name, Default is Go Ftp Server
	Name string `json:"name"`

	// The hostname that the FTP server should listen on. Optional, defaults to
	// "::", which means all hostnames on ipv4 and ipv6.
	Hostname string `json:"hostName"`

	// Public IP of the server
	PublicIp string `json:"publicIP"`

	// Passive ports
	PassivePorts string `json:"passivePorts"`

	// The port that the FTP should listen on. Optional, defaults to 3000. In
	// a production environment you will probably want to change this to 21.
	Port int `json:"port"`

	// use tls, default is false
	TLS bool `json:"tls"`

	// if tls used, cert file is required
	CertFile string `json:"certFile"`

	// if tls used, key file is required
	KeyFile string `json:"keyFile"`

	// If ture TLS is used in RFC4217 mode
	ExplicitFTPS bool `json:"explicitFTPS"`

	WelcomeMessage string `json:"welcomeMessage"`
}

// Server is the root of your FTP application. You should instantiate one
// of these and call ListenAndServe() to start accepting client connections.
//
// Always use the NewServer() method to create a new Server.
type Server struct {
	*ServerOpts
	listenTo  string
	logger    *Logger
	listener  net.Listener
	tlsConfig *tls.Config
}

// serverOptsWithDefaults copies an ServerOpts struct into a new struct,
// then adds any default values that are missing and returns the new data.
func serverOptsWithDefaults(opts *ServerOpts) *ServerOpts {
	var newOpts ServerOpts
	if opts == nil {
		opts = &ServerOpts{}
	}
	if opts.Hostname == "" {
		newOpts.Hostname = "::"
	} else {
		newOpts.Hostname = opts.Hostname
	}
	if opts.Port == 0 {
		newOpts.Port = 3000
	} else {
		newOpts.Port = opts.Port
	}
	newOpts.Factory = opts.Factory
	if opts.Name == "" {
		newOpts.Name = "Go FTP Server"
	} else {
		newOpts.Name = opts.Name
	}

	if opts.WelcomeMessage == "" {
		newOpts.WelcomeMessage = defaultWelcomeMessage
	} else {
		newOpts.WelcomeMessage = opts.WelcomeMessage
	}

	if opts.Auth != nil {
		newOpts.Auth = opts.Auth
	}

	newOpts.TLS = opts.TLS
	newOpts.KeyFile = opts.KeyFile
	newOpts.CertFile = opts.CertFile
	newOpts.ExplicitFTPS = opts.ExplicitFTPS

	newOpts.PublicIp = opts.PublicIp
	newOpts.PassivePorts = opts.PassivePorts

	return &newOpts
}

// NewServer initialises a new FTP server. Configuration options are provided
// via an instance of ServerOpts. Calling this function in your code will
// probably look something like this:
//
//     factory := &MyDriverFactory{}
//     server  := server.NewServer(&server.ServerOpts{ Factory: factory })
//
// or:
//
//     factory := &MyDriverFactory{}
//     opts    := &server.ServerOpts{
//       Factory: factory,
//       Port: 2000,
//       Hostname: "127.0.0.1",
//     }
//     server  := server.NewServer(opts)
//
func NewServer(opts *ServerOpts) *Server {
	opts = serverOptsWithDefaults(opts)
	s := new(Server)
	s.ServerOpts = opts
	s.listenTo = net.JoinHostPort(opts.Hostname, strconv.Itoa(opts.Port))
	s.logger = newLogger("")
	return s
}

// NewConn constructs a new object that will handle the FTP protocol over
// an active net.TCPConn. The TCP connection should already be open before
// it is handed to this functions. driver is an instance of FTPDriver that
// will handle all auth and persistence details.
func (server *Server) newConn(tcpConn net.Conn, driver Driver) *Conn {
	c := new(Conn)
	c.namePrefix = "/"
	c.conn = tcpConn
	c.controlReader = bufio.NewReader(tcpConn)
	c.controlWriter = bufio.NewWriter(tcpConn)
	c.driver = driver
	c.auth = server.Auth
	c.server = server
	c.sessionID = newSessionID()
	c.logger = newLogger(c.sessionID)
	c.tlsConfig = server.tlsConfig
	driver.Init(c)
	return c
}

func simpleTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	config := &tls.Config{}
	if config.NextProtos == nil {
		config.NextProtos = []string{"ftp"}
	}

	var err error
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// ListenAndServe asks a new Server to begin accepting client connections. It
// accepts no arguments - all configuration is provided via the NewServer
// function.
//
// If the server fails to start for any reason, an error will be returned. Common
// errors are trying to bind to a privileged port or something else is already
// listening on the same port.
//
func (server *Server) ListenAndServe() error {
	var listener net.Listener
	var err error

	if server.ServerOpts.TLS {
		server.tlsConfig, err = simpleTLSConfig(server.CertFile, server.KeyFile)
		if err != nil {
			return err
		}

		if server.ServerOpts.ExplicitFTPS {
			listener, err = net.Listen("tcp", server.listenTo)
		} else {
			listener, err = tls.Listen("tcp", server.listenTo, server.tlsConfig)
		}
	} else {
		listener, err = net.Listen("tcp", server.listenTo)
	}
	if err != nil {
		return err
	}

	server.logger.Printf("%s listening on %d", server.Name, server.Port)

	server.listener = listener
	for {
		tcpConn, err := server.listener.Accept()
		if err != nil {
			server.logger.Printf("listening error: %v", err)
			break
		}
		driver, err := server.Factory.NewDriver()
		if err != nil {
			server.logger.Printf("Error creating driver, aborting client connection: %v", err)
			tcpConn.Close()
		} else {
			ftpConn := server.newConn(tcpConn, driver)
			go ftpConn.Serve()
		}
	}
	return nil
}

// Gracefully stops a server. Already connected clients will retain their connections
func (server *Server) Shutdown() error {
	if server.listener != nil {
		return server.listener.Close()
	}
	// server wasnt even started
	return nil
}
