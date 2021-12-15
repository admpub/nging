package ua

import (
	"bytes"
	"regexp"
	"strings"
)

// UserAgent struct containg all determined datra from parsed user-agent string
type UserAgent struct {
	Name      string
	Version   string
	OS        string
	OSVersion string
	Device    string
	Mobile    bool
	Tablet    bool
	Desktop   bool
	Bot       bool
	URL       string
	String    string
	tokens    properties
}

func (ua UserAgent) Tokens() properties {
	return ua.tokens
}

var ignore = map[string]struct{}{
	"KHTML, like Gecko": struct{}{},
	"U":                 struct{}{},
	"compatible":        struct{}{},
	"Mozilla":           struct{}{},
	"WOW64":             struct{}{},
}

// Constants for browsers and operating systems for easier comparation
const (
	Windows      = "Windows"
	WindowsPhone = "Windows Phone"
	Android      = "Android"
	MacOS        = "macOS"
	IOS          = "iOS"
	Linux        = "Linux"

	Opera            = "Opera"
	OperaMini        = "Opera Mini"
	OperaTouch       = "Opera Touch"
	Chrome           = "Chrome"
	Firefox          = "Firefox"
	InternetExplorer = "Internet Explorer"
	Safari           = "Safari"
	Edge             = "Edge"
	Vivaldi          = "Vivaldi"

	WxWork         = "wxwork"         //企业微信浏览器
	MicroMessenger = "MicroMessenger" //微信浏览器
	AlipayClient   = "AlipayClient"   //支付宝客户端

	Googlebot           = "Googlebot"
	Twitterbot          = "Twitterbot"
	FacebookExternalHit = "facebookexternalhit"
	Bingbot             = "bingbot"
	Baidubot            = "Baiduspider"
	Sobot               = "360Spider"
	Yahoobot            = "Yahoo! Slurp"
	Sososbot            = "Sosospider"
	IAskbot             = "iaskspider" //Sina iask
	Sogoubot            = "Sogou web spider"
	Yodaobot            = "YodaoBot"
	MSNbot              = "msnbot"
	SemrushBot          = "SemrushBot"
)

