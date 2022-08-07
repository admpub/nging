package buildinfo

import "github.com/webx-top/echo"

func New(opts ...Option) *BuildInfo {
	b := &BuildInfo{}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

type BuildInfo struct {
	BUILD_TIME string
	BUILD_OS   string
	BUILD_ARCH string
	CLOUD_GOX  string
	COMMIT     string
	LABEL      string // beta/alpha/stable
	VERSION    string // `4.1.6`
	PACKAGE    string // `free`
	SCHEMA_VER float64
}

func (b *BuildInfo) Apply() {
	echo.Set(`BUILD_TIME`, b.BUILD_TIME)
	echo.Set(`BUILD_OS`, b.BUILD_OS)
	echo.Set(`BUILD_ARCH`, b.BUILD_ARCH)
	echo.Set(`COMMIT`, b.COMMIT)
	echo.Set(`LABEL`, b.LABEL)
	echo.Set(`VERSION`, b.VERSION)
	echo.Set(`PACKAGE`, b.PACKAGE)
	echo.Set(`SCHEMA_VER`, b.SCHEMA_VER)
}

type Option func(*BuildInfo)

func Time(buildTime string) Option {
	return func(b *BuildInfo) {
		b.BUILD_TIME = buildTime
	}
}

func OS(buildOS string) Option {
	return func(b *BuildInfo) {
		b.BUILD_OS = buildOS
	}
}

func Arch(buildArch string) Option {
	return func(b *BuildInfo) {
		b.BUILD_ARCH = buildArch
	}
}

func Commit(commit string) Option {
	return func(b *BuildInfo) {
		b.COMMIT = commit
	}
}

func Label(label string) Option {
	return func(b *BuildInfo) {
		b.LABEL = label
	}
}

func Version(version string) Option {
	return func(b *BuildInfo) {
		b.VERSION = version
	}
}

func Package(pkg string) Option {
	return func(b *BuildInfo) {
		b.PACKAGE = pkg
	}
}

func SchemaVer(schemaVer float64) Option {
	return func(b *BuildInfo) {
		b.SCHEMA_VER = schemaVer
	}
}

func CloudGox(cloudGox string) Option {
	return func(b *BuildInfo) {
		b.CLOUD_GOX = cloudGox
	}
}
