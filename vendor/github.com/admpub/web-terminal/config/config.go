package config

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kardianos/osext"
)

var (
	Default = &Config{
		SHExecute: `bash`,
	}
	ExecutableFolder string
)

func init() {
	executableFolder, e := osext.ExecutableFolder()
	if nil != e {
		log.Println(e)
		if len(os.Args) > 0 {
			executableFolder = filepath.Dir(os.Args[0])
		}
	}
	ExecutableFolder = executableFolder
}

type Config struct {
	SHExecute   string
	Listen      string
	Debug       bool
	MIBSDir     string
	LogDir      string
	ResourceDir string
	APPRoot     string
	Password    string
	IDFile      string
	SHFile      string
}

func (c *Config) SetDefault() *Config {
	if len(c.LogDir) == 0 {
		c.LogDir = AutoLogDir()
	}
	if len(c.MIBSDir) == 0 {
		c.MIBSDir = AutoMIBSDir()
	}
	if len(c.ResourceDir) == 0 {
		c.ResourceDir = AutoResourceDir()
	}
	if !strings.HasSuffix(c.APPRoot, "/") {
		c.APPRoot = c.APPRoot + "/"
	}
	if !strings.HasPrefix(c.APPRoot, "/") {
		c.APPRoot = "/" + c.APPRoot
	}
	return c
}

func AbsPath(s string) string {
	r, e := filepath.Abs(s)
	if nil != e {
		return s
	}
	return r
}

func AutoLogDir() string {
	files := []string{"logs",
		filepath.Join("..", "logs"),
		filepath.Join(ExecutableFolder, "logs"),
		filepath.Join(ExecutableFolder, "..", "logs")}
	for _, nm := range files {
		nm = AbsPath(nm)
		if st, e := os.Stat(nm); nil == e && nil != st && st.IsDir() {
			nm = nm + "/"
			log.Println("'logs' directory is '" + nm + "'")
			return nm
		}
	}
	return ``
}

func AutoMIBSDir() string {
	files := []string{"mibs",
		filepath.Join("lib", "mibs"),
		filepath.Join("tools", "mibs"),
		filepath.Join("..", "lib", "mibs"),
		filepath.Join("..", "tools", "mibs"),
		filepath.Join(ExecutableFolder, "mibs"),
		filepath.Join(ExecutableFolder, "tools", "mibs"),
		filepath.Join(ExecutableFolder, "lib", "mibs"),
		filepath.Join(ExecutableFolder, "..", "lib", "mibs"),
		filepath.Join(ExecutableFolder, "..", "tools", "mibs")}
	for _, nm := range files {
		nm = AbsPath(nm)
		if st, e := os.Stat(nm); nil == e && nil != st && st.IsDir() {
			log.Println("'mibs' directory is '" + nm + "'")
			return nm
		}
	}
	return ``
}

func AutoResourceDir() string {
	files := []string{"web-terminal",
		filepath.Join("lib", "web-terminal"),
		filepath.Join("..", "lib", "web-terminal"),
		filepath.Join(ExecutableFolder, "static"),
		filepath.Join(ExecutableFolder, "web-terminal"),
		filepath.Join(ExecutableFolder, "lib", "web-terminal"),
		filepath.Join(ExecutableFolder, "..", "lib", "web-terminal")}

	for _, nm := range files {
		nm = AbsPath(nm)
		if st, e := os.Stat(nm); nil == e && nil != st && st.IsDir() {
			return nm
		}
	}
	buffer := bytes.NewBuffer(make([]byte, 0, 2048))
	buffer.WriteString("[warn] root path is not found:\r\n")
	for _, nm := range files {
		buffer.WriteString("\t\t")
		buffer.WriteString(nm)
		buffer.WriteString("\r\n")
	}
	buffer.Truncate(buffer.Len() - 2)
	log.Println(buffer)
	return ``
}
