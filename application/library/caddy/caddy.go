package caddy

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"

	"github.com/mholt/caddy"
	_ "github.com/mholt/caddy/caddyhttp"
	"github.com/mholt/caddy/caddytls"
	"github.com/xenolf/lego/acme"
)

var (
	DefaultCAUrl  = `https://acme-v01.api.letsencrypt.org/directory`
	DefaultConfig = &Config{
		CAUrl:      DefaultCAUrl,
		CATimeout:  int64(acme.HTTPClient.Timeout.Seconds()),
		ServerType: `http`,
		CPU:        `100%`,
	}
)

func TrapSignals() {
	caddy.TrapSignals()
}

func Fixed(c *Config) {
	if c.CAUrl == `` {
		c.CAUrl = DefaultConfig.CAUrl
	}
	if c.CATimeout == 0 {
		c.CATimeout = DefaultConfig.CATimeout
	}
	if c.ServerType == `` {
		c.ServerType = DefaultConfig.ServerType
	}
	if c.CPU == `` {
		c.CPU = DefaultConfig.CPU
	}
}

type Config struct {
	Agreed     bool   `json:"agreed"`     //Agree to the CA's Subscriber Agreement
	CAUrl      string `json:"caURL"`      //URL to certificate authority's ACME server directory
	Caddyfile  string `json:"caddyFile"`  //Caddyfile to load (default caddy.DefaultConfigFile)
	CPU        string `json:"cpu"`        //CPU cap
	CAEmail    string `json:"caEmail"`    //Default ACME CA account email address
	CATimeout  int64  `json:"caTimeout"`  //Default ACME CA HTTP timeout
	LogFile    string `json:"logFile"`    //Process log file
	PidFile    string `json:"pidFile"`    //Path to write pid file
	Quiet      bool   `json:"quiet"`      //Quiet mode (no initialization output)
	Revoke     string `json:"revoke"`     //Hostname for which to revoke the certificate
	ServerType string `json:"serverType"` //Type of server to run

	//---
	Plugins bool `json:"plugins"` //List installed plugins
	Version bool `json:"version"` //Show version

	//---
	AppName    string `json:"appName"`
	AppVersion string `json:"appVersion"`

	instance *caddy.Instance
}

func (c *Config) Start() {
	// Get Caddyfile input
	caddyfile, err := caddy.LoadCaddyfile(c.ServerType)
	if err != nil {
		mustLogFatalf(err.Error())
	}

	// Start your engines
	c.instance, err = caddy.Start(caddyfile)
	if err != nil {
		mustLogFatalf(err.Error())
	}

	// Twiddle your thumbs
	c.instance.Wait()
}

func (c *Config) Restart() {
	if c.instance == nil {
		return
	}
	// Get Caddyfile input
	caddyfile, err := caddy.LoadCaddyfile(c.ServerType)
	if err != nil {
		mustLogFatalf(err.Error())
	}
	c.instance, err = c.instance.Restart(caddyfile)
	if err != nil {
		mustLogFatalf(err.Error())
	}

	c.instance.Wait()
}

func (c *Config) Stop() {
	if c.instance == nil {
		return
	}
	err := c.instance.Stop()
	if err != nil {
		mustLogFatalf(err.Error())
	}
}

func (c *Config) Init() *Config {
	caddytls.Agreed = c.Agreed
	caddytls.DefaultCAUrl = c.CAUrl
	caddytls.DefaultEmail = c.CAEmail
	acme.HTTPClient.Timeout = time.Duration(c.CATimeout)
	caddy.PidFile = c.PidFile
	caddy.Quiet = c.Quiet
	caddy.RegisterCaddyfileLoader("flag", caddy.LoaderFunc(c.confLoader))
	caddy.SetDefaultCaddyfileLoader("default", caddy.LoaderFunc(c.defaultLoader))

	acme.UserAgent = c.AppName + "/" + c.AppVersion

	// Set up process log before anything bad happens
	switch c.LogFile {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "stderr":
		log.SetOutput(os.Stderr)
	case "":
		log.SetOutput(ioutil.Discard)
	default:
		log.SetOutput(&lumberjack.Logger{
			Filename:   c.LogFile,
			MaxSize:    100,
			MaxAge:     14,
			MaxBackups: 10,
		})
	}

	// Check for one-time actions
	if c.Revoke != "" {
		err := caddytls.Revoke(c.Revoke)
		if err != nil {
			mustLogFatalf(err.Error())
		}
		fmt.Printf("Revoked certificate for %s\n", c.Revoke)
		os.Exit(0)
	}
	if c.Version {
		fmt.Printf("%s %s\n", c.AppName, c.AppVersion)
		os.Exit(0)
	}
	if c.Plugins {
		fmt.Println(caddy.DescribePlugins())
		os.Exit(0)
	}

	moveStorage() // TODO: This is temporary for the 0.9 release, or until most users upgrade to 0.9+

	// Set CPU cap
	err := setCPU(c.CPU)
	if err != nil {
		mustLogFatalf(err.Error())
	}
	return c
}

