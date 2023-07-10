package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/webx-top/com"
)

var p = buildParam{
	GoVersion:    `1.20.5`,
	Executor:     `nging`,
	NgingVersion: `5.1.0`,
	NgingLabel:   `stable`,
	Project:      `github.com/admpub/nging`,
	WorkDir:      ``,
}

var targetNames = map[string]string{
	`linux_386`:     `linux/386`,
	`linux_amd64`:   `linux/amd64`,
	`linux_arm5`:    `linux/arm-5`,
	`linux_arm6`:    `linux/arm-6`,
	`linux_arm7`:    `linux/arm-7`,
	`linux_arm64`:   `linux/arm64`,
	`darwin_amd64`:  `darwin/amd64`,
	`windows_386`:   `windows/386`,
	`windows_amd64`: `windows/amd64`,
}

var armRegexp = regexp.MustCompile(`/arm`)

func main() {
	flag.StringVar(&p.GoVersion, `goVersion`, p.GoVersion, `--goVersion `+p.GoVersion)
	flag.StringVar(&p.Executor, `executor`, p.Executor, `--executor `+p.Executor)
	flag.StringVar(&p.NgingVersion, `ngingVersion`, p.NgingVersion, `--ngingVersion `+p.NgingVersion)
	flag.StringVar(&p.NgingLabel, `ngingLabel`, p.NgingLabel, `--ngingLabel `+p.NgingLabel)
	flag.StringVar(&p.Project, `project`, p.Project, `--project `+p.Project)
	flag.Parse()
	var err error
	p.ProjectPath, err = com.GetSrcPath(p.Project)
	if err != nil {
		com.ExitOnFailure(err.Error(), 1)
	}
	p.WorkDir = strings.TrimSuffix(strings.TrimSuffix(p.ProjectPath, `/`), p.Project)
	fmt.Println(`WorkDir: `, p.WorkDir)
	var targets []string
	var armTargets []string
	addTarget := func(target string, notNames ...bool) {
		if len(notNames) == 0 || !notNames[0] {
			target = getTarget(target)
			if len(target) == 0 {
				return
			}
		}
		if armRegexp.MatchString(target) {
			armTargets = append(armTargets, target)
		} else {
			targets = append(targets, target)
		}
	}
	var minify bool
	switch len(os.Args) {
	case 3:
		minify = isMinified(os.Args[2])
		addTarget(os.Args[1])
	case 2:
		if isMinified(os.Args[1]) {
			minify = true
			for _, t := range targetNames {
				addTarget(t, true)
			}
		} else {
			addTarget(os.Args[1])
		}
	case 1:
		for _, t := range targetNames {
			addTarget(t, true)
		}
	default:
		com.ExitOnFailure(`参数错误`)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = os.Chdir(p.ProjectPath)
	if err != nil {
		com.ExitOnFailure(err.Error(), 1)
	}
	p.NgingCommitID = execGitCommitIDCommand(ctx)
	p.NgingBuildTime = time.Now().Format(`20060102150405`)
	if minify {
		p.MinifyFlags = []string{`-s`, `-w`}
	}
	distPath := filepath.Join(p.ProjectPath, `dist`)
	fmt.Println(`DistPath: `, distPath)
	allTargets := append(targets, armTargets...)
	fmt.Printf("Building %s for %+v\n", p.Executor, allTargets)
	for _, target := range allTargets {
		parts := strings.SplitN(target, `/`, 2)
		if len(parts) != 2 {
			continue
		}
		pCopy := p
		pCopy.Target = target
		pCopy.PureGoTags = []string{`osusergo`}
		osName := parts[0]
		archName := parts[1]
		pCopy.ReleaseDir = filepath.Join(distPath, p.Executor+`_`+osName+`_`+archName)
		pCopy.goos = osName
		pCopy.goarch = archName
		if osName != `darwin` {
			pCopy.LdFlags = []string{`-extldflags`, `'-static'`}
		}
		if osName != `windows` {
			pCopy.PureGoTags = append(pCopy.PureGoTags, `netgo`)
		} else {
			pCopy.Extension = `.exe`
		}
		execGenerateCommand(ctx, pCopy)
		err := com.MkdirAll(pCopy.ReleaseDir, os.ModePerm)
		if err != nil {
			com.ExitOnFailure(err.Error(), 1)
		}
		execBuildCommand(ctx, pCopy)
		packFiles(pCopy)
	}
}

func getTarget(target string) string {
	if t, y := targetNames[target]; y {
		return t
	}
	for _, t := range targetNames {
		if t == target {
			return t
		}
	}
	return ``
}

func isMinified(arg string) bool {
	return arg == `m` || arg == `min`
}

type buildParam struct {
	GoVersion      string
	Target         string //${GOOS}/${GOARCH}
	ReleaseDir     string
	Executor       string
	Extension      string
	PureGoTags     []string
	BuildTags      []string
	NgingBuildTime string
	NgingCommitID  string
	NgingVersion   string
	NgingLabel     string
	MinifyFlags    []string
	LdFlags        []string
	Project        string
	ProjectPath    string
	WorkDir        string
	goos           string
	goarch         string
}

func (p buildParam) genLdFlagsString() string {
	ldflags := []string{}
	ldflags = append(ldflags, p.MinifyFlags...)
	ldflags = append(ldflags, p.LdFlags...)
	return `-X main.BUILD_TIME=` + p.NgingBuildTime + ` -X main.COMMIT=` + p.NgingCommitID + ` -X main.VERSION=` + p.NgingVersion + ` -X main.LABEL=` + p.NgingLabel + ` -X main.BUILD_OS=` + runtime.GOOS + ` -X main.BUILD_ARCH=` + runtime.GOARCH + ` ` + strings.Join(ldflags, ` `)
}

func execBuildCommand(ctx context.Context, p buildParam) {
	tags := []string{}
	tags = append(tags, p.PureGoTags...)
	tags = append(tags, `bindata`, `sqlite`)
	tags = append(tags, p.BuildTags...)
	cmd := exec.CommandContext(ctx, `xgo`,
		`-go`, p.GoVersion,
		`-goproxy`, `https://goproxy.cn,direct`,
		`-image`, `localhost/crazymax/xgo:`+p.GoVersion,
		`-targets`, p.Target,
		`-dest`, p.ReleaseDir,
		`-out`, p.Executor,
		`-tags`, strings.Join(tags, ` `),
		`-ldflags`, p.genLdFlagsString(),
		`./`+p.Project,
	)
	cmd.Dir = p.WorkDir
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		com.ExitOnFailure(err.Error(), 1)
	}
}

