module github.com/admpub/nging/v3

go 1.16

exclude github.com/gomodule/redigo v2.0.0+incompatible

require (
	gitee.com/admpub/certmagic v0.8.8
	github.com/PuerkitoBio/goquery v1.7.1 // indirect
	github.com/abh/errorutil v0.0.0-20130729183701-f9bd360d00b9
	github.com/admpub/9t v0.0.0-20190605154903-a68069ace5e1
	github.com/admpub/archiver v1.1.4
	github.com/admpub/bindata/v3 v3.1.5
	github.com/admpub/caddy v1.1.11
	github.com/admpub/ccs-gm v0.0.3
	github.com/admpub/checksum v1.0.1
	github.com/admpub/color v1.7.0
	github.com/admpub/confl v0.0.3
	github.com/admpub/cr v0.0.2
	github.com/admpub/cron v0.0.1
	github.com/admpub/dgoogauth v0.0.0-20170926052827-752650e076f2
	github.com/admpub/email v2.3.1+incompatible
	github.com/admpub/errors v0.8.2
	github.com/admpub/events v1.2.0
	github.com/admpub/fasthttp v0.0.2 // indirect
	github.com/admpub/frp v0.37.1
	github.com/admpub/go-bindata-assetfs v0.0.0-20170428090253-36eaa4c19588
	github.com/admpub/go-download/v2 v2.1.9
	github.com/admpub/go-figure v0.0.0-20180619031829-18b2b544842c
	github.com/admpub/go-isatty v0.0.9
	github.com/admpub/go-password v0.1.3
	github.com/admpub/go-phantomjs-fetcher v0.0.0-20180924162325-bb2ae1648e33
	github.com/admpub/go-pretty/v6 v6.0.3
	github.com/admpub/go-ps v0.0.1 // indirect
	github.com/admpub/go-sshclient v0.0.1
	github.com/admpub/godownloader v2.0.3+incompatible
	github.com/admpub/goforever v0.1.1
	github.com/admpub/gohls v0.0.0-20191013012052-b7505aaf3c90 // indirect
	github.com/admpub/gohls-server v0.3.4 // indirect
	github.com/admpub/gohttp v0.0.0-20190322032039-b55c707b8f1e
	github.com/admpub/gopiper v1.0.1
	github.com/admpub/goseaweedfs v0.1.2
	github.com/admpub/highwayhash v0.0.0-20180501080913-85fc8a2dacad
	github.com/admpub/httpscerts v0.0.0-20180907121630-a2990e2af45c
	github.com/admpub/i18n v0.0.2 // indirect
	github.com/admpub/imageproxy v0.9.3
	github.com/admpub/ip2region v1.2.6
	github.com/admpub/license_gen v0.0.0-20201028104329-fe31fcc255a8
	github.com/admpub/log v1.1.0
	github.com/admpub/logcool v0.3.2
	github.com/admpub/mahonia v0.0.0-20151019004008-c528b747d92d
	github.com/admpub/mail v0.0.0-20170408110349-d63147b0317b
	github.com/admpub/marmot v0.0.0-20200702042226-2170d9ff59f5
	github.com/admpub/metrohash v0.0.0-20160821164112-8d1c8b6bed28
	github.com/admpub/mysql-schema-sync v0.1.2
	github.com/admpub/null v8.0.3+incompatible
	github.com/admpub/pester v0.0.0-20200411024648-005672a2bd48 // indirect
	github.com/admpub/phantomjs v0.0.0-20180924162111-8a5af756140e
	github.com/admpub/qrcode v0.0.2
	github.com/admpub/randomize v0.0.2 // indirect
	github.com/admpub/redistore v1.1.0 // indirect
	github.com/admpub/regexp2 v1.1.7
	github.com/admpub/resty/v2 v2.6.0
	github.com/admpub/securecookie v1.1.2
	github.com/admpub/service v0.0.1
	github.com/admpub/sessions v0.1.1 // indirect
	github.com/admpub/snowflake v0.0.0-20180412010544-68117e6bbede
	github.com/admpub/sockjs-go/v3 v3.0.0
	github.com/admpub/sonyflake v0.0.1
	github.com/admpub/tail v1.0.3
	github.com/admpub/useragent v0.0.0-20190806155403-63e85649b0f2
	github.com/admpub/web-terminal v0.0.1
	github.com/admpub/websocket v1.0.4
	github.com/apache/thrift v0.14.2 // indirect
	github.com/arl/statsviz v0.4.0
	github.com/armon/go-metrics v0.3.9 // indirect
	github.com/aws/aws-sdk-go v1.40.22
	github.com/axgle/mahonia v0.0.0-20180208002826-3358181d7394 // indirect
	github.com/bitly/go-simplejson v0.5.0 // indirect
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/caddy-plugins/caddy-expires v1.1.2
	github.com/caddy-plugins/caddy-filter v0.15.2
	github.com/caddy-plugins/caddy-locale v0.0.2
	github.com/caddy-plugins/caddy-prometheus v0.0.3
	github.com/caddy-plugins/caddy-rate-limit v1.6.3
	github.com/caddy-plugins/caddy-s3browser v0.1.2
	github.com/caddy-plugins/cors v0.0.3
	github.com/caddy-plugins/ipfilter v1.1.6
	github.com/caddy-plugins/nobots v0.1.3
	github.com/caddy-plugins/webdav v1.2.10
	github.com/chromedp/cdproto v0.0.0-20210808225517-c36c1bd4c35e
	github.com/chromedp/chromedp v0.7.4
	github.com/codegangsta/inject v0.0.0-20150114235600-33e0aa1cb7c0 // indirect
	github.com/coscms/forms v1.9.2
	github.com/coscms/go-imgparse v0.0.0-20150925144422-3e3a099f7856
	github.com/edwingeng/doublejump v0.0.0-20210724020454-c82f1bcb3280 // indirect
	github.com/fatedier/beego v1.7.2 // indirect
	github.com/fatedier/golib v0.2.0
	github.com/fatih/color v1.12.0
	github.com/fd/go-shellwords v0.0.0-20130603174837-6a119423524d // indirect
	github.com/garyburd/redigo v1.6.2 // indirect
	github.com/go-ini/ini v1.62.0 // indirect
	github.com/go-ping/ping v0.0.0-20210506233800-ff8be3320020 // indirect
	github.com/go-playground/validator/v10 v10.9.0 // indirect
	github.com/go-redis/redis/v8 v8.11.3 // indirect
	github.com/golang/glog v0.0.0-20210429001901-424d2337a529 // indirect
	github.com/gomodule/redigo v1.8.5
	github.com/grafov/m3u8 v0.11.1 // indirect
	github.com/grandcat/zeroconf v1.0.0 // indirect
	github.com/hashicorp/consul/api v1.9.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-hclog v0.16.2 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/kardianos/service v1.2.0 // indirect
	github.com/kisielk/errcheck v1.6.0 // indirect
	github.com/klauspost/compress v1.13.4 // indirect
	github.com/klauspost/pgzip v1.2.5 // indirect
	github.com/klauspost/reedsolomon v1.9.13 // indirect
	github.com/linxGnu/gumble v1.0.5 // indirect
	github.com/lucas-clemente/quic-go v0.22.1 // indirect
	github.com/markbates/goth v1.68.0 // indirect
	github.com/maruel/rs v1.0.0 // indirect
	github.com/mattn/go-runewidth v0.0.13
	github.com/microcosm-cc/bluemonday v1.0.15
	github.com/minio/minio-go v6.0.14+incompatible
	github.com/mrjones/oauth v0.0.0-20190623134757-126b35219450 // indirect
	github.com/nwaples/rardecode v1.1.2 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pkg/sftp v1.13.2
	github.com/rs/cors v1.8.0 // indirect
	github.com/russross/blackfriday v1.6.0
	github.com/saintfish/chardet v0.0.0-20120816061221-3af4cd4741ca // indirect
	github.com/shirou/gopsutil/v3 v3.21.7
	github.com/shivakar/metrohash v0.0.0-20160821164112-8d1c8b6bed28 // indirect
	github.com/shivakar/xxhash v0.0.0-20160821164220-5ea66fb45566 // indirect
	github.com/smallnest/quick v0.0.0-20210406061658-4bf95e372fbd // indirect
	github.com/smallnest/rpcx v1.6.9
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/syndtr/goleveldb v1.0.0
	github.com/tdewolff/minify v2.3.6+incompatible // indirect
	github.com/tdewolff/parse v2.3.4+incompatible // indirect
	github.com/tdewolff/test v1.0.6 // indirect
	github.com/tebeka/selenium v0.9.9
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/webx-top/captcha v0.0.1
	github.com/webx-top/chardet v0.0.0-20180930194453-2f90d95f7b7f // indirect
	github.com/webx-top/client v0.3.5
	github.com/webx-top/codec v0.1.1
	github.com/webx-top/com v0.2.7
	github.com/webx-top/db v1.15.11
	github.com/webx-top/echo v2.19.4+incompatible
	github.com/webx-top/image v0.0.8
	github.com/webx-top/pagination v0.1.0
	gocloud.dev v0.23.0
	goftp.io/server/v2 v2.0.0
	golang.org/x/crypto v0.0.0-20210813211128-0a44fdfbc16e
	golang.org/x/mod v0.5.0 // indirect
	golang.org/x/net v0.0.0-20210813160813-60bc85c4be6d
	golang.org/x/oauth2 v0.0.0-20210810183815-faf39c7919d5 // indirect
	google.golang.org/genproto v0.0.0-20210813162853-db860fec028c // indirect
	google.golang.org/grpc v1.40.0 // indirect
	gopkg.in/fsnotify/fsnotify.v1 v1.4.7 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)
