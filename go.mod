module github.com/admpub/nging

go 1.14

replace github.com/caddyserver/caddy => github.com/admpub/caddy v1.1.5

replace google.golang.org/protobuf => github.com/protocolbuffers/protobuf-go v1.25.0

replace golang.org/x/crypto => github.com/golang/crypto v0.0.0-20200820211705-5c72a883971a

replace golang.org/x/image => github.com/golang/image v0.0.0-20200924062109-4578eab98f00

replace golang.org/x/net => github.com/golang/net v0.0.0-20200925080053-05aa5d4ee321

replace golang.org/x/oauth2 => github.com/golang/oauth2 v0.0.0-20200902213428-5d25da1a8d43

replace golang.org/x/protobuf => github.com/golang/protobuf v1.4.2

replace golang.org/x/time => github.com/golang/time v0.0.0-20200630173020-3af7569d3a1e

replace golang.org/x/sync => github.com/golang/sync v0.0.0-20200625203802-6e8e738ad208

replace golang.org/x/sys => github.com/golang/sys v0.0.0-20200923182605-d9f96fdee20d

replace golang.org/x/text => github.com/golang/text v0.3.3

replace golang.org/x/tools => github.com/golang/tools v0.0.0-20200925191224-5d1fdd8fa346

replace golang.org/x/xerrors => github.com/golang/xerrors v0.0.0-20200804184101-5ec99f83aff1

replace github.com/admpub/xgo => github.com/admpub/xgo v1.14.6

replace github.com/admpub/bindata/v3 => github.com/admpub/bindata/v3 v3.1.5