func execGenerateCommand(ctx context.Context, p buildParam) {
	cmd := exec.CommandContext(ctx, `go`, `generate`)
	cmd.Dir = p.ProjectPath
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		com.ExitOnFailure(err.Error(), 1)
	}
}

func execGitCommitIDCommand(ctx context.Context) string {
	cmd := exec.CommandContext(ctx, `git`, `rev-parse`, `HEAD`)
	cmd.Dir = p.ProjectPath
	out, err := cmd.CombinedOutput()
	if err != nil {
		com.ExitOnFailure(err.Error(), 1)
	}
	return string(out)
}

func packFiles(p buildParam) {
	files, err := filepath.Glob(filepath.Join(p.ReleaseDir, p.Executor+`-`+p.goos+`*`))
	if err != nil {
		com.ExitOnFailure(err.Error(), 1)
	}
	for _, file := range files {
		com.Rename(file, filepath.Join(p.ReleaseDir, p.Executor+p.Extension))
		break
	}
	err = com.MkdirAll(filepath.Join(p.ReleaseDir, `data`, `logs`), os.ModePerm)
	if err != nil {
		com.ExitOnFailure(err.Error(), 1)
	}
	err = com.CopyDir(filepath.Join(p.ProjectPath, `data`, `ip2region`), filepath.Join(p.ReleaseDir, `data`, `ip2region`))
	if err != nil {
		com.ExitOnFailure(err.Error(), 1)
	}
	err = com.MkdirAll(filepath.Join(p.ReleaseDir, `config`, `vhosts`), os.ModePerm)
	if err != nil {
		com.ExitOnFailure(err.Error(), 1)
	}
	err = com.Copy(filepath.Join(p.ProjectPath, `config`, `config.yaml.sample`), filepath.Join(p.ReleaseDir, `config`, `config.yaml.sample`))
	if err != nil {
		com.ExitOnFailure(err.Error(), 1)
	}
	err = com.Copy(filepath.Join(p.ProjectPath, `config`, `config.yaml.sample`), filepath.Join(p.ReleaseDir, `config`, `config.yaml.sample`))
	if err != nil {
		com.ExitOnFailure(err.Error(), 1)
	}
	files, err = filepath.Glob(filepath.Join(p.ReleaseDir, `config`, `preupgrade.*`))
	if err != nil {
		com.ExitOnFailure(err.Error(), 1)
	}
	for _, file := range files {
		err = com.Copy(file, filepath.Join(p.ReleaseDir, `config`, filepath.Base(file)))
		if err != nil {
			com.ExitOnFailure(err.Error(), 1)
		}
	}
	err = com.Copy(filepath.Join(p.ProjectPath, `config`, `ua.txt`), filepath.Join(p.ReleaseDir, `config`, `ua.txt`))
	if err != nil {
		com.ExitOnFailure(err.Error(), 1)
	}
	err = com.MkdirAll(filepath.Join(p.ReleaseDir, `public`, `upload`), os.ModePerm)
	if err != nil {
		com.ExitOnFailure(err.Error(), 1)
	}
	err = com.TarGz(p.ReleaseDir, p.ReleaseDir+`.tar.gz`)
	if err != nil {
		com.ExitOnFailure(err.Error(), 1)
	}
	err = os.RemoveAll(p.ReleaseDir)
	if err != nil {
		com.ExitOnFailure(err.Error(), 1)
	}
}
