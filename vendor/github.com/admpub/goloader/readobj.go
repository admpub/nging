package goloader

import (
	"cmd/objfile/goobj"
	"fmt"
	"os"
	"strings"
)

func readObj(f *os.File, reloc *CodeReloc, objSymMap map[string]objSym, pkgpath *string) error {
	if pkgpath == nil || *pkgpath == EMPTY_STRING {
		defaultPkgPath := DEFAULT_PKGPATH
		pkgpath = &defaultPkgPath
	}
	obj, err := goobj.Parse(f, *pkgpath)
	if err != nil {
		return fmt.Errorf("read error: %v", err)
	}
	if len(reloc.Arch) != 0 && reloc.Arch != obj.Arch {
		return fmt.Errorf("read obj error: Arch %s != Arch %s", reloc.Arch, obj.Arch)
	}
	reloc.Arch = obj.Arch
	for _, sym := range obj.Syms {
		objSymMap[sym.Name] = objSym{
			sym:     sym,
			file:    f,
			pkgpath: *pkgpath,
		}
		for index, loc := range sym.Reloc {
			sym.Reloc[index].Sym.Name = strings.Replace(loc.Sym.Name, EMPTY_PKGPATH, *pkgpath, -1)
		}
		if sym.Func != nil {
			for index, FuncData := range sym.Func.FuncData {
				sym.Func.FuncData[index].Sym.Name = strings.Replace(FuncData.Sym.Name, EMPTY_PKGPATH, *pkgpath, -1)
			}
		}
	}
	return nil
}

func ReadObj(f *os.File) (*CodeReloc, error) {
	reloc := &CodeReloc{symMap: make(map[string]*Sym), stkmaps: make(map[string][]byte), namemap: make(map[string]int)}
	reloc.pclntable = append(reloc.pclntable, x86moduleHead...)
	objSymMap := make(map[string]objSym)
	err := readObj(f, reloc, objSymMap, nil)
	if err != nil {
		return nil, err
	}
	//static_tmp is 0, golang compile not allocate memory.
	reloc.data = append(reloc.data, make([]byte, IntSize)...)
	for _, objSym := range objSymMap {
		if objSym.sym.Kind == STEXT && objSym.sym.DupOK == false {
			_, err := relocSym(reloc, objSym.sym.Name, objSymMap)
			if err != nil {
				return nil, err
			}
		}
	}
	if reloc.Arch == ARCH_ARM32 || reloc.Arch == ARCH_ARM64 {
		copy(reloc.pclntable, armmoduleHead)
	}
	return reloc, err
}

func ReadObjs(files []string, pkgPath []string) (*CodeReloc, error) {
	reloc := &CodeReloc{symMap: make(map[string]*Sym), stkmaps: make(map[string][]byte), namemap: make(map[string]int)}
	reloc.pclntable = append(reloc.pclntable, x86moduleHead...)
	objSymMap := make(map[string]objSym)
	for i, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		err = readObj(f, reloc, objSymMap, &(pkgPath[i]))
		if err != nil {
			return nil, err
		}
	}
	//static_tmp is 0, golang compile not allocate memory.
	reloc.data = append(reloc.data, make([]byte, IntSize)...)
	for _, objSym := range objSymMap {
		if objSym.sym.Kind == STEXT && objSym.sym.DupOK == false {
			_, err := relocSym(reloc, objSym.sym.Name, objSymMap)
			if err != nil {
				return nil, err
			}
		}
	}
	if reloc.Arch == ARCH_ARM32 || reloc.Arch == ARCH_ARM64 {
		copy(reloc.pclntable, armmoduleHead)
	}
	return reloc, nil
}
