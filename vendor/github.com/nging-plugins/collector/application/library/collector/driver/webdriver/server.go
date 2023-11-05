package webdriver

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/admpub/marmot/miner"
	"github.com/nging-plugins/collector/application/library/collector"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"github.com/tebeka/selenium/firefox"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

var (
	servers     = sync.Map{}
	serivceAPI  = "http://127.0.0.1:%d/wd/hub"
	defaultPort = 4444
)

func ChromeDriverDefaultPath() string {
	wd := echo.Wd()
	if !com.IsExist(filepath.Join(wd, `support`)) {
		gopath := os.Getenv(`GOPATH`)
		if len(gopath) > 0 {
			wd = filepath.Join(gopath, `src/github.com/admpub/nging`)
		}
	}
	switch runtime.GOOS {
	case `windows`:
		return filepath.Join(wd, `support`, `chromedriver_386.exe`)
	case `linux`, `darwin`:
		return filepath.Join(wd, `support`, `chromedriver_`+runtime.GOOS+`_`+runtime.GOARCH)
	default:
		return ``
	}
}

func StartChromeServer(driverPath string, port int, opts ...selenium.ServiceOption) (service *selenium.Service, err error) {
	//chromedriver 可从 http://npm.taobao.org/mirrors/chromedriver/ 下载
	if len(driverPath) == 0 {
		driverPath = `chromedriver`
	}
	// 启动chromedriver，端口号可自定义
	service, err = selenium.NewChromeDriverService(driverPath, port, opts...)
	if err != nil {
		log.Printf("Error starting the ChromeDriver server: %v", err)
	}
	//defer service.Stop()
	return
}

func chromeDriverOption(caps selenium.Capabilities, cfg *collector.Base) {
	// 禁止加载图片，加快渲染速度
	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}
	var userAgent string
	if len(cfg.UserAgent) > 0 {
		userAgent = cfg.UserAgent
	} else {
		userAgent = miner.RandomUserAgent()
	}
	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Path:  "",
		Args: []string{
			"--headless", // 设置Chrome无头模式
			"--no-sandbox",
			"--disable-gpu",
			"--user-agent=" + userAgent, // 模拟user-agent，防反爬
			//"window-size=1200x600"
		},
	}
	if cfg.Header != nil {
		for k, v := range cfg.Header {
			chromeCaps.Args = append(chromeCaps.Args, k+`=`+fmt.Sprintf(`%q`, v))
		}
	}
	caps.AddChrome(chromeCaps)
}

func firefoxDriverOption(caps selenium.Capabilities, cfg *collector.Base) {
	imagCaps := map[string]interface{}{}
	var userAgent string
	if len(cfg.UserAgent) > 0 {
		userAgent = cfg.UserAgent
	} else {
		userAgent = miner.RandomUserAgent()
	}
	firefoxCaps := firefox.Capabilities{
		Prefs:  imagCaps,
		Binary: "",
		Args: []string{
			"--headless", // 设置Firefox无头模式
			"--disable-gpu",
			"--user-agent=" + userAgent, // 模拟user-agent，防反爬
			//"window-size=1200x600"
		},
	}
	if cfg.Header != nil {
		for k, v := range cfg.Header {
			firefoxCaps.Args = append(firefoxCaps.Args, k+`=`+fmt.Sprintf(`%q`, v))
		}
	}
	caps.AddFirefox(firefoxCaps)
}

func StartServer(cfg *collector.Base) (service *selenium.Service, err error) {
	opt := cfg.Extra
	browserName := opt.String(`browserName`, `chrome`)
	if v, y := servers.Load(browserName); y && v != nil {
		service = v.(*selenium.Service)
		return
	}
	//chromedriver 可从 http://npm.taobao.org/mirrors/chromedriver/ 下载
	driverPath := opt.String(`driverPath`, ChromeDriverDefaultPath())
	servicePort := opt.Int(`servicePort`, defaultPort)

	var opts []selenium.ServiceOption
	switch v := opt.Get(`serviceOption`).(type) {
	case []selenium.ServiceOption:
		opts = v
	case selenium.ServiceOption:
		opts = append(opts, v)
	}

	// set browser as chrome
	caps := selenium.Capabilities{"browserName": browserName}
	switch browserName {
	case `chrome`:
		chromeDriverOption(caps, cfg)
		service, err = selenium.NewChromeDriverService(driverPath, servicePort, opts...)
	case `firefox`:
		firefoxDriverOption(caps, cfg)
		service, err = selenium.NewGeckoDriverService(driverPath, servicePort, opts...)
	default:
		if chromeDriverPath := opt.String(`chromeDriverPath`); len(chromeDriverPath) > 0 {
			opts = append(opts, selenium.ChromeDriver(chromeDriverPath))
		}
		if geckoDriverPath := opt.String(`geckoDriverPath`); len(geckoDriverPath) > 0 {
			opts = append(opts, selenium.GeckoDriver(geckoDriverPath))
		}
		if javaPath := opt.String(`javaPath`); len(javaPath) > 0 {
			opts = append(opts, selenium.JavaPath(javaPath))
		}
		if htmlUnitPath := opt.String(`htmlUnitPath`); len(htmlUnitPath) > 0 {
			opts = append(opts, selenium.HTMLUnit(htmlUnitPath))
		}
		service, err = selenium.NewSeleniumService(driverPath, servicePort, opts...)
	}
	if err == nil {
		servers.Store(browserName, service)
	}
	return
}

func CloseServer(browserName ...string) error {
	if len(browserName) < 1 {
		servers.Range(func(key, val interface{}) bool {
			val.(*selenium.Service).Stop()
			servers.Delete(key)
			return true
		})
		return nil
	}
	if dri, ok := servers.Load(browserName[0]); ok {
		servers.Delete(browserName[0])
		return dri.(*selenium.Service).Stop()
	}
	return nil
}
