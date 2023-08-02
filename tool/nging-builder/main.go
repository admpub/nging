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

	"github.com/admpub/confl"
	"github.com/webx-top/com"
)

var p = buildParam{}

var c = Config{
	GoVersion:    `1.20.6`,
	Executor:     `nging`,
	NgingVersion: `5.1.1`,
	NgingLabel:   `stable`,
	Project:      `github.com/admpub/nging`,
	VendorMiscDirs: map[string][]string{
		`*`: {
			`vendor/github.com/nging-plugins/caddymanager/template/`,
			`vendor/github.com/nging-plugins/collector/template/`,
			`vendor/github.com/nging-plugins/dbmanager/template/`,
			`vendor/github.com/nging-plugins/dbmanager/public/assets/`,
			`vendor/github.com/nging-plugins/ddnsmanager/template/`,
			`vendor/github.com/nging-plugins/dlmanager/template/`,
			`vendor/github.com/nging-plugins/frpmanager/template/`,
			`vendor/github.com/nging-plugins/ftpmanager/template/`,
			`vendor/github.com/nging-plugins/servermanager/template/`,
			`vendor/github.com/nging-plugins/sshmanager/template/`,
			`vendor/github.com/nging-plugins/webauthn/template/`,
		},
		`linux`: {
			`vendor/github.com/nging-plugins/firewallmanager/template/`,
		},
		`!linux`: {},
	},
	BuildTags: []string{`bindata`, `sqlite`},
	CopyFiles: []string{`config/ua.txt`, `config/config.yaml.sample`, `data/ip2region`, `config/preupgrade.*`},
	MakeDirs:  []string{`public/upload`, `config/vhosts`, `data/logs`},
	Compiler:  `xgo`,
}

var targetNames = map[string]string{
	`linux_386`:     `linux/386`,
	`linux_amd64`:   `linux/amd64`,
	`linux_arm5`:    `linux/arm-5`,
	`linux_arm6`:    `linux/arm-6`,
	`linux_arm7`:    `linux/arm-7`,
	`linux_arm64`:   `linux/arm64`,
	`darwin_amd64`:  `darwin/amd64`,
	`darwin_arm64`:  `darwin/arm64`,
	`windows_386`:   `windows/386`,
	`windows_amd64`: `windows/amd64`,
}

var armRegexp = regexp.MustCompile(`/arm`)
var configFile = `./builder.conf`

func main() {
	flag.StringVar(&configFile, `conf`, configFile, `--conf `+configFile)
	flag.Parse()

	_, err := confl.DecodeFile(configFile, &c)
	if err != nil {
		com.ExitOnFailure(err.Error(), 1)
	}
	c.apply()
	p.ProjectPath, err = com.GetSrcPath(p.Project)
	if err != nil {
		com.ExitOnFailure(err.Error(), 1)
	}
	p.WorkDir = strings.TrimSuffix(strings.TrimSuffix(p.ProjectPath, `/`), p.Project)
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
		switch {
		case isMinified(os.Args[1]):
			minify = true
			for _, t := range targetNames {
				addTarget(t, true)
			}
		case os.Args[1] == `genConfig`:
			b, err := confl.Marshal(c)
			if err != nil {
				com.ExitOnFailure(err.Error(), 1)
			}
			err = os.WriteFile(configFile, b, os.ModePerm)
			if err != nil {
				com.ExitOnFailure(err.Error(), 1)
			}
			com.ExitOnSuccess(`successully generate config file: ` + configFile)
			return
		case os.Args[1] == `makeGen`:
			makeGenerateCommandComment(c)
			return
		default:
			addTarget(os.Args[1])
		}
	case 1:
		for _, t := range targetNames {
			addTarget(t, true)
		}
	default:
		com.ExitOnFailure(`参数错误`)
	}
	fmt.Println(`WorkDir: `, p.WorkDir)
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
	Config
	Target         string //${GOOS}/${GOARCH}
	ReleaseDir     string
	Extension      string
	PureGoTags     []string
	NgingBuildTime string
	NgingCommitID  string
	MinifyFlags    []string
	LdFlags        []string
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

func (p buildParam) genEnvVars() []string {
	env := []string{`GOOS=` + p.goos}
	parts := strings.SplitN(p.goarch, `-`, 2)
	if parts[0] == `arm` {
		env = append(env, `GOARCH=`+parts[0])
		if len(parts) == 0 {
			env = append(env, `GOARM=`+parts[1])
		}
	} else {
		env = append(env, `GOARCH=`+p.goarch)
	}
	return env
}