// confLoader loads the Caddyfile using the -conf flag.
func (c *Config) confLoader(serverType string) (caddy.Input, error) {
	contents, err := ioutil.ReadFile(c.Caddyfile)
	if err != nil {
		return nil, err
	}
	return caddy.CaddyfileInput{
		Contents:       contents,
		Filepath:       c.Caddyfile,
		ServerTypeName: serverType,
	}, nil
}

// defaultLoader loads the Caddyfile from the current working directory.
func (c *Config) defaultLoader(serverType string) (caddy.Input, error) {
	contents, err := ioutil.ReadFile(caddy.DefaultConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return caddy.CaddyfileInput{
		Contents:       contents,
		Filepath:       caddy.DefaultConfigFile,
		ServerTypeName: serverType,
	}, nil
}

// mustLogFatalf wraps log.Fatalf() in a way that ensures the
// output is always printed to stderr so the user can see it
// if the user is still there, even if the process log was not
// enabled. If this process is an upgrade, however, and the user
// might not be there anymore, this just logs to the process
// log and exits.
func mustLogFatalf(format string, args ...interface{}) {
	if !caddy.IsUpgrade() {
		log.SetOutput(os.Stderr)
	}
	log.Fatalf(format, args...)
}

// moveStorage moves the old certificate storage location by
// renaming the "letsencrypt" folder to the hostname of the
// CA URL. This is TEMPORARY until most users have upgraded to 0.9+.
func moveStorage() {
	oldPath := filepath.Join(caddy.AssetsPath(), "letsencrypt")
	_, err := os.Stat(oldPath)
	if os.IsNotExist(err) {
		return
	}
	// Just use a default config to get default (file) storage
	fileStorage, err := new(caddytls.Config).StorageFor(caddytls.DefaultCAUrl)
	if err != nil {
		mustLogFatalf("[ERROR] Unable to get new path for certificate storage: %v", err)
	}
	newPath := fileStorage.(*caddytls.FileStorage).Path
	err = os.MkdirAll(string(newPath), 0700)
	if err != nil {
		mustLogFatalf("[ERROR] Unable to make new certificate storage path: %v\n\nPlease follow instructions at:\nhttps://github.com/mholt/caddy/issues/902#issuecomment-228876011", err)
	}
	err = os.Rename(oldPath, string(newPath))
	if err != nil {
		mustLogFatalf("[ERROR] Unable to migrate certificate storage: %v\n\nPlease follow instructions at:\nhttps://github.com/mholt/caddy/issues/902#issuecomment-228876011", err)
	}
	// convert mixed case folder and file names to lowercase
	var done bool // walking is recursive and preloads the file names, so we must restart walk after a change until no changes
	for !done {
		done = true
		filepath.Walk(string(newPath), func(path string, info os.FileInfo, err error) error {
			// must be careful to only lowercase the base of the path, not the whole thing!!
			base := filepath.Base(path)
			if lowerBase := strings.ToLower(base); base != lowerBase {
				lowerPath := filepath.Join(filepath.Dir(path), lowerBase)
				err = os.Rename(path, lowerPath)
				if err != nil {
					mustLogFatalf("[ERROR] Unable to lower-case: %v\n\nPlease follow instructions at:\nhttps://github.com/mholt/caddy/issues/902#issuecomment-228876011", err)
				}
				// terminate traversal and restart since Walk needs the updated file list with new file names
				done = false
				return errors.New("start over")
			}
			return nil
		})
	}
}

// setCPU parses string cpu and sets GOMAXPROCS
// according to its value. It accepts either
// a number (e.g. 3) or a percent (e.g. 50%).
func setCPU(cpu string) error {
	var numCPU int

	availCPU := runtime.NumCPU()

	if strings.HasSuffix(cpu, "%") {
		// Percent
		var percent float32
		pctStr := cpu[:len(cpu)-1]
		pctInt, err := strconv.Atoi(pctStr)
		if err != nil || pctInt < 1 || pctInt > 100 {
			return errors.New("invalid CPU value: percentage must be between 1-100")
		}
		percent = float32(pctInt) / 100
		numCPU = int(float32(availCPU) * percent)
	} else {
		// Number
		num, err := strconv.Atoi(cpu)
		if err != nil || num < 1 {
			return errors.New("invalid CPU value: provide a number or percent greater than 0")
		}
		numCPU = num
	}

	if numCPU > availCPU {
		numCPU = availCPU
	}

	runtime.GOMAXPROCS(numCPU)
	return nil
}
