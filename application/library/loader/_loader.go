package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/admpub/goloader"
	"github.com/admpub/log"
	"github.com/webx-top/echo"
)

var (
	symPtr  = make(map[string]uintptr)
	relocs  = make(map[string]*goloader.CodeReloc)
	modules = make(map[string]*goloader.CodeModule)
)

func init() {
	err := goloader.RegSymbol(symPtr)
	if err != nil {
		panic(err)
	}
}

func RegTypes(interfaces ...interface{}) {
	goloader.RegTypes(symPtr, interfaces...)
}

func Register(moduleName string, libFiles ...string) error {
	var (
		files    []string
		pkgPaths []string
	)
	for _, f := range libFiles {
		a := strings.SplitN(f, `:`, 2)
		switch len(a) {
		case 2:
			pkgPaths = append(pkgPaths, a[1])
			files = append(files, a[0])
		case 1:
			pkgPaths = append(pkgPaths, ``)
			files = append(files, a[0])
		}
	}
	reloc, err := goloader.ReadObjs(files, pkgPaths)
	if err != nil {
		return err
	}
	relocs[moduleName] = reloc
	_, err = Load(moduleName)
	return err
}

func Unregister(moduleName string) {
	if _, y := relocs[moduleName]; y {
		Unload(moduleName)
		delete(relocs, moduleName)
	}
}

func Load(moduleName string) (*goloader.CodeModule, error) {
	if m, y := modules[moduleName]; y {
		return m, nil
	}
	reloc, ok := relocs[moduleName]
	if !ok {
		return nil, fmt.Errorf(`Load error! not find reloc: %s`, moduleName)
	}
	codeModule, err := goloader.Load(reloc, symPtr)
	if err != nil {
		return nil, fmt.Errorf(`Load error: %w`, err)
	}
	modules[moduleName] = codeModule
	return codeModule, nil
}

func Unload(moduleName string) {
	m, y := modules[moduleName]
	if !y {
		return
	}
	m.Unload()
	delete(modules, moduleName)
}

func Exec(moduleName string, funcName ...string) error {
	codeModule, err := Load(moduleName)
	if err != nil {
		return err
	}
	var run string
	if len(funcName) > 0 {
		run = funcName[0]
	}
	if len(run) == 0 {
		run = `main.main`
	}
	runFuncPtr := codeModule.Syms[run]
	if runFuncPtr == 0 {
		return fmt.Errorf(`Load error! not find function: %s`, run)
	}
	funcPtrContainer := (uintptr)(unsafe.Pointer(&runFuncPtr))
	runFunc := *(*func())(unsafe.Pointer(&funcPtrContainer))
	runFunc()
	return nil
}

func LoadPlugins() error {
	pluginsDir := filepath.Join(echo.Wd(), `support/plugins`)
	files, err := os.ReadDir(pluginsDir)
	if err != nil {
		return err
	}

	// Load Plugin
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if !strings.HasSuffix(f.Name(), ".o") {
			continue
		}
		moduleName := strings.TrimSuffix(f.Name(), ".o")
		Register(moduleName, pluginsDir+`/`+f.Name())
		log.Info(`Loading plugin: `, moduleName)
		if err := Exec(moduleName); err != nil {
			log.Error(`[`+moduleName+`] `, `Initialize plugin err: `, err)
		}
	}
	return err
}

func LoadPlugin(moduleName string) error {
	plugin := filepath.Join(echo.Wd(), `support/plugins/`+moduleName+`.o`)
	Register(moduleName, plugin)
	log.Info(`Loading plugin: `, moduleName)
	err := Exec(moduleName)
	if err != nil {
		log.Error(`[`+moduleName+`] `, `Initialize plugin err: `, err)
	}
	return err
}
