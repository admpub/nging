package goloader

import (
	"cmd/objfile/objfile"
	"os"
	"reflect"
	"strings"
	"unsafe"
)

// See reflect/value.go emptyInterface
type emptyInterface struct {
	typ  unsafe.Pointer
	word unsafe.Pointer
}

// See reflect/value.go sliceHeader
type sliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}

func typelinksinit(symPtr map[string]uintptr) {
	md := firstmoduledata
	for _, tl := range md.typelinks {
		t := (*_type)(adduintptr(md.types, int(tl)))
		if md.typemap != nil {
			t = (*_type)(adduintptr(md.typemap[typeOff(tl)], 0))
		}
		switch t.Kind() {
		case reflect.Ptr:
			name := t.nameOff(t.str).name()
			element := *(**_type)(add(unsafe.Pointer(t), unsafe.Sizeof(_type{})))
			pkgpath := t.PkgPath()
			if element != nil && pkgpath == EMPTY_STRING {
				pkgpath = element.PkgPath()
			}
			name = strings.Replace(name, pkgname(pkgpath), pkgpath, 1)
			if element != nil {
				symPtr[TYPE_PREFIX+name[1:]] = uintptr(unsafe.Pointer(element))
			}
			symPtr[TYPE_PREFIX+name] = uintptr(unsafe.Pointer(t))
		default:
		}
	}
	for _, f := range md.ftab {
		_func := (*_func)(unsafe.Pointer((&md.pclntable[f.funcoff])))
		name := gostringnocopy(&md.pclntable[_func.nameoff])
		if !strings.HasPrefix(name, TYPE_DOUBLE_DOT_PREFIX) && _func.entry < md.etext {
			symPtr[name] = _func.entry
		}
	}
}

func RegSymbol(symPtr map[string]uintptr) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	f, err := objfile.Open(exe)
	if err != nil {
		return err
	}
	defer f.Close()

	typelinksinit(symPtr)
	syms, err := f.Symbols()
	for _, sym := range syms {
		if sym.Name == OS_STDOUT {
			symPtr[sym.Name] = uintptr(sym.Addr)
		}
	}
	addroff := int64(uintptr(unsafe.Pointer(&os.Stdout))) - int64(symPtr[OS_STDOUT])
	for _, sym := range syms {
		code := strings.ToUpper(string(sym.Code))
		if code == "B" || code == "D" {
			symPtr[sym.Name] = uintptr(int64(sym.Addr) + addroff)
		}
		if strings.HasPrefix(sym.Name, ITAB_PREFIX) {
			symPtr[sym.Name] = uintptr(int64(sym.Addr) + addroff)
		}
	}
	return nil
}

func regTLS(symPtr map[string]uintptr, offset int) {
	//FUNCTION HEADER
	//x86/amd64
	//asm:		MOVQ (TLS), CX
	//bytes:	0x488b0c2500000000
	funcptr := getFunctionPtr(regTLS)
	tlsptr := *(*uint32)(adduintptr(funcptr, offset))
	symPtr[TLSNAME] = uintptr(tlsptr)
}

func regFunc(symPtr map[string]uintptr, name string, function interface{}) {
	symPtr[name] = getFunctionPtr(function)
}

func getFunctionPtr(function interface{}) uintptr {
	return *(*uintptr)((*emptyInterface)(unsafe.Pointer(&function)).word)
}
