package utils

import (
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/admpub/web-terminal/config"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var CharsetList = map[string]encoding.Encoding{
	`GB18030`:       simplifiedchinese.GB18030,
	`GB2312`:        simplifiedchinese.HZGB2312,
	`HZ-GB2312`:     simplifiedchinese.HZGB2312,
	`GBK`:           simplifiedchinese.GBK,
	`BIG5`:          traditionalchinese.Big5,
	`EUC-JP`:        japanese.EUCJP,
	`ISO2022JP`:     japanese.ISO2022JP,
	`SHIFTJIS`:      japanese.ShiftJIS,
	`EUC-KR`:        korean.EUCKR,
	`UTF-8`:         encoding.Nop,
	`UTF-16`:        unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM),
	`UTF-16-BE`:     unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM),
	`UTF-16-LE`:     unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM),
	`UTF-16-BOM`:    unicode.UTF16(unicode.BigEndian, unicode.UseBOM),
	`UTF-16-BE-BOM`: unicode.UTF16(unicode.BigEndian, unicode.UseBOM),
	`UTF-16-LE-BOM`: unicode.UTF16(unicode.LittleEndian, unicode.UseBOM),
}

func SupportedCharsets() []string {
	r := []string{}
	for k := range CharsetList {
		r = append(r, k)
	}
	sort.Strings(r)
	return r
}

func CharsetEncoding(charset string) encoding.Encoding {
	charset = strings.ToUpper(charset)
	switch charset {
	case "UTF8":
		charset = "UTF-8"
	case "UTF16-BOM":
		charset = "UTF-16-BOM"
	case "UTF16-BE-BOM":
		charset = "UTF-16-BE-BOM"
	case "UTF16-LE-BOM":
		charset = "UTF-16-LE-BOM"
	case "UTF16":
		charset = "UTF-16"
	case "UTF16-BE":
		charset = "UTF-16-BE"
	case "UTF16-LE":
		charset = "UTF-16-LE"
	case "UTF32":
		charset = "UTF-32"
	}
	if enc, ok := CharsetList[charset]; ok {
		return enc
	}
	return nil
}

func Warp(dst io.ReadCloser, dump io.Writer) io.ReadCloser {
	if nil == dump {
		return dst
	}
	return &ConsoleReader{out: dump, dst: dst}
}

func DecodeBy(charset string, dst io.Writer) io.Writer {
	switch strings.ToUpper(charset) {
	case "UTF-8", "UTF8":
		return dst
	}
	cs := CharsetEncoding(charset)
	if nil == cs {
		panic("charset '" + charset + "' is not exists.")
	}

	return transform.NewWriter(dst, cs.NewDecoder())
}

func MatchBy(dst io.Writer, excepted string, cb func()) io.Writer {
	return &MatchWriter{
		out:      dst,
		excepted: []byte(excepted),
		cb:       cb,
	}
}

func ToInt(s string, v int) int {
	if value, e := strconv.ParseInt(s, 10, 0); nil == e {
		return int(value)
	}
	return v
}

func LogString(ws io.Writer, msg string) {
	if nil != ws {
		io.WriteString(ws, "%tpt%"+msg)
	}
	log.Println(msg)
}

func SaveSessionKey(pa string, args []string, wd string) {
	args = RemoveBatchOption(args)
	var cmd = exec.Command(pa, args...)
	if len(wd) > 0 {
		cmd.Dir = wd
	}

	timer := time.AfterFunc(1*time.Minute, func() {
		defer recover()
		cmd.Process.Kill()
	})
	cmd.Stdin = strings.NewReader("y\ny\ny\ny\ny\ny\ny\ny\n")
	cmd.Run()
	timer.Stop()
}

func LookPath(executableFolder string, alias ...string) (string, bool) {
	var names []string
	for _, aliasName := range alias {
		if runtime.GOOS == "windows" {
			names = append(names, aliasName, aliasName+".bat", aliasName+".com", aliasName+".exe")
		} else {
			names = append(names, aliasName, aliasName+".sh")
		}
	}

	for _, nm := range names {
		files := []string{nm,
			filepath.Join("bin", nm),
			filepath.Join("tools", nm),
			filepath.Join("runtime_env", nm),
			filepath.Join("..", nm),
			filepath.Join("..", "bin", nm),
			filepath.Join("..", "tools", nm),
			filepath.Join("..", "runtime_env", nm),
			filepath.Join(executableFolder, nm),
			filepath.Join(executableFolder, "bin", nm),
			filepath.Join(executableFolder, "tools", nm),
			filepath.Join(executableFolder, "runtime_env", nm),
			filepath.Join(executableFolder, "..", nm),
			filepath.Join(executableFolder, "..", "bin", nm),
			filepath.Join(executableFolder, "..", "tools", nm),
			filepath.Join(executableFolder, "..", "runtime_env", nm)}
		for _, file := range files {
			// fmt.Println("====", file)
			file = config.AbsPath(file)
			if st, e := os.Stat(file); nil == e && nil != st && !st.IsDir() {
				//fmt.Println("1=====", file, e)
				return file, true
			}
		}
	}

	for _, nm := range names {
		_, err := exec.LookPath(nm)
		if nil == err {
			return nm, true
		}
	}
	return "", false
}

func RemoveBatchOption(args []string) []string {
	offset := 0
	for idx, s := range args {
		if strings.ToLower(s) == "-batch" {
			continue
		}
		if offset != idx {
			args[offset] = s
		}
		offset++
	}
	return args[:offset]
}

func AddMibDir(args []string) []string {
	hasMIBSDir := false
	for _, argument := range args {
		if "-M" == argument {
			hasMIBSDir = true
		}
	}

	if !hasMIBSDir {
		newArgs := make([]string, len(args)+2)
		newArgs[0] = "-M"
		newArgs[1] = config.Default.MIBSDir
		copy(newArgs[2:], args)
		args = newArgs
	}
	return args
}