// Parse user agent string returning UserAgent struct
func Parse(userAgent string) UserAgent {
	ua := UserAgent{
		String: userAgent,
	}

	tokens := parse(userAgent)

	// check is there URL
	for k := range tokens {
		if strings.HasPrefix(k, "http://") || strings.HasPrefix(k, "https://") {
			ua.URL = k
			delete(tokens, k)
			break
		}
	}

	// OS lookup
	switch {
	case tokens.exists("Android"):
		ua.OS = Android
		ua.OSVersion = tokens[Android]
		for s := range tokens {
			if strings.HasSuffix(s, "Build") {
				ua.Device = strings.TrimSpace(s[:len(s)-5])
				ua.Tablet = strings.Contains(strings.ToLower(ua.Device), "tablet")
			}
		}

	case tokens.exists("iPhone"):
		ua.OS = IOS
		ua.OSVersion = tokens.findMacOSVersion()
		ua.Device = "iPhone"
		ua.Mobile = true

	case tokens.exists("iPad"):
		ua.OS = IOS
		ua.OSVersion = tokens.findMacOSVersion()
		ua.Device = "iPad"
		ua.Tablet = true

	case tokens.exists("Windows NT"):
		ua.OS = Windows
		ua.OSVersion = tokens["Windows NT"]
		ua.Desktop = true

	case tokens.exists("Windows Phone OS"):
		ua.OS = WindowsPhone
		ua.OSVersion = tokens["Windows Phone OS"]
		ua.Mobile = true

	case tokens.exists("Macintosh"):
		ua.OS = MacOS
		ua.OSVersion = tokens.findMacOSVersion()
		ua.Desktop = true

	case tokens.exists("Linux"):
		ua.OS = Linux
		ua.OSVersion = tokens[Linux]
		ua.Desktop = true

	}

	// for s, val := range sys {
	// 	fmt.Println(s, "--", val)
	// }

	switch {

	case tokens.exists("Googlebot"):
		ua.Name = Googlebot
		ua.Version = tokens[Googlebot]
		ua.Bot = true
		ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")

	case tokens["AlipayClient"] != "":
		ua.Name = AlipayClient
		ua.Version = tokens[AlipayClient]
		ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")

	case tokens["wxwork"] != "":
		ua.Name = WxWork
		ua.Version = tokens[WxWork]
		ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")

	case tokens["MicroMessenger"] != "":
		ua.Name = MicroMessenger
		ua.Version = tokens[MicroMessenger]
		ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")

	case tokens["Opera Mini"] != "":
		ua.Name = OperaMini
		ua.Version = tokens[OperaMini]
		ua.Mobile = true

	case tokens["Opera Touch"] != "":
		ua.Name = OperaTouch
		ua.Version = tokens[OperaTouch]
		ua.Tablet = true

	case tokens["Opera"] != "":
		ua.Name = Opera
		ua.Version = tokens[Opera]
		ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")

	case tokens["OPR"] != "":
		ua.Name = Opera
		ua.Version = tokens["OPR"]
		ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")

	case tokens["OPT"] != "":
		ua.Name = OperaTouch
		ua.Version = tokens["OPT"]
		ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")

	// Opera on iOS
	case tokens["OPiOS"] != "":
		ua.Name = Opera
		ua.Version = tokens["OPiOS"]
		ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")

	// Chrome on iOS
	case tokens["CriOS"] != "":
		ua.Name = Chrome
		ua.Version = tokens["CriOS"]
		ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")

	// Firefox on iOS
	case tokens["FxiOS"] != "":
		ua.Name = Firefox
		ua.Version = tokens["FxiOS"]
		ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")

	case tokens["Firefox"] != "":
		ua.Name = Firefox
		ua.Version = tokens[Firefox]
		_, ua.Mobile = tokens["Mobile"]
		_, ua.Tablet = tokens["Tablet"]

	case tokens["Vivaldi"] != "":
		ua.Name = Vivaldi
		ua.Version = tokens[Vivaldi]

	case tokens.exists("MSIE"):
		ua.Name = InternetExplorer
		ua.Version = tokens["MSIE"]

	case tokens["Edge"] != "":
		ua.Name = Edge
		ua.Version = tokens["Edge"]
		ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")

	case tokens["EdgA"] != "":
		ua.Name = Edge
		ua.Version = tokens["EdgA"]
		ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")

	case tokens["bingbot"] != "":
		ua.Name = Bingbot
		ua.Version = tokens["bingbot"]
		ua.Bot = true
		ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")

	case tokens["SamsungBrowser"] != "":
		ua.Name = "Samsung Browser"
		ua.Version = tokens["SamsungBrowser"]
		ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")

	// if chrome and Safari defined, find any other tokensent descr
	case tokens.exists(Chrome) && tokens.exists(Safari):
		name := tokens.findBestMatch(true)
		if name != "" {
			ua.Name = name
			ua.Version = tokens[name]
			break
		}
		fallthrough

	case tokens.exists("Chrome"):
		ua.Name = Chrome
		ua.Version = tokens["Chrome"]
		ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")

	case tokens.exists("Safari"):
		ua.Name = Safari
		if v, ok := tokens["Version"]; ok {
			ua.Version = v
		} else {
			ua.Version = tokens["Safari"]
		}
		ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")

	default:
		if ua.OS == "Android" && tokens["Version"] != "" {
			ua.Name = "Android browser"
			ua.Version = tokens["Version"]
			ua.Mobile = true
		} else {
			if name := tokens.findBestMatch(false); name != "" {
				ua.Name = name
				ua.Version = tokens[name]
			} else {
				ua.Name = ua.String
			}
			name := strings.ToLower(ua.Name)
			ua.Bot = strings.Contains(name, "bot") || strings.Contains(name, "spider")
			ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")
		}
	}

	// if tabler, switch mobile to off
	if ua.Tablet {
		ua.Mobile = false
	}

	// if not already bot, check some popular bots and weather URL is set
	if !ua.Bot {
		ua.Bot = ua.URL != ""
	}

	if !ua.Bot {
		switch ua.Name {
		case FacebookExternalHit, Yahoobot:
			ua.Bot = true
		}
	}

	ua.tokens = tokens
	return ua
}

