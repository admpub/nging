module github.com/admpub/nging/v5

go 1.24.5

exclude github.com/gomodule/redigo v2.0.0+incompatible

replace github.com/fatedier/golib => github.com/fatedier/golib v0.2.0

// replace github.com/nging-plugins/caddymanager => ../../nging-plugins/caddymanager

// replace github.com/nging-plugins/collector => ../../nging-plugins/collector

// replace github.com/nging-plugins/dbmanager => ../../nging-plugins/dbmanager

// replace github.com/nging-plugins/ddnsmanager => ../../nging-plugins/ddnsmanager

// replace github.com/nging-plugins/dlmanager => ../../nging-plugins/dlmanager

// replace github.com/nging-plugins/frpmanager => ../../nging-plugins/frpmanager

// replace github.com/nging-plugins/ftpmanager => ../../nging-plugins/ftpmanager

// replace github.com/nging-plugins/servermanager => ../../nging-plugins/servermanager

// replace github.com/nging-plugins/sshmanager => ../../nging-plugins/sshmanager

// replace github.com/nging-plugins/webauthn => ../../nging-plugins/webauthn

// replace github.com/webx-top/client => ../../webx-top/client

// replace github.com/admpub/web-terminal => ../../admpub/web-terminal

// replace github.com/webx-top/echo => ../../webx-top/echo

// replace github.com/coscms/webcore => ../../coscms/webcore

require (
	github.com/admpub/copier v0.1.1
	github.com/admpub/go-ps v0.0.1
	github.com/admpub/regexp2 v1.1.8
	github.com/admpub/sse v0.0.0-20160126180136-ee05b128a739
	github.com/coscms/webcore v0.12.4-0.20250820043758-88039f05c360
	github.com/nging-plugins/caddymanager v1.8.18
	github.com/nging-plugins/collector v1.8.2
	github.com/nging-plugins/dbmanager v1.8.7
	github.com/nging-plugins/ddnsmanager v1.8.0
	github.com/nging-plugins/dlmanager v1.8.2
	github.com/nging-plugins/firewallmanager v1.8.3
	github.com/nging-plugins/frpmanager v1.8.4
	github.com/nging-plugins/ftpmanager v1.8.2
	github.com/nging-plugins/servermanager v1.8.10
	github.com/nging-plugins/sshmanager v1.8.7
	github.com/nging-plugins/webauthn v1.8.1
)

