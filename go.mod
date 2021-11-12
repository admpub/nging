module github.com/admpub/nging/v3

go 1.17

exclude github.com/gomodule/redigo v2.0.0+incompatible

require (
	gitee.com/admpub/certmagic v0.8.8
	github.com/PuerkitoBio/goquery v1.8.0 // indirect
	github.com/abh/errorutil v0.0.0-20130729183701-f9bd360d00b9
	github.com/admpub/9t v0.0.0-20190605154903-a68069ace5e1
	github.com/admpub/archiver v1.1.4
	github.com/admpub/bindata/v3 v3.1.5
	github.com/admpub/caddy v1.1.11
	github.com/admpub/ccs-gm v0.0.3
	github.com/admpub/checksum v1.0.1
	github.com/admpub/color v1.8.0
	github.com/admpub/confl v0.1.0
	github.com/admpub/cr v0.0.2
	github.com/admpub/cron v0.0.1
	github.com/admpub/dgoogauth v0.0.1
	github.com/admpub/email v2.4.1+incompatible
	github.com/admpub/errors v0.8.2
	github.com/admpub/events v1.2.0
	github.com/admpub/fasthttp v0.0.3 // indirect
	github.com/admpub/frp v0.37.1
	github.com/admpub/go-bindata-assetfs v0.0.0-20170428090253-36eaa4c19588
	github.com/admpub/go-download/v2 v2.1.12
	github.com/admpub/go-figure v0.0.0-20180619031829-18b2b544842c
	github.com/admpub/go-isatty v0.0.10
	github.com/admpub/go-password v0.1.3
	github.com/admpub/go-phantomjs-fetcher v0.0.0-20180924162325-bb2ae1648e33
	github.com/admpub/go-pretty/v6 v6.0.3
	github.com/admpub/go-ps v0.0.1 // indirect
	github.com/admpub/go-sshclient v0.0.1
	github.com/admpub/go-utility v0.0.1 // indirect
	github.com/admpub/godotenv v1.4.2
	github.com/admpub/godownloader v2.0.3+incompatible
	github.com/admpub/goforever v0.1.1
	github.com/admpub/gohls v1.0.7 // indirect
	github.com/admpub/gohls-server v0.3.5 // indirect
	github.com/admpub/gohttp v0.0.0-20190322032039-b55c707b8f1e // indirect
	github.com/admpub/gopiper v1.0.2
	github.com/admpub/goseaweedfs v0.1.2
	github.com/admpub/httpscerts v0.0.0-20180907121630-a2990e2af45c
	github.com/admpub/i18n v0.0.2 // indirect
	github.com/admpub/identicon v1.0.2 // indirect
	github.com/admpub/imageproxy v0.9.3
	github.com/admpub/imaging v1.5.0 // indirect
	github.com/admpub/ip2region v1.2.11
	github.com/admpub/license_gen v0.1.0
	github.com/admpub/log v1.3.2
	github.com/admpub/logcool v0.3.2
	github.com/admpub/mahonia v0.0.0-20151019004008-c528b747d92d
	github.com/admpub/mail v0.0.0-20170408110349-d63147b0317b
	github.com/admpub/marmot v0.0.0-20200702042226-2170d9ff59f5
	github.com/admpub/mysql-schema-sync v0.2.1
	github.com/admpub/null v8.0.4+incompatible
	github.com/admpub/once v0.0.1
	github.com/admpub/pester v0.0.0-20200411024648-005672a2bd48 // indirect
	github.com/admpub/phantomjs v0.0.0-20180924162111-8a5af756140e
	github.com/admpub/qrcode v0.0.2
	github.com/admpub/randomize v0.0.2 // indirect
	github.com/admpub/redistore v1.1.0 // indirect
	github.com/admpub/regexp2 v1.1.7
	github.com/admpub/resty/v2 v2.7.0
	github.com/admpub/securecookie v1.1.2
	github.com/admpub/service v0.0.2
	github.com/admpub/sessions v0.1.2 // indirect
	github.com/admpub/snowflake v0.0.0-20180412010544-68117e6bbede
	github.com/admpub/sockjs-go/v3 v3.0.0
	github.com/admpub/sonyflake v0.0.1
	github.com/admpub/tail v1.1.0
	github.com/admpub/timeago v1.1.3
	github.com/admpub/useragent v0.0.0-20190806155403-63e85649b0f2
	github.com/admpub/web-terminal v0.0.1
	github.com/admpub/websocket v1.0.4
	github.com/apache/thrift v0.15.0 // indirect
	github.com/arl/statsviz v0.4.0
	github.com/armon/go-metrics v0.3.10 // indirect
	github.com/aws/aws-sdk-go v1.42.3
	github.com/axgle/mahonia v0.0.0-20180208002826-3358181d7394 // indirect
	github.com/bitly/go-simplejson v0.5.0 // indirect
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/caddy-plugins/caddy-expires v1.1.2
	github.com/caddy-plugins/caddy-filter v0.15.2
	github.com/caddy-plugins/caddy-locale v0.0.2
	github.com/caddy-plugins/caddy-prometheus v0.1.0
	github.com/caddy-plugins/caddy-rate-limit v1.6.5
	github.com/caddy-plugins/caddy-s3browser v0.1.2
	github.com/caddy-plugins/cors v0.0.3
	github.com/caddy-plugins/ipfilter v1.1.6
	github.com/caddy-plugins/nobots v0.1.3
	github.com/caddy-plugins/webdav v1.2.10
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/chromedp/cdproto v0.0.0-20211107215322-d3760c2dc57b
	github.com/chromedp/chromedp v0.7.4
	github.com/codegangsta/inject v0.0.0-20150114235600-33e0aa1cb7c0 // indirect
	github.com/coscms/forms v1.10.0
	github.com/coscms/go-imgparse v0.0.0-20150925144422-3e3a099f7856
	github.com/dsoprea/go-logging v0.0.0-20200710184922-b02d349568dd // indirect
	github.com/edwingeng/doublejump v0.0.0-20210724020454-c82f1bcb3280 // indirect
	github.com/fatedier/beego v1.7.2 // indirect
	github.com/fatedier/golib v0.2.0
	github.com/fatih/color v1.13.0
	github.com/fd/go-shellwords v0.0.0-20130603174837-6a119423524d // indirect
	github.com/garyburd/redigo v1.6.3 // indirect
	github.com/go-errors/errors v1.4.1 // indirect
	github.com/go-ping/ping v0.0.0-20211014180314-6e2b003bffdd // indirect
	github.com/go-playground/validator/v10 v10.9.0 // indirect
	github.com/go-redis/redis/v8 v8.11.4 // indirect
	github.com/goccy/go-json v0.7.10 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/gomodule/redigo v1.8.5
	github.com/grafov/m3u8 v0.11.1 // indirect
	github.com/grandcat/zeroconf v1.0.0 // indirect
	github.com/hashicorp/consul/api v1.11.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-hclog v1.0.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/yamux v0.0.0-20211028200310-0bc27b27de87 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/kardianos/service v1.2.0 // indirect
	github.com/kisielk/errcheck v1.6.0 // indirect
	github.com/klauspost/pgzip v1.2.5 // indirect
	github.com/klauspost/reedsolomon v1.9.14 // indirect
	github.com/linxGnu/gumble v1.0.5 // indirect
	github.com/markbates/goth v1.68.0 // indirect
	github.com/maruel/rs v1.0.0 // indirect
	github.com/mattn/go-runewidth v0.0.13
	github.com/microcosm-cc/bluemonday v1.0.16
	github.com/miekg/dns v1.1.43
	github.com/minio/minio-go/v7 v7.0.15
	github.com/mrjones/oauth v0.0.0-20190623134757-126b35219450 // indirect
	github.com/nwaples/rardecode v1.1.2 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pkg/sftp v1.13.4
	github.com/rs/cors v1.8.0 // indirect
	github.com/russross/blackfriday v1.6.0
	github.com/shirou/gopsutil/v3 v3.21.10
	github.com/smallnest/quick v0.0.0-20210406061658-4bf95e372fbd // indirect
	github.com/smallnest/rpcx v1.6.11
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/syndtr/goleveldb v1.0.0
	github.com/tdewolff/minify v2.3.6+incompatible // indirect
	github.com/tdewolff/parse v2.3.4+incompatible // indirect
	github.com/tebeka/selenium v0.9.9
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/webx-top/captcha v0.0.1
	github.com/webx-top/chardet v0.0.1 // indirect
	github.com/webx-top/client v0.3.7
	github.com/webx-top/codec v0.1.1
	github.com/webx-top/com v0.3.2
	github.com/webx-top/db v1.17.0
	github.com/webx-top/echo v2.22.17+incompatible
	github.com/webx-top/echo-prometheus v1.0.10 // indirect
	github.com/webx-top/image v0.0.9
	github.com/webx-top/pagination v0.1.1
	gocloud.dev v0.24.0
	goftp.io/server/v2 v2.0.0
	golang.org/x/crypto v0.0.0-20211108221036-ceb1ce70b4fa
	golang.org/x/net v0.0.0-20211111160137-58aab5ef257a
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	google.golang.org/genproto v0.0.0-20211111162719-482062a4217b // indirect
	gopkg.in/ini.v1 v1.63.2 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

require (
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/admpub/copier v0.0.2 // indirect
	github.com/admpub/decimal v0.0.2 // indirect
	github.com/admpub/fsnotify v1.5.0 // indirect
	github.com/admpub/gifresize v1.0.2 // indirect
	github.com/admpub/go-reuseport v0.0.4 // indirect
	github.com/admpub/humanize v0.0.0-20190501023926-5f826e92c8ca // indirect
	github.com/admpub/json5 v0.0.1 // indirect
	github.com/admpub/map2struct v0.1.0 // indirect
	github.com/admpub/realip v0.0.0-20210421084339-374cf5df122d // indirect
	github.com/admpub/statik v0.1.7 // indirect
	github.com/admpub/webdav/v4 v4.1.2 // indirect
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/armon/go-socks5 v0.0.0-20160902184237-e75332964ef5 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/boombuler/barcode v1.0.1 // indirect
	github.com/cenk/backoff v2.2.1+incompatible // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cenkalti/backoff/v4 v4.1.2 // indirect
	github.com/cheekybits/genny v1.0.0 // indirect
	github.com/chromedp/sysutil v1.0.0 // indirect
	github.com/coreos/go-oidc v2.2.1+incompatible // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-jump v0.0.0-20211018200510-ba001c3ffce0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/disintegration/imaging v1.6.2 // indirect
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/facebookgo/clock v0.0.0-20150410010913-600d898af40a // indirect
	github.com/fatedier/kcp-go v2.0.4-0.20190803094908-fe8645b0a904+incompatible // indirect
	github.com/fcjr/aia-transport-go v1.2.2 // indirect
	github.com/flynn/go-shlex v0.0.0-20150515145356-3f9db97f8568 // indirect
	github.com/francoispqt/gojay v1.2.13 // indirect
	github.com/friendsofgo/errors v0.9.2 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-acme/lego/v4 v4.5.3 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/go-task/slim-sprig v0.0.0-20210107165309-348f09dbbbc0 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.1.0 // indirect
	github.com/gofrs/uuid v4.1.0+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/gax-go/v2 v2.1.1 // indirect
	github.com/gorilla/css v1.0.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-syslog v1.0.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/serf v0.9.5 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jimstudt/http-authentication v0.0.0-20140401203705-3eca13d6893a // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/juju/ratelimit v1.0.1 // indirect
	github.com/julienschmidt/httprouter v1.3.0 // indirect
	github.com/kavu/go_reuseport v1.5.0 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/klauspost/cpuid v1.3.1 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/lucas-clemente/quic-go v0.24.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/marten-seemann/qpack v0.2.1 // indirect
	github.com/marten-seemann/qtls-go1-16 v0.1.4 // indirect
	github.com/marten-seemann/qtls-go1-17 v0.1.0 // indirect
	github.com/mattn/go-colorable v0.1.11 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-sqlite3 v1.14.9 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.4.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/muesli/smartcrop v0.3.0 // indirect
	github.com/naoina/go-stringutil v0.1.0 // indirect
	github.com/naoina/toml v0.1.1 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/oschwald/maxminddb-golang v1.8.0 // indirect
	github.com/pires/go-proxyproto v0.6.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/pquerna/cachecontrol v0.1.0 // indirect
	github.com/prometheus/client_golang v1.11.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/rpcxio/libkv v0.5.1-0.20210420120011-1fceaedca8a5 // indirect
	github.com/rs/xid v1.3.0 // indirect
	github.com/rubyist/circuitbreaker v2.2.1+incompatible // indirect
	github.com/rwcarlsen/goexif v0.0.0-20190401172101-9e8deecbddbd // indirect
	github.com/samuel/go-zookeeper v0.0.0-20201211165307-7117e9ea2414 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/tdewolff/test v1.0.6 // indirect
	github.com/templexxx/cpufeat v0.0.0-20180724012125-cef66df7f161 // indirect
	github.com/templexxx/xor v0.0.0-20191217153810-f85b25db303b // indirect
	github.com/tjfoc/gmsm v1.4.1 // indirect
	github.com/tklauser/go-sysconf v0.3.9 // indirect
	github.com/tklauser/numcpus v0.3.0 // indirect
	github.com/tuotoo/qrcode v0.0.0-20190222102259-ac9c44189bf2 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fastrand v1.1.0 // indirect
	github.com/vmihailenco/msgpack/v5 v5.3.5 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	github.com/volatiletech/inflect v0.0.1 // indirect
	github.com/volatiletech/strmangle v0.0.1 // indirect
	github.com/webx-top/tagfast v0.0.0-20161020041435-9a2065ce3dd2 // indirect
	github.com/webx-top/validation v0.0.3 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	github.com/xtaci/kcp-go v5.4.20+incompatible // indirect
	go.opencensus.io v0.23.0 // indirect
	golang.org/x/image v0.0.0-20211028202545-6944b10bf410 // indirect
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/mod v0.5.1 // indirect
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8 // indirect
	golang.org/x/sys v0.0.0-20211111213525-f221eed1c01e // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	golang.org/x/tools v0.1.7 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/api v0.60.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/grpc v1.42.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
