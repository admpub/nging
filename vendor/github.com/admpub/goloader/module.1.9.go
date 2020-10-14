// +build go1.9
// +build !go1.10,!go1.11,!go1.12,!go1.13,!go1.14,!go1.15,!go1.16

package goloader

import (
	"cmd/objfile/goobj"
)

// PCDATA and FUNCDATA table indexes.
//
// See funcdata.h and ../cmd/internal/obj/funcdata.go.
const (
	_PCDATA_StackMapIndex       = 0
	_PCDATA_InlTreeIndex        = 1
	_FUNCDATA_ArgsPointerMaps   = 0
	_FUNCDATA_LocalsPointerMaps = 1
	_FUNCDATA_InlTree           = 2
	_ArgsSizeUnknown            = -0x80000000
)

// moduledata records information about the layout of the executable
// image. It is written by the linker. Any changes here must be
// matched changes to the code in cmd/internal/ld/symtab.go:symtab.
// moduledata is stored in read-only memory; none of the pointers here
// are visible to the garbage collector.
type moduledata struct {
	pclntable    []byte
	ftab         []functab
	filetab      []uint32
	findfunctab  uintptr
	minpc, maxpc uintptr

	text, etext           uintptr
	noptrdata, enoptrdata uintptr
	data, edata           uintptr
	bss, ebss             uintptr
	noptrbss, enoptrbss   uintptr
	end, gcdata, gcbss    uintptr
	types, etypes         uintptr

	textsectmap []textsect
	typelinks   []int32 // offsets from types
	itablinks   []*itab

	ptab []ptabEntry

	pluginpath string
	pkghashes  []modulehash

	modulename   string
	modulehashes []modulehash

	gcdatamask, gcbssmask bitvector

	typemap map[typeOff]uintptr // offset to *_rtype in previous module

	next *moduledata
}

type _func struct {
	entry   uintptr // start pc
	nameoff int32   // function name

	args int32 // in/out args size
	_    int32 // previously legacy frame size; kept for layout compatibility

	pcsp      int32
	pcfile    int32
	pcln      int32
	npcdata   int32
	nfuncdata int32
}

func init_func(symbol *goobj.Sym, nameOff, spOff, pcfileOff, pclnOff int) _func {
	fdata := _func{
		entry:     uintptr(0),
		nameoff:   int32(nameOff),
		args:      int32(symbol.Func.Args),
		pcsp:      int32(spOff),
		pcfile:    int32(pcfileOff),
		pcln:      int32(pclnOff),
		npcdata:   int32(len(symbol.Func.PCData)),
		nfuncdata: int32(len(symbol.Func.FuncData)),
	}
	return fdata
}

// inlinedCall is the encoding of entries in the FUNCDATA_InlTree table.
type inlinedCall struct {
	parent int32 // index of parent in the inltree, or < 0
	file   int32 // fileno index into filetab
	line   int32 // line number of the call site
	func_  int32 // offset into pclntab for name of called function
}

func initInlinedCall(codereloc *CodeReloc, inl goobj.InlinedCall, _func *_func) inlinedCall {
	return inlinedCall{
		parent: int32(inl.Parent),
		file:   int32(findFileTab(codereloc, inl.File)),
		line:   int32(inl.Line),
		func_:  int32(codereloc.namemap[inl.Func.Name])}
}

func addInlineTree(codereloc *CodeReloc, _func *_func, objsym objSym) (err error) {
	return _addInlineTree(codereloc, _func, objsym)
}