func parse(userAgent string) (clients properties) {
	clients = make(map[string]string, 0)
	slash := false
	isURL := false
	var buff, val bytes.Buffer
	addToken := func() {
		if buff.Len() != 0 {
			s := strings.TrimSpace(buff.String())
			if _, ign := ignore[s]; !ign {
				if isURL {
					s = strings.TrimPrefix(s, "+")
				}

				if val.Len() == 0 { // only if value don't exists
					var ver string
					s, ver = checkVer(s) // determin version string and split
					clients[s] = ver
				} else {
					clients[s] = strings.TrimSpace(val.String())
				}
			}
		}
		buff.Reset()
		val.Reset()
		slash = false
		isURL = false
	}

	parOpen := false

	bua := []byte(userAgent)
	for i, c := range bua {

		//fmt.Println(string(c), c)
		switch {
		case c == 41: // )
			addToken()
			parOpen = false

		case parOpen && c == 59: // ;
			addToken()

		case c == 40: // (
			addToken()
			parOpen = true

		case slash && c == 32:
			addToken()

		case slash:
			val.WriteByte(c)

		case c == 47 && !isURL: //   /
			if i != len(bua)-1 && bua[i+1] == 47 && (bytes.HasSuffix(buff.Bytes(), []byte("http:")) || bytes.HasSuffix(buff.Bytes(), []byte("https:"))) {
				buff.WriteByte(c)
				isURL = true
			} else {
				slash = true
			}

		default:
			buff.WriteByte(c)
		}
	}
	addToken()

	return clients
}

func checkVer(s string) (name, v string) {
	i := strings.LastIndex(s, " ")
	if i == -1 {
		return s, ""
	}

	//v = s[i+1:]

	switch s[:i] {
	case "Linux", "Windows NT", "Windows Phone OS", "MSIE", "Android":
		return s[:i], s[i+1:]
	default:
		return s, ""
	}

	// for _, c := range v {
	// 	if (c >= 48 && c <= 57) || c == 46 {
	// 	} else {
	// 		return s, ""
	// 	}
	// }

	// return s[:i], s[i+1:]

}

type properties map[string]string

func (p properties) exists(key string) bool {
	_, ok := p[key]
	return ok
}

func (p properties) existsAny(keys ...string) bool {
	for _, k := range keys {
		if _, ok := p[k]; ok {
			return true
		}
	}
	return false
}

func (p properties) findMacOSVersion() string {
	for k, v := range p {
		if strings.Contains(k, "OS") {
			if ver := findVersion(v); ver != "" {
				return ver
			} else if ver = findVersion(k); ver != "" {
				return ver
			}
		}
	}
	return ""
}

// findBestMatch from the rest of the bunch
// in first cycle only return key vith version value
// if withVerValue is false, do another cycle and return any token
func (p properties) findBestMatch(withVerOnly bool) string {
	n := 2
	if withVerOnly {
		n = 1
	}
	for i := 0; i < n; i++ {
		for k, v := range p {
			switch k {
			case Chrome, Firefox, Safari, Opera, Edge, "Version", "Mobile", "Mobile Safari", "Mozilla", "AppleWebKit", "Windows NT", "Windows Phone OS", Android, "Macintosh", Linux, "GSA":
			default:
				if i == 0 {
					if v != "" { // in first check, only return  keys with value
						return k
					}
				} else {
					return k
				}
			}
		}
	}
	return ""
}

var rxMacOSVer = regexp.MustCompile("[_\\d\\.]+")

func findVersion(s string) string {
	if ver := rxMacOSVer.FindString(s); ver != "" {
		return strings.Replace(ver, "_", ".", -1)
	}
	return ""
}
