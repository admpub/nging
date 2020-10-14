// +build go1.14
// +build !go1.16

package goloader

import (
	"errors"
	"fmt"
)

const (
	R_PCREL = 16
	// R_TLS_LE, used on 386, amd64, and ARM, resolves to the offset of the
	// thread-local symbol from the thread local base and is used to implement the
	// "local exec" model for tls access (r.Sym is not set on intel platforms but is
	// set to a TLS symbol -- runtime.tlsg -- in the linker when externally linking).
	R_TLS_LE = 17
	// R_METHODOFF resolves to a 32-bit offset from the beginning of the section
	// holding the data being relocated to the referenced symbol.
	// It is a variant of R_ADDROFF used when linking from the uncommonType of a
	// *rtype, and may be set to zero by the linker if it determines the method
	// text is unreachable by the linked program.
	R_METHODOFF = 25
)

// copy from $GOROOT/src/cmd/internal/objabi/symkind.go
const (
	// An otherwise invalid zero value for the type
	Sxxx = iota
	// Executable instructions
	STEXT
	// Read only static data
	SRODATA
	// Static data that does not contain any pointers
	SNOPTRDATA
	// Static data
	SDATA
	// Statically data that is initially all 0s
	SBSS
	// Statically data that is initially all 0s and does not contain pointers
	SNOPTRBSS
	// Thread-local data that is initially all 0s
	STLSBSS
	// Debugging data
	SDWARFINFO
	SDWARFRANGE
	SDWARFLOC
	SDWARFLINES
	// ABI alias. An ABI alias symbol is an empty symbol with a
	// single relocation with 0 size that references the native
	// function implementation symbol.
	//
	// TODO(austin): Remove this and all uses once the compiler
	// generates real ABI wrappers rather than symbol aliases.
	SABIALIAS
	// Coverage instrumentation counter for libfuzzer.
	SLIBFUZZER_EXTRA_COUNTER
	// Update cmd/link/internal/sym/AbiSymKindToSymKind for new SymKind values.

)

func addStackObject(codereloc *CodeReloc, funcname string, symbolMap map[string]uintptr) (err error) {
	return _addStackObject(codereloc, funcname, symbolMap)
}

func addDeferReturn(codereloc *CodeReloc, _func *_func) (err error) {
	funcname := gostringnocopy(&codereloc.pclntable[_func.nameoff])
	Func := codereloc.symMap[funcname].Func
	if Func != nil && len(Func.FuncData) > _FUNCDATA_OpenCodedDeferInfo {
		sym := codereloc.symMap[funcname]
		for _, r := range sym.Reloc {
			if r.Sym == codereloc.symMap[RUNTIME_DEFERRETURN] {
				//../cmd/link/internal/ld/pcln.go:pclntab
				switch codereloc.Arch {
				case ARCH_386, ARCH_AMD64:
					_func.deferreturn = uint32(r.Offset) - uint32(sym.Offset) - 1
				case ARCH_ARM32, ARCH_ARM64:
					_func.deferreturn = uint32(r.Offset) - uint32(sym.Offset)
				default:
					err = errors.New(fmt.Sprintf("not support arch:%s", codereloc.Arch))
				}
				break
			}
		}
	}
	return err
}
