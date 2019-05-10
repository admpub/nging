package imageproxy

import (
	"io/ioutil"
	"log"
	"net/url"
	"strings"

	"github.com/admpub/httpcache"
	"github.com/admpub/httpcache/diskcache"
	"github.com/peterbourgon/diskv"
)

var Default *Proxy

type Config struct {
	ProxyURL     string
	CacheDir     string
	CacheSize    uint64
	Whitelist    string
	Referrers    string
	SignatureKey string
	BaseURL      string
	ScaleUp      bool
	ResrcVKey    string
	ValidToken   bool
}

func (c *Config) URL(imageURL string, option string) (r string) {

	if len(c.SignatureKey) <= 0 {
		r = c.ProxyURL + option + "/" + url.QueryEscape(imageURL)
		return
	}
	if c.ValidToken {
		text := option + "/" + imageURL
		token := GenerateSignature([]byte(c.SignatureKey), []byte(text))
		r = c.ProxyURL + token + "/" + option + "/" + url.QueryEscape(imageURL)
		return
	}

	token := GenerateSignature([]byte(c.SignatureKey), []byte(imageURL))
	r = c.ProxyURL + option + "," + optSignaturePrefix + token + "/" + url.QueryEscape(imageURL)
	return
}

func Init(config *Config) {
	var c httpcache.Cache
	if config.CacheSize > 0 {
		if len(config.CacheDir) > 0 {
			d := diskv.New(diskv.Options{
				BasePath:     config.CacheDir,
				CacheSizeMax: config.CacheSize * 1024 * 1024,
			})
			c = diskcache.NewWithDiskv(d)
		} else {
			c = httpcache.NewMemoryCache()
		}
	}

	p := NewProxy(nil, c)
	if len(config.Whitelist) > 0 {
		p.Whitelist = strings.Split(config.Whitelist, ",")
	}
	if len(config.Referrers) > 0 {
		p.Referrers = strings.Split(config.Referrers, ",")
	}
	if len(config.SignatureKey) > 0 {
		key := []byte(config.SignatureKey)
		if strings.HasPrefix(config.SignatureKey, "@") {
			file := strings.TrimPrefix(config.SignatureKey, "@")
			var err error
			key, err = ioutil.ReadFile(file)
			if err != nil {
				log.Fatalf("error reading signature file: %v", err)
			}
		}
		p.SignatureKey = key
	}
	if len(config.BaseURL) > 0 {
		var err error
		p.DefaultBaseURL, err = url.Parse(config.BaseURL)
		if err != nil {
			log.Fatalf("error parsing baseURL: %v", err)
		}
	}

	p.ScaleUp = config.ScaleUp
	p.ResrcVKey = config.ResrcVKey
	p.ValidToken = config.ValidToken
	p.CrossParams = true
	p.CleanPrefix = "/imageproxy"
	println(p.BuildURL("200x300", "http://www.coscms.com/images/201411/source_img/13815_G_1416788524324.jpg"))
	Default = p
}
