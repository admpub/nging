module github.com/admpub/nging/v5

go 1.23

toolchain go1.23.0

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
	github.com/coscms/webcore v0.2.3
	github.com/nging-plugins/caddymanager v1.6.0
	github.com/nging-plugins/collector v1.6.0
	github.com/nging-plugins/dbmanager v1.6.3
	github.com/nging-plugins/ddnsmanager v1.6.0
	github.com/nging-plugins/dlmanager v1.6.2
	github.com/nging-plugins/firewallmanager v1.6.1
	github.com/nging-plugins/frpmanager v1.6.0
	github.com/nging-plugins/ftpmanager v1.6.0
	github.com/nging-plugins/servermanager v1.6.1
	github.com/nging-plugins/sshmanager v1.6.0
	github.com/nging-plugins/webauthn v1.6.0
)

require (
	cloud.google.com/go/compute/metadata v0.5.1 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	gitee.com/admpub/certmagic v0.8.8 // indirect
	github.com/BurntSushi/toml v1.4.0 // indirect
	github.com/GehirnInc/crypt v0.0.0-20230320061759-8cc1b52080c5 // indirect
	github.com/PuerkitoBio/goquery v1.10.0 // indirect
	github.com/abbot/go-http-auth v0.4.0 // indirect
	github.com/abh/errorutil v1.0.0 // indirect
	github.com/admpub/9t v0.0.0-20190605154903-a68069ace5e1 // indirect
	github.com/admpub/archiver v1.1.4 // indirect
	github.com/admpub/bindata/v3 v3.2.1
	github.com/admpub/boltstore v1.1.2 // indirect
	github.com/admpub/caddy v1.2.7 // indirect
	github.com/admpub/captcha-go v0.0.1 // indirect
	github.com/admpub/ccs-gm v0.0.5 // indirect
	github.com/admpub/checksum v1.1.0
	github.com/admpub/collate v1.1.0 // indirect
	github.com/admpub/color v1.8.1
	github.com/admpub/confl v0.2.4 // indirect
	github.com/admpub/conpty v0.2.0 // indirect
	github.com/admpub/cr v0.0.5 // indirect
	github.com/admpub/cron v0.1.1 // indirect
	github.com/admpub/decimal v1.3.1 // indirect
	github.com/admpub/dgoogauth v0.0.1
	github.com/admpub/email v2.4.1+incompatible // indirect
	github.com/admpub/errors v0.8.2
	github.com/admpub/events v1.3.6
	github.com/admpub/fasthttp v0.0.5 // indirect
	github.com/admpub/frp v0.37.7 // indirect
	github.com/admpub/fsnotify v1.7.0 // indirect
	github.com/admpub/gerberos v0.1.1 // indirect
	github.com/admpub/gifresize v1.0.2 // indirect
	github.com/admpub/go-bindata-assetfs v0.0.0-20170428090253-36eaa4c19588 // indirect
	github.com/admpub/go-download/v2 v2.1.15 // indirect
	github.com/admpub/go-figure v0.0.0-20180619031829-18b2b544842c // indirect
	github.com/admpub/go-iptables v0.6.5 // indirect
	github.com/admpub/go-isatty v0.0.11 // indirect
	github.com/admpub/go-password v0.1.3 // indirect
	github.com/admpub/go-pretty/v6 v6.0.4 // indirect
	github.com/admpub/go-reuseport v0.0.4 // indirect
	github.com/admpub/go-sshclient v0.0.3 // indirect
	github.com/admpub/go-ttlmap v1.1.0 // indirect
	github.com/admpub/go-utility v0.0.1 // indirect
	github.com/admpub/godotenv v1.4.3 // indirect
	github.com/admpub/godownloader v2.2.2+incompatible // indirect
	github.com/admpub/goforever v0.3.5 // indirect
	github.com/admpub/gohls v1.3.3 // indirect
	github.com/admpub/gohls-server v0.3.10 // indirect
	github.com/admpub/gopiper v1.1.2 // indirect
	github.com/admpub/gopty v0.1.2 // indirect
	github.com/admpub/goth v0.0.4 // indirect
	github.com/admpub/httpscerts v0.0.0-20180907121630-a2990e2af45c // indirect
	github.com/admpub/humanize v0.0.0-20190501023926-5f826e92c8ca // indirect
	github.com/admpub/i18n v0.3.2 // indirect
	github.com/admpub/identicon v1.0.2 // indirect
	github.com/admpub/imageproxy v0.10.0
	github.com/admpub/imaging v1.6.3 // indirect
	github.com/admpub/ip2region/v2 v2.0.1 // indirect
	github.com/admpub/json5 v0.0.1 // indirect
	github.com/admpub/license_gen v0.1.1 // indirect
	github.com/admpub/log v1.3.6
	github.com/admpub/mahonia v0.0.0-20151019004008-c528b747d92d // indirect
	github.com/admpub/mail v0.0.0-20170408110349-d63147b0317b // indirect
	github.com/admpub/marmot v0.0.0-20200702042226-2170d9ff59f5 // indirect
	github.com/admpub/mysql-schema-sync v0.2.6 // indirect
	github.com/admpub/nftablesutils v0.3.4 // indirect
	github.com/admpub/null v8.0.4+incompatible
	github.com/admpub/oauth2/v4 v4.0.2 // indirect
	github.com/admpub/once v0.0.1 // indirect
	github.com/admpub/osinfo v0.0.2 // indirect
	github.com/admpub/packer v0.0.3 // indirect
	github.com/admpub/pester v0.0.0-20200411024648-005672a2bd48 // indirect
	github.com/admpub/pp v0.0.7 // indirect
	github.com/admpub/qrcode v0.0.3 // indirect
	github.com/admpub/randomize v0.0.2 // indirect
	github.com/admpub/realip v0.2.7 // indirect
	github.com/admpub/redistore v1.2.1 // indirect
	github.com/admpub/resty/v2 v2.7.1 // indirect
	github.com/admpub/safesvg v0.0.8 // indirect
	github.com/admpub/securecookie v1.3.0 // indirect
	github.com/admpub/service v0.0.5 // indirect
	github.com/admpub/sessions v0.2.1 // indirect
	github.com/admpub/sockjs-go/v3 v3.0.1 // indirect
	github.com/admpub/sonyflake v0.0.1 // indirect
	github.com/admpub/statik v0.1.7 // indirect
	github.com/admpub/tail v1.1.0 // indirect
	github.com/admpub/timeago v1.2.1 // indirect
	github.com/admpub/useragent v0.0.2 // indirect
	github.com/admpub/web-terminal v0.2.1 // indirect
	github.com/admpub/webdav/v4 v4.1.3 // indirect
	github.com/admpub/websocket v1.0.4
	github.com/alicebob/gopher-json v0.0.0-20230218143504-906a9b012302 // indirect
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/andybalholm/cascadia v1.3.2 // indirect
	github.com/armon/go-socks5 v0.0.0-20160902184237-e75332964ef5 // indirect
	github.com/aws/aws-sdk-go v1.55.5 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bitly/go-simplejson v0.5.1 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/boombuler/barcode v1.0.2 // indirect
	github.com/caddy-plugins/caddy-expires v1.1.3 // indirect
	github.com/caddy-plugins/caddy-filter v0.15.5 // indirect
	github.com/caddy-plugins/caddy-jwt/v3 v3.8.2 // indirect
	github.com/caddy-plugins/caddy-locale v0.0.2 // indirect
	github.com/caddy-plugins/caddy-prometheus v0.1.0 // indirect
	github.com/caddy-plugins/caddy-rate-limit v1.7.0 // indirect
	github.com/caddy-plugins/caddy-s3browser v0.1.3 // indirect
	github.com/caddy-plugins/cors v0.0.3 // indirect
	github.com/caddy-plugins/ipfilter v1.1.8 // indirect
	github.com/caddy-plugins/loginsrv v0.1.9 // indirect
	github.com/caddy-plugins/nobots v0.2.1 // indirect
	github.com/caddy-plugins/webdav v1.2.10 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/chromedp/cdproto v0.0.0-20240919203636-12af5e8a671f // indirect
	github.com/chromedp/chromedp v0.10.0 // indirect
	github.com/chromedp/sysutil v1.0.0 // indirect
	github.com/coreos/go-oidc/v3 v3.11.0 // indirect
	github.com/coscms/forms v1.12.2 // indirect
	github.com/coscms/go-imgparse v0.0.1 // indirect
	github.com/coscms/oauth2s v0.4.1 // indirect
	github.com/coscms/webauthn v0.3.1 // indirect
	github.com/creack/pty v1.1.23 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/dsoprea/go-logging v0.0.0-20200710184922-b02d349568dd // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fatedier/beego v0.0.0-20171024143340-6c6a4f5bd5eb // indirect
	github.com/fatedier/golib v0.5.0 // indirect
	github.com/fatedier/kcp-go v2.0.4-0.20190803094908-fe8645b0a904+incompatible // indirect
	github.com/fatih/color v1.17.0 // indirect
	github.com/fcjr/aia-transport-go v1.2.2 // indirect
	github.com/fd/go-shellwords v0.0.0-20130603174837-6a119423524d // indirect
	github.com/flynn/go-shlex v0.0.0-20150515145356-3f9db97f8568 // indirect
	github.com/francoispqt/gojay v1.2.13 // indirect
	github.com/friendsofgo/errors v0.9.2 // indirect
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/fynelabs/selfupdate v0.2.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.5 // indirect
	github.com/gaissmai/extnetip v1.1.0 // indirect
	github.com/garyburd/redigo v1.6.4 // indirect
	github.com/geoffgarside/ber v1.1.0 // indirect
	github.com/glebarez/go-sqlite v1.22.0 // indirect
	github.com/go-acme/lego/v4 v4.18.0 // indirect
	github.com/go-errors/errors v1.5.1 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-jose/go-jose/v4 v4.0.4 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.22.1 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/go-webauthn/webauthn v0.11.2 // indirect
	github.com/go-webauthn/x v0.1.14 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.4.0 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gomodule/redigo v1.9.2 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/go-tpm v0.9.1 // indirect
	github.com/google/nftables v0.2.0 // indirect
	github.com/google/pprof v0.0.0-20240910150728-a0b0bb1d4134 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/css v1.0.1 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/gorilla/sessions v1.4.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/grafov/m3u8 v0.12.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/h2non/filetype v1.1.3 // indirect
	github.com/h2non/go-is-svg v0.0.0-20160927212452-35e8c4b0612c // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-syslog v1.0.0 // indirect
	github.com/hashicorp/go-version v1.7.0 // indirect
	github.com/hashicorp/yamux v0.1.1 // indirect
	github.com/hirochachacha/go-smb2 v1.1.0 // indirect
	github.com/iamacarpet/go-winpty v1.0.4 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jimstudt/http-authentication v0.0.0-20140401203705-3eca13d6893a // indirect
	github.com/jlaffaye/ftp v0.2.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/josharian/native v1.1.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/kisielk/errcheck v1.7.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/klauspost/cpuid v1.3.1 // indirect
	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
	github.com/klauspost/pgzip v1.2.6 // indirect
	github.com/klauspost/reedsolomon v1.12.4 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20240909124753-873cd0166683 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/maruel/rs v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mattn/go-sqlite3 v1.14.23 // indirect
	github.com/mdlayher/netlink v1.7.2 // indirect
	github.com/mdlayher/socket v0.5.1 // indirect
	github.com/microcosm-cc/bluemonday v1.0.27 // indirect
	github.com/miekg/dns v1.1.62 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/minio/minio-go/v7 v7.0.76 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/muesli/smartcrop v0.3.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/naoina/go-stringutil v0.1.0 // indirect
	github.com/naoina/toml v0.1.1 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/nwaples/rardecode v1.1.3 // indirect
	github.com/onsi/ginkgo/v2 v2.20.2 // indirect
	github.com/oschwald/maxminddb-golang v1.13.1 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pires/go-proxyproto v0.7.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pkg/sftp v1.13.6 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55 // indirect
	github.com/prometheus/client_golang v1.20.4 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.59.1 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/quic-go/qpack v0.5.1 // indirect
	github.com/quic-go/quic-go v0.47.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/rs/xid v1.6.0 // indirect
	github.com/russross/blackfriday v1.6.0 // indirect
	github.com/rwcarlsen/goexif v0.0.0-20190401172101-9e8deecbddbd // indirect
	github.com/segmentio/fasthash v1.0.3 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible // indirect
	github.com/shirou/gopsutil/v3 v3.24.5 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/spf13/cobra v1.8.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.9.0
	github.com/studio-b12/gowebdav v0.9.0 // indirect
	github.com/syndtr/goleveldb v1.0.0
	github.com/tarent/lib-compose/v2 v2.0.1 // indirect
	github.com/tebeka/selenium v0.9.9 // indirect
	github.com/templexxx/cpufeat v0.0.0-20180724012125-cef66df7f161 // indirect
	github.com/templexxx/xor v0.0.0-20191217153810-f85b25db303b // indirect
	github.com/tidwall/btree v1.7.0 // indirect
	github.com/tidwall/buntdb v1.3.2 // indirect
	github.com/tidwall/gjson v1.17.3 // indirect
	github.com/tidwall/grect v0.1.4 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/rtred v0.1.2 // indirect
	github.com/tidwall/tinyqueue v0.1.1 // indirect
	github.com/tjfoc/gmsm v1.4.1 // indirect
	github.com/tklauser/go-sysconf v0.3.14 // indirect
	github.com/tklauser/numcpus v0.8.0 // indirect
	github.com/tuotoo/qrcode v0.0.0-20220425170535-52ccc2bebf5d // indirect
	github.com/ulikunitz/xz v0.5.12 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/vishvananda/netlink v1.3.0 // indirect
	github.com/vishvananda/netns v0.0.4 // indirect
	github.com/volatiletech/inflect v0.0.1 // indirect
	github.com/volatiletech/strmangle v0.0.6 // indirect
	github.com/webx-top/captcha v0.1.0 // indirect
	github.com/webx-top/chardet v0.0.2 // indirect
	github.com/webx-top/client v0.9.5
	github.com/webx-top/codec v0.3.0 // indirect
	github.com/webx-top/com v1.3.3
	github.com/webx-top/db v1.27.10
	github.com/webx-top/echo v1.10.3
	github.com/webx-top/echo-prometheus v1.1.2 // indirect
	github.com/webx-top/image v0.1.1
	github.com/webx-top/pagination v0.3.1 // indirect
	github.com/webx-top/poolx v0.0.0-20210912044716-5cfa2d58e380 // indirect
	github.com/webx-top/restyclient v0.0.4 // indirect
	github.com/webx-top/tagfast v0.0.1 // indirect
	github.com/webx-top/validation v0.0.3 // indirect
	github.com/webx-top/validator v0.3.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	github.com/yuin/gopher-lua v1.1.1 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.etcd.io/bbolt v1.3.11 // indirect
	go.uber.org/mock v0.4.0 // indirect
	goftp.io/server/v2 v2.0.1 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/exp v0.0.0-20240909161429-701f63a606c0 // indirect
	golang.org/x/image v0.20.0 // indirect
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/mod v0.21.0 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/oauth2 v0.23.0 // indirect
	golang.org/x/sync v0.8.0
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	golang.org/x/time v0.6.0 // indirect
	golang.org/x/tools v0.25.0 // indirect
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	modernc.org/libc v1.61.0 // indirect
	modernc.org/mathutil v1.6.0 // indirect
	modernc.org/memory v1.8.0 // indirect
	modernc.org/sqlite v1.33.1 // indirect
	rsc.io/qr v0.2.0 // indirect
)
