package namecheap

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager/ddnserrors"
	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager/domain/dnsdomain"
	"github.com/webx-top/echo"
)

const (
	namecheapEndpoint       = "https://dynamicdns.park-your-domain.com/update"
	signUpURL               = ``
	docLineType             = ``
	defaultTTL              = 300
	defaultTimeout    int64 = 30
)

type Namecheap struct {
	clientSecret string
	Domains      []*dnsdomain.Domain
	//TTL          int
	client *http.Client
}

const Name = `Namecheap`

func (*Namecheap) Name() string {
	return Name
}

func (*Namecheap) Description() string {
	return ``
}

func (*Namecheap) SignUpURL() string {
	return signUpURL
}

func (*Namecheap) LineTypeURL() string {
	return docLineType
}

var configItems = echo.KVList{
	echo.NewKV(`clientSecret`, `Password`).SetHKV(`inputType`, `text`).SetHKV(`helpBlock`, `密码`).SetHKV(`required`, true),
	echo.NewKV(`timeout`, `接口超时(秒)`).SetHKV(`inputType`, `number`).SetX(defaultTimeout),
	//echo.NewKV(`ttl`, `域名TTL`).SetHKV(`inputType`, `number`).SetX(defaultTTL),
}

func (*Namecheap) ConfigItems() echo.KVList {
	return configItems
}

var support = dnsdomain.Support{
	A:    true,
	AAAA: false,
}

func (*Namecheap) Support() dnsdomain.Support {
	return support
}

// Init 初始化
func (p *Namecheap) Init(settings echo.H, domains []*dnsdomain.Domain) error {
	p.clientSecret = settings.String(`clientSecret`)
	// p.TTL = settings.Int(`ttl`)
	// if p.TTL < 1 {
	// 	p.TTL = defaultTTL
	// }
	p.Domains = domains
	timeout := settings.Int64(`timeout`)
	if timeout < 1 {
		timeout = defaultTimeout
	}
	p.client = &http.Client{Timeout: time.Duration(timeout) * time.Second}
	return nil
}

func (p *Namecheap) Update(ctx context.Context, recordType string, ipAddr string) error {
	for _, domain := range p.Domains {
		values := url.Values{}
		values.Set("host", domain.GetSubDomain())
		values.Set("domain", domain.DomainName)
		values.Set("password", p.clientSecret)
		ip := domain.IP(ipAddr)
		values.Set("ip", ip)

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, namecheapEndpoint+`?`+values.Encode(), nil)
		if err != nil {
			domain.UpdateStatus = dnsdomain.UpdatedFailed
			return err
		}
		request.Header.Set(`Accept`, "application/xml")

		response, err := p.client.Do(request)
		if err != nil {
			domain.UpdateStatus = dnsdomain.UpdatedFailed
			return err
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(response.Body)
			return fmt.Errorf("%w: %d: %s",
				ddnserrors.ErrBadHTTPStatus, response.StatusCode, string(b))
		}

		decoder := xml.NewDecoder(response.Body)
		var parsedXML struct {
			Errors struct {
				Error string `xml:"errors.Err1"`
			} `xml:"errors"`
			IP string `xml:"IP"`
		}
		if err := decoder.Decode(&parsedXML); err != nil {
			domain.UpdateStatus = dnsdomain.UpdatedFailed
			return fmt.Errorf("%w: %s", ddnserrors.ErrUnmarshalResponse, err)
		}

		if len(parsedXML.Errors.Error) > 0 {
			domain.UpdateStatus = dnsdomain.UpdatedFailed
			return fmt.Errorf("%w: %s", ddnserrors.ErrUnsuccessfulResponse, parsedXML.Errors.Error)
		}
		if parsedXML.IP != ip {
			domain.UpdateStatus = dnsdomain.UpdatedFailed
			return fmt.Errorf("%w: %s", ddnserrors.ErrIPReceivedMismatch, parsedXML.IP)
		}
		domain.UpdateStatus = dnsdomain.UpdatedSuccess
	}
	return nil
}