func execBuildCommand(ctx context.Context, p buildParam) {
	tags := []string{}
	tags = append(tags, p.PureGoTags...)
	tags = append(tags, p.BuildTags...)
	var args []string
	var env []string
	var workDir string
	var compiler string
	switch p.Compiler {
	case `go`:
		workDir = filepath.Join(p.WorkDir, p.Project)
		compiler = p.Compiler
		com.MkdirAll(p.ReleaseDir, os.ModePerm)
		args = []string{`build`,
			`-tags`, strings.Join(tags, ` `),
			`-ldflags`, p.genLdFlagsString(),
			`-o`, filepath.Join(p.ReleaseDir, p.Executor+`-`+p.goos+`-`+p.goarch),
		}
		env = append(env, os.Environ()...)
		env = append(env, p.genEnvVars()...)
		if p.CgoEnabled {
			env = append(env, `CGO_ENABLED=1`)
		} else {
			env = append(env, `CGO_ENABLED=0`)
		}
	case `xgo`:
		fallthrough
	default:
		workDir = p.WorkDir
		compiler = `xgo`
		args = []string{
			`-go`, p.GoVersion,
			`-goproxy`, `https://goproxy.cn,direct`,
			`-image`, `localhost/crazymax/xgo:` + p.GoVersion,
			`-targets`, p.Target,
			`-dest`, p.ReleaseDir,
			`-out`, p.Executor,
			`-tags`, strings.Join(tags, ` `),
			`-ldflags`, p.genLdFlagsString(),
			`./` + p.Project,
		}
	}
	cmd := exec.CommandContext(ctx, compiler, args...)
	cmd.Dir = workDir
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Env = env
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
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, p.genEnvVars()...)
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
	for _, copyFile := range p.CopyFiles {
		f := filepath.Join(p.ProjectPath, copyFile)
		if strings.Contains(f, `*`) {
			files, err = filepath.Glob(f)
			if err != nil {
				com.ExitOnFailure(err.Error(), 1)
			}
			for _, file := range files {
				destFile := filepath.Join(p.ReleaseDir, strings.TrimPrefix(file, p.ProjectPath))
				com.MkdirAll(filepath.Dir(destFile), os.ModePerm)
				err = com.Copy(file, destFile)
				if err != nil {
					com.ExitOnFailure(err.Error(), 1)
				}
			}
			continue
		}
		if com.IsDir(f) {
			err = com.CopyDir(f, filepath.Join(p.ReleaseDir, copyFile))
			if err != nil {
				com.ExitOnFailure(err.Error(), 1)
			}
			continue
		}
		destFile := filepath.Join(p.ReleaseDir, copyFile)
		com.MkdirAll(filepath.Dir(destFile), os.ModePerm)
		err = com.Copy(f, destFile)
		if err != nil {
			com.ExitOnFailure(err.Error(), 1)
		}
	}
	for _, newDir := range p.MakeDirs {
		err = com.MkdirAll(filepath.Join(p.ReleaseDir, newDir), os.ModePerm)
		if err != nil {
			com.ExitOnFailure(err.Error(), 1)
		}
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

func genComment(vendorMiscDirs ...string) string {
	comment := "//go:generate go install github.com/admpub/bindata/v3/go-bindata@latest\n"
	comment += `//go:generate go-bindata -fs -o bindata_assetfs.go -ignore "\\.(git|svn|DS_Store|less|scss|gitkeep)$" -minify "\\.(js|css)$" -tags bindata`
	prefixes := []string{}
	miscDirs := []string{`public/assets/`, `template/`, `config/i18n/`}
	miscDirs = append(miscDirs, vendorMiscDirs...)
	uniquePrefixes := map[string]struct{}{}
	for k, v := range miscDirs {
		if !strings.HasSuffix(v, `/...`) {
			if !strings.HasSuffix(v, `/`) {
				v += `/`
			}
			v += `...`
		}
		if strings.HasPrefix(v, `vendor/`) {
			parts := strings.SplitN(v, `/`, 5)
			if len(parts) == 5 {
				prefix := strings.Join(parts[0:4], `/`) + `/`
				if _, ok := uniquePrefixes[prefix]; !ok {
					uniquePrefixes[prefix] = struct{}{}
					prefixes = append(prefixes, prefix)
				}
			}
		}
		miscDirs[k] = v
	}
	comment += ` -prefix "` + strings.Join(prefixes, `|`) + `" `
	comment += strings.Join(miscDirs, ` `)
	return comment
}

func makeGenerateCommandComment(c Config) {
	dfts := c.VendorMiscDirs[`*`]
	for osName, miscDirs := range c.VendorMiscDirs {
		if osName == `*` {
			continue
		}
		dirs := make([]string, 0, len(dfts)+len(miscDirs))
		dirs = append(dirs, dfts...)
		dirs = append(dirs, miscDirs...)
		fileName := `main_`
		if strings.HasPrefix(osName, `!`) {
			fileName += `non` + strings.TrimPrefix(osName, `!`)
		} else {
			fileName += osName
		}
		fileName += `.go`
		filePath := filepath.Join(p.ProjectPath, fileName)
		fileContent := "//go:build " + osName + "\n\n"
		fileContent += "package main\n\n"
		fileContent += genComment(dirs...) + "\n\n"
		b, err := os.ReadFile(filePath)
		if err == nil {
			old := string(b)
			pos := strings.Index(old, `import `)
			if pos > -1 {
				fileContent += old[pos:]
			}
		}
		err = os.WriteFile(filePath, []byte(fileContent), os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

type Config struct {
	GoVersion      string
	Executor       string
	NgingVersion   string
	NgingLabel     string
	Project        string
	VendorMiscDirs map[string][]string // key: GOOS
	BuildTags      []string
	CopyFiles      []string
	MakeDirs       []string
	Compiler       string
	CgoEnabled     bool
}

func (a Config) apply() {
	if len(a.GoVersion) > 0 {
		p.GoVersion = a.GoVersion
	}
	if len(a.Executor) > 0 {
		p.Executor = a.Executor
	}
	if len(a.NgingVersion) > 0 {
		p.NgingVersion = a.NgingVersion
	}
	if len(a.NgingLabel) > 0 {
		p.NgingLabel = a.NgingLabel
	}
	if len(a.Project) > 0 {
		p.Project = a.Project
	}
	p.BuildTags = a.BuildTags
	p.CopyFiles = a.CopyFiles
	p.MakeDirs = a.MakeDirs
	p.Compiler = a.Compiler
	p.CgoEnabled = a.CgoEnabled
}
