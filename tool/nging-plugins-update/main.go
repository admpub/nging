package main

import (
	"bytes"
	"context"
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/webx-top/com"
)

var (
	gendbschema bool
	gittag      string // `v1.9.0`
	updatedeps  bool
)

func main() {
	flag.BoolVar(&gendbschema, "gendbschema", false, "Generate DB Schema")
	flag.StringVar(&gittag, "gittag", "", "Git Tag Version")
	flag.BoolVar(&updatedeps, "updatedeps", false, "Update Dependencies")
	flag.Parse()

	UpdateNgingPlugins()
}

func UpdateNgingPlugins() {
	dir := "../../nging-plugins"
	dir, _ = filepath.Abs(dir)
	dirs, err := os.ReadDir(dir)
	if err != nil {
		log.Panic(err)
	}
	ctx := context.Background()
	for _, f := range dirs {
		if !f.IsDir() {
			continue
		}
		pluginDir := filepath.Join(dir, f.Name())
		os.Chdir(pluginDir)
		if gendbschema {
			scriptPath := filepath.Join(pluginDir, "gen_dbschema.sh")
			if fi, err := os.Stat(scriptPath); err == nil && !fi.IsDir() {
				b, err := os.ReadFile(scriptPath)
				if err != nil {
					log.Panic(err)
				}
				if len(b) == 0 {
					continue
				}
				b = bytes.Replace(b, []byte(`dbgenerator `), []byte(`dbgenerator -container mysql8 `), 1)
				scriptPathNew := scriptPath + `.new.sh`
				err = os.WriteFile(scriptPathNew, b, 0755)
				if err != nil {
					log.Panic(err)
				}
				log.Printf("Exec %s\n", scriptPathNew)
				execScriptCommand(ctx, scriptPathNew)
				os.Remove(scriptPathNew)
			}
		}
		if updatedeps {
			log.Printf("Update Dependencies for %s\n", f.Name())
			execCommand(ctx, `go`, `get`, `-u`)
		}
		execGoModCommand(ctx)
		execGitAddCommand(ctx)
		execGitCommitCommand(ctx)
		execGitPushCommand(ctx)
		if gittag != "" {
			execGitTagCommand(ctx, gittag)
			execGitTagPushCommand(ctx, gittag)
		}
		//execGitTagPushCommand(ctx, `v1.9.0`)
	}
}

func execScriptCommand(ctx context.Context, scriptPath string) {
	execCommand(ctx, `bash`, `-c`, scriptPath)
}

func execCommand(ctx context.Context, exe string, args ...string) {
	cmd := exec.CommandContext(ctx, exe, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		com.ExitOnFailure(err.Error(), 1)
	}
}

func execGoModCommand(ctx context.Context) {
	execCommand(ctx, `go`, `mod`, `tidy`, `-compat=1.17`)
}

func execGitAddCommand(ctx context.Context) {
	execCommand(ctx, `git`, `add`, `.`)
}

func execGitCommitCommand(ctx context.Context) {
	execCommand(ctx, `git`, `commit`, `-m`, `Update`)
}

func execGitPushCommand(ctx context.Context) {
	execCommand(ctx, `git`, `push`)
}

func execGitTagCommand(ctx context.Context, ver string) {
	execCommand(ctx, `git`, `tag`, `-a`, ver, `-m`, `Update`)
}

func execGitTagPushCommand(ctx context.Context, ver string) {
	execCommand(ctx, `git`, `push`, `origin`, ver)
}