require (
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/admpub/9t v0.0.0-20190605154903-a68069ace5e1
	github.com/admpub/archiver v1.1.4
	github.com/admpub/bindata/v3 v3.1.5
	github.com/admpub/ccs-gm v0.0.3
	github.com/admpub/checksum v1.0.1
	github.com/admpub/color v1.7.0
	github.com/admpub/confl v0.0.0-20190331072055-254deeac709e
	github.com/admpub/copier v0.0.0-20200812014131-931651b20f74 // indirect
	github.com/admpub/cr v0.0.0-20200630080251-9947244796af
	github.com/admpub/cron v0.0.1
	github.com/admpub/decimal v0.0.0-20180709203117-cd690d0c9e24 // indirect
	github.com/admpub/dgoogauth v0.0.0-20170926052827-752650e076f2
	github.com/admpub/email v0.0.0-20180510084211-36989e527d89
	github.com/admpub/errors v0.8.2
	github.com/admpub/events v0.0.0-20190913050400-383beb0843c6
	github.com/admpub/fasthttp v0.0.0-20200705060903-f0f441048f4a // indirect
	github.com/admpub/frp v0.28.1
	github.com/admpub/fsnotify v1.4.4 // indirect
	github.com/admpub/ftpserver v0.0.0-20180821023651-6b950445d653
	github.com/admpub/gifresize v1.0.2 // indirect
	github.com/admpub/go-bindata-assetfs v0.0.0-20170428090253-36eaa4c19588
	github.com/admpub/go-download v2.1.2+incompatible
	github.com/admpub/go-figure v0.0.0-20180619031829-18b2b544842c
	github.com/admpub/go-isatty v0.0.9
	github.com/admpub/go-phantomjs-fetcher v0.0.0-20180924162325-bb2ae1648e33
	github.com/admpub/go-pretty v3.3.3+incompatible
	github.com/admpub/go-ps v0.0.1 // indirect
	github.com/admpub/go-reuseport v0.0.3 // indirect
	github.com/admpub/go-utility v0.0.1 // indirect
	github.com/admpub/godownloader v0.0.0-20191013090831-d509172649e4
	github.com/admpub/goforever v0.1.1
	github.com/admpub/gohls v0.0.0-20191013012052-b7505aaf3c90 // indirect
	github.com/admpub/gohls-server v0.3.3 // indirect
	github.com/admpub/gohttp v0.0.0-20190322032039-b55c707b8f1e
	github.com/admpub/goloader v0.0.0-20200821104313-abaac9b4b83f
	github.com/admpub/gopiper v1.0.1
	github.com/admpub/goseaweedfs v0.1.2
	github.com/admpub/highwayhash v0.0.0-20180501080913-85fc8a2dacad
	github.com/admpub/httpscerts v0.0.0-20180907121630-a2990e2af45c
	github.com/admpub/humanize v0.0.0-20190501023926-5f826e92c8ca // indirect
	github.com/admpub/i18n v0.0.0-20190425064330-383d24a2fded // indirect
	github.com/admpub/identicon v1.0.2 // indirect
	github.com/admpub/imageproxy v0.9.1
	github.com/admpub/imaging v1.5.0 // indirect
	github.com/admpub/ini v1.38.2
	github.com/admpub/ip2region v1.2.5
	github.com/admpub/license_gen v0.0.0-20181201145035-5d7646743eec
	github.com/admpub/log v0.0.0-20191027043925-a6c03a7421a3
	github.com/admpub/logcool v0.3.1
	github.com/admpub/mahonia v0.0.0-20151019004008-c528b747d92d
	github.com/admpub/mail v0.0.0-20170408110349-d63147b0317b
	github.com/admpub/marmot v0.0.0-20200702042226-2170d9ff59f5
	github.com/admpub/metrohash v0.0.0-20160821164112-8d1c8b6bed28
	github.com/admpub/mysql-schema-sync v0.0.4
	github.com/admpub/null v8.0.1+incompatible
	github.com/admpub/pester v0.0.0-20200411024648-005672a2bd48 // indirect
	github.com/admpub/phantomjs v0.0.0-20180924162111-8a5af756140e
	github.com/admpub/qrcode v0.0.0-20190512103923-0a2b75fa9edc
	github.com/admpub/queueChan v0.0.0-20151001074356-79908f7a499f // indirect
	github.com/admpub/regexp2 v1.1.7
	github.com/admpub/restful v0.0.0-20180328144945-2e8d62d39607
	github.com/admpub/securecookie v0.0.0-20170722041919-69560e375596
	github.com/admpub/service v0.0.0-20200628035946-b971165ffff2
	github.com/admpub/sessions v0.0.0-20200615143229-e7dd49c7c4f5 // indirect
	github.com/admpub/snowflake v0.0.0-20180412010544-68117e6bbede
	github.com/admpub/sockjs-go v0.0.0-20170208085255-715e9716fc23
	github.com/admpub/sonyflake v0.0.0-20160530021500-fa881fb1052b
	github.com/admpub/sqlboiler v3.0.1+incompatible // indirect
	github.com/admpub/statik v0.1.7 // indirect
	github.com/admpub/tail v1.0.2
	github.com/admpub/useragent v0.0.0-20190806155403-63e85649b0f2
	github.com/admpub/web-terminal v0.0.0-20190705124712-e503d17936e9
	github.com/admpub/websocket v1.0.1
	github.com/araddon/gou v0.0.0-20190110011759-c797efecbb61 // indirect
	github.com/aws/aws-sdk-go v1.35.11
	github.com/axgle/mahonia v0.0.0-20180208002826-3358181d7394 // indirect
	github.com/bitly/go-simplejson v0.5.0 // indirect
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/caddy-plugins/caddy-expires v1.1.1
	github.com/caddy-plugins/caddy-filter v0.15.0
	github.com/caddy-plugins/caddy-locale v0.0.0-20190704155156-288438ce0a5e
	github.com/caddy-plugins/caddy-prometheus v0.0.0-20190704154614-d29127a2871c
	github.com/caddy-plugins/caddy-rate-limit v1.6.1
	github.com/caddy-plugins/caddy-s3browser v0.0.1
	github.com/caddy-plugins/cors v0.0.0-20190704155148-3c98079f1197
	github.com/caddy-plugins/ipfilter v1.1.4
	github.com/caddy-plugins/nobots v0.1.1
	github.com/caddyserver/caddy v1.1.5
	github.com/chromedp/cdproto v0.0.0-20201009231348-1c6a710e77de
	github.com/chromedp/chromedp v0.5.3
	github.com/codegangsta/inject v0.0.0-20150114235600-33e0aa1cb7c0 // indirect
	github.com/coscms/go-imgparse v0.0.0-20150925144422-3e3a099f7856
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/disintegration/imaging v1.6.2 // indirect
	github.com/fatedier/beego v1.7.2 // indirect
	github.com/fatedier/golib v0.2.0 // indirect
	github.com/fatedier/kcp-go v2.0.4-0.20190803094908-fe8645b0a904+incompatible // indirect
	github.com/fatih/color v1.9.0
	github.com/fd/go-shellwords v0.0.0-20130603174837-6a119423524d // indirect
	github.com/francoispqt/gojay v1.2.13 // indirect
	github.com/friendsofgo/errors v0.9.2 // indirect
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/go-openapi/strfmt v0.19.5 // indirect
	github.com/goftp/file-driver v0.0.0-20180502053751-5d604a0fc0c9 // indirect
	github.com/goftp/server v0.0.0-20200708154336-f64f7c2d8a42 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/gomodule/redigo v1.8.2
	github.com/grafov/m3u8 v0.11.1 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/hashicorp/yamux v0.0.0-20200609203250-aecfd211c9ce // indirect
	github.com/igm/sockjs-go v3.0.0+incompatible // indirect
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/jedib0t/go-pretty v4.3.0+incompatible // indirect
	github.com/jesseduffield/lazygit v0.22.9
	github.com/jlaffaye/ftp v0.0.0-20200812143550-39e3779af0db // indirect
	github.com/kardianos/service v1.1.0 // indirect
	github.com/markbates/goth v1.64.2 // indirect
	github.com/mattn/go-colorable v0.1.8
	github.com/mattn/go-runewidth v0.0.9
	github.com/mattn/go-sqlite3 v1.14.0 // indirect
	github.com/mholt/certmagic v0.8.3
	github.com/microcosm-cc/bluemonday v1.0.4
	github.com/minio/minio-go v6.0.14+incompatible
	github.com/muesli/smartcrop v0.3.0 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/oschwald/maxminddb-golang v1.7.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.12.0
	github.com/prometheus/client_golang v1.7.1 // indirect
	github.com/russross/blackfriday v1.5.2
	github.com/rwcarlsen/goexif v0.0.0-20190401172101-9e8deecbddbd // indirect
	github.com/saintfish/chardet v0.0.0-20120816061221-3af4cd4741ca // indirect
	github.com/shirou/gopsutil v2.20.9+incompatible
	github.com/shivakar/metrohash v0.0.0-20160821164112-8d1c8b6bed28 // indirect
	github.com/shivakar/xxhash v0.0.0-20160821164220-5ea66fb45566 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.6.1
	github.com/syndtr/goleveldb v1.0.0
	github.com/tdewolff/minify v2.3.6+incompatible // indirect
	github.com/tdewolff/parse v2.3.4+incompatible // indirect
	github.com/tdewolff/test v1.0.6 // indirect
	github.com/tebeka/selenium v0.9.9
	github.com/tuotoo/qrcode v0.0.0-20190222102259-ac9c44189bf2 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/volatiletech/inflect v0.0.1 // indirect
	github.com/volatiletech/sqlboiler v3.7.1+incompatible // indirect
	github.com/webx-top/captcha v0.0.0-20161202061115-29e9e7f30aa0
	github.com/webx-top/chardet v0.0.0-20180930194453-2f90d95f7b7f // indirect
	github.com/webx-top/client v0.1.3
	github.com/webx-top/codec v0.0.0-20200914105801-3782d81a0302
	github.com/webx-top/com v0.0.4
	github.com/webx-top/db v1.1.2
	github.com/webx-top/echo v2.3.8+incompatible
	github.com/webx-top/image v0.0.1
	github.com/webx-top/pagination v0.0.1
	github.com/webx-top/tagfast v0.0.0-20161020041435-9a2065ce3dd2 // indirect
	github.com/webx-top/validation v0.0.1 // indirect
	github.com/willf/bitset v1.1.11 // indirect
	github.com/xtaci/lossyconn v0.0.0-20200209145036-adba10fffc37 // indirect
	github.com/xu42/youzan-sdk-go v0.0.3
	gocloud.dev v0.20.0
	golang.org/x/crypto v0.0.0-20201012173705-84dcc777aaee
	golang.org/x/net v0.0.0-20201010224723-4f7140c49acb // indirect
	golang.org/x/sys v0.0.0-20201013132646-2da7054afaeb // indirect
	gopkg.in/fsnotify/fsnotify.v1 v1.4.7 // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
)
