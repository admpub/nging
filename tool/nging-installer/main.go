package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	godl "github.com/admpub/go-download/v2"
	"github.com/admpub/go-download/v2/progressbar"
	"github.com/webx-top/com"
)

var osName = runtime.GOOS
var archName = runtime.GOARCH
var operate = `install`
var version = `5.1.0`
var saveDir = `nging`
var softwareURL = `https://img.nging.coscms.com/nging/v%s/`
var binName = "nging"
var fileName = fmt.Sprintf("nging_%s_%s", osName, archName)
var fileFullName = fileName + ".tar.gz"
var softwareFullURL string
var workDir string
var local string

func parseArgs() {
	if len(os.Args) > 4 {
		os.Args = os.Args[0:4]
	}
	switch len(os.Args) {
	case 4:
		saveDir = os.Args[3]
		fallthrough
	case 3:
		version = os.Args[2]
		fallthrough
	case 2:
		operate = os.Args[1]
	}
}

var supports = map[string][]string{
	`darwin`:  {`amd64`, `arm64`},
	`linux`:   {`386`, `amd64`, `arm64`, `arm-7`, `arm-6`, `arm-5`},
	`windows`: {`386`, `amd64`},
}

func main() {
	flag.StringVar(&local, `local`, ``, `--local ./nging_darwin_amd64.tar.gz`)
	flag.Parse()
	if len(local) > 0 {
		fmt.Println(`local: `, local)
	}

	if _, ok := supports[osName]; !ok {
		com.ExitOnFailure(`Unsupported System:`+osName, 1)
	}
	switch archName {
	case `x86_64`:
		archName = "amd64"
	case "i386", "i686":
		archName = "386"
	case "aarch64_be", "aarch64", "armv8b", "armv8l", "armv8", "arm64":
		archName = "arm64"
	case "armv7":
		archName = "arm-7"
	case "armv7l":
		archName = "arm-6"
	case "armv6":
		archName = "arm-6"
	case "armv5", "arm":
		archName = "arm-5"
	}
	if !com.InSlice(archName, supports[osName]) {
		com.ExitOnFailure(`Unsupported Arch:`+archName, 1)
	}
	parseArgs()
	softwareFullURL = fmt.Sprintf(softwareURL, version) + fileFullName
	var err error
	workDir, err = filepath.Abs(saveDir)
	if err != nil {
		com.ExitOnFailure(err.Error(), 1)
	}
	switch operate {
	case `un`, `uninstall`:
		uninstall()
	case `up`, `upgrade`:
		upgrade()
	case `install`:
		install()
	default:
		install()
	}
}

func downloadAndExtract() {
	compressedFile := fileFullName
	if len(local) == 0 {
		godlOpt := &godl.Options{}
		progress := progressbar.New(godlOpt, 50)
		defer progress.Wait()

		_, err := godl.Download(softwareFullURL, compressedFile, godlOpt)
		if err != nil {
			com.ExitOnFailure(err.Error(), 1)
		}
	} else {
		compressedFile = local
	}
	err := com.MkdirAll(saveDir, os.ModePerm)
	if err != nil {
		com.ExitOnFailure(err.Error(), 1)
	}
	_, err = com.UnTarGz(compressedFile, saveDir)
	if err != nil {
		com.ExitOnFailure(err.Error(), 1)
	}
	distDir := filepath.Join(saveDir, fileName)
	err = com.CopyDir(distDir, saveDir)
	if err != nil {
		com.ExitOnFailure(err.Error(), 1)
	}
	os.RemoveAll(distDir)
	if len(local) == 0 {
		os.Remove(compressedFile)
	}
	os.Chmod(filepath.Join(saveDir, binName), os.ModePerm)
}

func execServiceCommand(op string, mustSucceed ...bool) error {
	cmd := exec.Command(`./`+binName, `service`, op)
	cmd.Dir = workDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		if len(mustSucceed) == 0 || mustSucceed[0] {
			com.ExitOnFailure(err.Error(), 1)
		}
		return err
	}
	fmt.Println(len(out))
	return err
}

func install() {
	downloadAndExtract()
	execServiceCommand(`install`)
	execServiceCommand(`start`)
	fmt.Println(`🎉 Congratulations! Installed successfully.`)
}

func uninstall() {
	execServiceCommand(`stop`)
	execServiceCommand(`uninstall`)
	fmt.Println(`🎉 Congratulations! Successfully uninstalled.`)
	err := os.RemoveAll(saveDir)
	if err != nil {
		com.ExitOnFailure(err.Error(), 1)
	}
	fmt.Println(`🎉 Congratulations! File deleted successfully.`)
}

func upgrade() {
	execServiceCommand(`stop`)
	execServiceCommand(`stop`, false)
	downloadAndExtract()
	execServiceCommand(`start`)
	fmt.Println(`🎉 Congratulations! Successfully upgraded.`)
}
