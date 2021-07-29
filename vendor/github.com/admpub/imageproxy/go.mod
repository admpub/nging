module github.com/admpub/imageproxy

go 1.16

require (
	cloud.google.com/go v0.88.0 // indirect
	cloud.google.com/go/storage v1.16.0
	github.com/Azure/azure-sdk-for-go v56.0.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.19 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.14 // indirect
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/PaulARoy/azurestoragecache v0.0.0-20170906084534-3c249a3ba788
	github.com/admpub/gifresize v1.0.2
	github.com/aws/aws-sdk-go v1.40.10
	github.com/die-net/lrucache v0.0.0-20210724224853-653a274e85b0
	github.com/disintegration/imaging v1.6.2
	github.com/dnaeon/go-vcr v1.2.0 // indirect
	github.com/fcjr/aia-transport-go v1.2.2
	github.com/form3tech-oss/jwt-go v3.2.3+incompatible // indirect
	github.com/gofrs/uuid v4.0.0+incompatible // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/btree v1.0.1 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79
	github.com/jamiealquiza/envy v1.1.0
	github.com/muesli/smartcrop v0.3.0
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/peterbourgon/diskv v0.0.0-20171120014656-2973218375c3
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.30.0 // indirect
	github.com/prometheus/procfs v0.7.1 // indirect
	github.com/rwcarlsen/goexif v0.0.0-20190401172101-9e8deecbddbd
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d
	golang.org/x/net v0.0.0-20210726213435-c6fcb2dbf985 // indirect
	google.golang.org/api v0.51.0 // indirect
	google.golang.org/genproto v0.0.0-20210728212813-7823e685a01f // indirect
	willnorris.com/go/gifresize v1.0.0 // indirect
)

// local copy of envy package without cobra support
replace github.com/jamiealquiza/envy => ./third_party/envy