require (
	cloud.google.com/go/compute/metadata v0.7.0 // indirect
	github.com/BurntSushi/toml v1.5.0 // indirect
	github.com/GehirnInc/crypt v0.0.0-20230320061759-8cc1b52080c5 // indirect
	github.com/PuerkitoBio/goquery v1.10.3 // indirect
	github.com/VividCortex/ewma v1.2.0 // indirect
	github.com/abbot/go-http-auth v0.4.0 // indirect
	github.com/acarl005/stripansi v0.0.0-20180116102854-5a71ef0e047d // indirect
	github.com/admpub/bart v0.0.2 // indirect
	github.com/admpub/caddy v1.2.9 // indirect
	github.com/admpub/collate v1.1.0 // indirect
	github.com/admpub/conpty v0.2.0 // indirect
	github.com/admpub/cr v0.0.5 // indirect
	github.com/admpub/frp v0.37.8 // indirect
	github.com/admpub/gerberos v0.2.0 // indirect
	github.com/admpub/go-captcha-assets v0.0.0-20250122071745-baa7da4bda0d // indirect
	github.com/admpub/go-captcha/v2 v2.0.7 // indirect
	github.com/admpub/go-iptables v0.6.5 // indirect
	github.com/admpub/go-sshclient v0.0.3 // indirect
	github.com/admpub/go-ttlmap v1.1.0 // indirect
	github.com/admpub/godownloader v2.2.2+incompatible // indirect
	github.com/admpub/goforever v0.3.7 // indirect
	github.com/admpub/gohls v1.3.3 // indirect
	github.com/admpub/gohls-server v0.3.10 // indirect
	github.com/admpub/gopiper v1.1.2 // indirect
	github.com/admpub/gopty v0.1.2 // indirect
	github.com/admpub/map2struct v0.1.3 // indirect
	github.com/admpub/nftablesutils v0.3.4 // indirect
	github.com/admpub/osinfo v0.0.2 // indirect
	github.com/admpub/packer v0.0.3 // indirect
	github.com/admpub/sockjs-go/v3 v3.0.1 // indirect
	github.com/admpub/statik v0.1.7 // indirect
	github.com/admpub/useragent v0.0.2 // indirect
	github.com/admpub/web-terminal v0.2.1 // indirect
	github.com/admpub/webdav/v4 v4.1.3 // indirect
	github.com/admpub/xencoding v0.0.3 // indirect
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/armon/go-socks5 v0.0.0-20160902184237-e75332964ef5 // indirect
	github.com/bitly/go-simplejson v0.5.1 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/caddy-plugins/caddy-expires v1.1.3 // indirect
	github.com/caddy-plugins/caddy-filter v0.15.5 // indirect
	github.com/caddy-plugins/caddy-jwt/v3 v3.8.3 // indirect
	github.com/caddy-plugins/caddy-locale v0.0.2 // indirect
	github.com/caddy-plugins/caddy-prometheus v0.1.0 // indirect
	github.com/caddy-plugins/caddy-rate-limit v1.7.0 // indirect
	github.com/caddy-plugins/caddy-s3browser v0.2.2 // indirect
	github.com/caddy-plugins/cors v0.0.3 // indirect
	github.com/caddy-plugins/ipfilter v1.1.8 // indirect
	github.com/caddy-plugins/loginsrv v0.1.11 // indirect
	github.com/caddy-plugins/nobots v0.2.1 // indirect
	github.com/caddy-plugins/webdav v1.2.10 // indirect
	github.com/chromedp/cdproto v0.0.0-20250530212709-4dcc110a7b92 // indirect
	github.com/chromedp/chromedp v0.13.6 // indirect
	github.com/chromedp/sysutil v1.1.0 // indirect
	github.com/coreos/go-oidc/v3 v3.14.1 // indirect
	github.com/coscms/captcha v0.2.3 // indirect
	github.com/coscms/session-boltstore v0.0.0-20250617034717-a58d8848fa61 // indirect
	github.com/coscms/session-mysqlstore v0.0.0-20250617035706-a20b648443b1 // indirect
	github.com/coscms/session-redisstore v0.0.0-20250624032337-117cf04cfaf4 // indirect
	github.com/coscms/session-sqlitestore v0.0.4 // indirect
	github.com/coscms/session-sqlstore v0.0.1 // indirect
	github.com/coscms/webauthn v0.3.2 // indirect
	github.com/creack/pty v1.1.24 // indirect
	github.com/ebitengine/purego v0.8.4 // indirect
	github.com/fatedier/beego v0.0.0-20171024143340-6c6a4f5bd5eb // indirect
	github.com/fatedier/golib v0.5.1 // indirect
	github.com/fatedier/kcp-go v2.0.4-0.20190803094908-fe8645b0a904+incompatible // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/fd/go-shellwords v0.0.0-20130603174837-6a119423524d // indirect
	github.com/flynn/go-shlex v0.0.0-20150515145356-3f9db97f8568 // indirect
	github.com/fxamacker/cbor/v2 v2.8.0 // indirect
	github.com/gaissmai/extnetip v1.1.0 // indirect
	github.com/go-json-experiment/json v0.0.0-20250517221953-25912455fbc8 // indirect
	github.com/go-webauthn/webauthn v0.13.0 // indirect
	github.com/go-webauthn/x v0.1.21 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.4.0 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/go-tpm v0.9.5 // indirect
	github.com/google/nftables v0.3.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/gorilla/sessions v1.4.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/grafov/m3u8 v0.12.1 // indirect
	github.com/hashicorp/go-syslog v1.0.0 // indirect
	github.com/hashicorp/yamux v0.1.2 // indirect
	github.com/iamacarpet/go-winpty v1.0.4 // indirect
	github.com/jimstudt/http-authentication v0.0.0-20140401203705-3eca13d6893a // indirect
	github.com/klauspost/reedsolomon v1.12.4 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/mdlayher/netlink v1.7.3-0.20250113171957-fbb4dce95f42 // indirect
	github.com/mdlayher/socket v0.5.1 // indirect
	github.com/minio/crc64nvme v1.1.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/naoina/go-stringutil v0.1.0 // indirect
	github.com/naoina/toml v0.1.1 // indirect
	github.com/oschwald/maxminddb-golang v1.13.1 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/pires/go-proxyproto v0.8.1 // indirect
	github.com/pkg/sftp v1.13.9 // indirect
	github.com/quic-go/qpack v0.5.1 // indirect
	github.com/quic-go/quic-go v0.54.0 // indirect
	github.com/russross/blackfriday v1.6.0 // indirect
	github.com/shirou/gopsutil/v4 v4.25.7 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/tarent/lib-compose/v2 v2.0.1 // indirect
	github.com/tebeka/selenium v0.9.9 // indirect
	github.com/templexxx/cpufeat v0.0.0-20180724012125-cef66df7f161 // indirect
	github.com/templexxx/xor v0.0.0-20191217153810-f85b25db303b // indirect
	github.com/tinylib/msgp v1.3.0 // indirect
	github.com/tjfoc/gmsm v1.4.1 // indirect
	github.com/vbauerster/mpb/v6 v6.0.4 // indirect
	github.com/vishvananda/netlink v1.3.1 // indirect
	github.com/vishvananda/netns v0.0.5 // indirect
	github.com/webx-top/echo-prometheus v1.1.2 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	go.uber.org/mock v0.5.2 // indirect
	goftp.io/server/v2 v2.0.1 // indirect
	golang.org/x/time v0.12.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	gitee.com/admpub/certmagic v0.8.9 // indirect
	github.com/abh/errorutil v1.0.0 // indirect
	github.com/admpub/9t v0.0.1 // indirect
	github.com/admpub/bindata/v3 v3.2.1
	github.com/admpub/boltstore v1.2.0 // indirect
	github.com/admpub/captcha-go v0.0.1 // indirect
	github.com/admpub/ccs-gm v0.0.5 // indirect
	github.com/admpub/checksum v1.1.0
	github.com/admpub/color v1.8.1
	github.com/admpub/confl v0.2.4 // indirect
	github.com/admpub/cron v0.1.1 // indirect
	github.com/admpub/decimal v1.3.2 // indirect
	github.com/admpub/dgoogauth v0.0.1
	github.com/admpub/email v2.4.1+incompatible // indirect
	github.com/admpub/errors v0.8.2
	github.com/admpub/events v1.3.6
	github.com/admpub/fasthttp v0.0.7 // indirect
	github.com/admpub/fsnotify v1.7.1 // indirect
	github.com/admpub/gifresize v1.0.2 // indirect
	github.com/admpub/go-bindata-assetfs v0.0.1 // indirect
	github.com/admpub/go-download/v2 v2.2.0 // indirect
	github.com/admpub/go-figure v0.0.2 // indirect
	github.com/admpub/go-isatty v0.0.11 // indirect
	github.com/admpub/go-password v0.1.3 // indirect
	github.com/admpub/go-pretty/v6 v6.0.4 // indirect
	github.com/admpub/go-reuseport v0.5.0 // indirect
	github.com/admpub/go-utility v0.0.1 // indirect
	github.com/admpub/godotenv v1.4.3 // indirect
	github.com/admpub/goth v0.0.4 // indirect
	github.com/admpub/httpscerts v0.0.0-20180907121630-a2990e2af45c // indirect
	github.com/admpub/humanize v0.0.0-20190501023926-5f826e92c8ca // indirect
	github.com/admpub/i18n v0.5.2 // indirect
	github.com/admpub/identicon v1.0.2 // indirect
	github.com/admpub/imageproxy v0.10.1
	github.com/admpub/imaging v1.6.3 // indirect
	github.com/admpub/ip2region/v2 v2.0.1 // indirect
	github.com/admpub/json5 v0.0.1 // indirect
	github.com/admpub/license_gen v0.1.1 // indirect
	github.com/admpub/log v1.4.0
	github.com/admpub/mahonia v0.0.0-20151019004008-c528b747d92d // indirect
	github.com/admpub/mail v0.0.0-20170408110349-d63147b0317b // indirect
	github.com/admpub/marmot v0.0.0-20200702042226-2170d9ff59f5 // indirect
	github.com/admpub/mysql-schema-sync v0.2.6 // indirect
	github.com/admpub/null v8.0.5+incompatible
	github.com/admpub/oauth2/v4 v4.0.3 // indirect
	github.com/admpub/once v0.0.1 // indirect
	github.com/admpub/pester v0.0.0-20200411024648-005672a2bd48 // indirect
	github.com/admpub/pp v0.0.7 // indirect
	github.com/admpub/qrcode v0.0.3 // indirect
	github.com/admpub/realip v0.2.7 // indirect
	github.com/admpub/redistore v1.2.2 // indirect
	github.com/admpub/resty/v2 v2.7.2 // indirect
	github.com/admpub/safesvg v0.0.8 // indirect
	github.com/admpub/securecookie v1.3.0 // indirect
	github.com/admpub/service v0.0.8 // indirect
	github.com/admpub/sessions v0.3.0 // indirect
	github.com/admpub/sonyflake v0.0.1 // indirect
	github.com/admpub/tail v1.1.1 // indirect
	github.com/admpub/timeago v1.3.0 // indirect
	github.com/admpub/websocket v1.0.4
	github.com/alicebob/gopher-json v0.0.0-20230218143504-906a9b012302 // indirect
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/aws/aws-sdk-go v1.55.8 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/boombuler/barcode v1.1.0 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/coscms/forms v1.13.10 // indirect
	github.com/coscms/go-imgparse v0.0.1 // indirect
	github.com/coscms/oauth2s v0.4.3 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dsoprea/go-logging v0.0.0-20200710184922-b02d349568dd // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fcjr/aia-transport-go v1.3.0 // indirect
	github.com/francoispqt/gojay v1.2.13 // indirect
	github.com/fynelabs/selfupdate v0.2.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.9 // indirect
	github.com/geoffgarside/ber v1.2.0 // indirect
	github.com/glebarez/go-sqlite v1.22.0 // indirect
	github.com/go-acme/lego/v4 v4.25.2 // indirect
	github.com/go-errors/errors v1.5.1 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-jose/go-jose/v4 v4.1.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.27.0 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/gomodule/redigo v1.9.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/css v1.0.1 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/h2non/filetype v1.1.3 // indirect
	github.com/h2non/go-is-svg v0.0.0-20160927212452-35e8c4b0612c // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-version v1.7.0 // indirect
	github.com/hirochachacha/go-smb2 v1.1.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jlaffaye/ftp v0.2.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/kisielk/errcheck v1.9.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/klauspost/cpuid v1.3.1 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20250317134145-8bc96cf8fc35 // indirect
	github.com/maruel/rs v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mattn/go-sqlite3 v1.14.32 // indirect
	github.com/microcosm-cc/bluemonday v1.0.27 // indirect
	github.com/miekg/dns v1.1.68 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/minio/minio-go/v7 v7.0.95 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/muesli/smartcrop v0.3.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55 // indirect
	github.com/prometheus/client_golang v1.23.0 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.65.0 // indirect
	github.com/prometheus/procfs v0.17.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/rs/xid v1.6.0 // indirect
	github.com/rwcarlsen/goexif v0.0.0-20190401172101-9e8deecbddbd // indirect
	github.com/segmentio/fasthash v1.0.3 // indirect
	github.com/spf13/cobra v1.9.1 // indirect
	github.com/spf13/pflag v1.0.7 // indirect
	github.com/stretchr/testify v1.10.0
	github.com/studio-b12/gowebdav v0.10.0 // indirect
	github.com/syndtr/goleveldb v1.0.0
	github.com/tidwall/btree v1.8.1 // indirect
	github.com/tidwall/buntdb v1.3.2 // indirect
	github.com/tidwall/gjson v1.18.0 // indirect
	github.com/tidwall/grect v0.1.4 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/rtred v0.1.2 // indirect
	github.com/tidwall/tinyqueue v0.1.1 // indirect
	github.com/tklauser/go-sysconf v0.3.15 // indirect
	github.com/tklauser/numcpus v0.10.0 // indirect
	github.com/tuotoo/qrcode v0.0.0-20220425170535-52ccc2bebf5d // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/webx-top/captcha v0.1.0 // indirect
	github.com/webx-top/chardet v0.0.2 // indirect
	github.com/webx-top/client v0.9.6
	github.com/webx-top/codec v0.3.0 // indirect
	github.com/webx-top/com v1.4.0
	github.com/webx-top/db v1.28.6
	github.com/webx-top/echo v1.21.0
	github.com/webx-top/image v0.1.2
	github.com/webx-top/pagination v0.3.2 // indirect
	github.com/webx-top/poolx v0.0.0-20210912044716-5cfa2d58e380 // indirect
	github.com/webx-top/restyclient v0.0.6 // indirect
	github.com/webx-top/tagfast v0.0.1 // indirect
	github.com/webx-top/validation v0.0.3 // indirect
	github.com/webx-top/validator v0.3.0 // indirect
	github.com/yuin/gopher-lua v1.1.1 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.etcd.io/bbolt v1.4.2 // indirect
	golang.org/x/crypto v0.41.0 // indirect
	golang.org/x/exp v0.0.0-20250813145105-42675adae3e6 // indirect
	golang.org/x/image v0.30.0 // indirect
	golang.org/x/lint v0.0.0-20241112194109-818c5a804067 // indirect
	golang.org/x/mod v0.27.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/oauth2 v0.30.0 // indirect
	golang.org/x/sync v0.16.0
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	golang.org/x/tools v0.36.0 // indirect
	google.golang.org/protobuf v1.36.7 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	modernc.org/libc v1.66.7 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.11.0 // indirect
	modernc.org/sqlite v1.38.2 // indirect
	rsc.io/qr v0.2.0 // indirect
)
