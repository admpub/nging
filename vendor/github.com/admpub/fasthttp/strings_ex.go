package fasthttp

func init() {
	defaultServerName = []byte("webx")
	defaultUserAgent = []byte("webx")
	//defaultContentType = []byte("text/html; charset=utf-8")
}

func SetDefaultServerName(name []byte) {
	defaultServerName = name
}

func SetDefaultUserAgent(userAgent []byte) {
	defaultUserAgent = userAgent
}

func SetDefaultContentType(contentType []byte) {
	defaultContentType = contentType
}
