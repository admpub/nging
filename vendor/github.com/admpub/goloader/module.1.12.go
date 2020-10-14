// +build go1.12
// +build !go1.13,!go1.14,!go1.15,!go1.16

package goloader

import (
	"cmd/objfile/goobj"
	"cmd/objfile/objabi"
	"strings"
)

// PCDATA and FUNCDATA table indexes.
//
// See funcdata.h and ../cmd/internal/objabi/funcdata.go.
const (
	_PCDATA_StackMapIndex       = 0
	_PCDATA_InlTreeIndex        = 1
	_PCDATA_RegMapIndex         = 2
	_FUNCDATA_ArgsPointerMaps   = 0
	_FUNCDATA_LocalsPointerMaps = 1
	_FUNCDATA_InlTree           = 2
	_FUNCDATA_RegPointerMaps    = 3
	_FUNCDATA_StackObjects      = 4
	_ArgsSizeUnknown            = -0x80000000
)

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

	hasmain uint8 // 1 if module contains the main function, 0 otherwise

	gcdatamask, gcbssmask bitvector

	typemap map[typeOff]uintptr // offset to *_rtype in previous module

	bad bool // module failed to load and should be ignored

	next *moduledata
}

// A funcID identifies particular functions that need to be treated
// specially by the runtime.
// Note that in some situations involving plugins, there may be multiple
// copies of a particular special runtime function.
// Note: this list must match the list in cmd/internal/objabi/funcid.go.
type funcID uint8

type _func struct {
	entry   uintptr // start pc
	nameoff int32   // function name

	args int32 // in/out args size
	_    int32 // previously legacy frame size; kept for layout compatibility

	pcsp      int32
	pcfile    int32
	pcln      int32
	npcdata   int32
	funcID    funcID  // set for certain special runtime functions
	_         [2]int8 // unused
	nfuncdata uint8
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
		funcID:    funcID(objabi.GetFuncID(symbol.Name, strings.TrimLeft(symbol.Func.File[0], FILE_SYM_PREFIX))),
		nfuncdata: uint8(len(symbol.Func.FuncData)),
	}
	return fdata
}

// inlinedCall is the encoding of entries in the FUNCDATA_InlTree table.
type inlinedCall struct {
	parent   int16  // index of parent in the inltree, or < 0
	funcID   funcID // type of the called function
	_        byte
	file     int32 // fileno index into filetab
	line     int32 // line number of the call site
	func_    int32 // offset into pclntab for name of called function
	parentPc int32 // position of an instruction whose source position is the call site (offset from entry)
}

func initInlinedCall(codereloc *CodeReloc, inl goobj.InlinedCall, _func *_func) inlinedCall {
	return inlinedCall{
		parent:   int16(inl.Parent),
		funcID:   _func.funcID,
		file:     int32(findFileTab(codereloc, inl.File)),
		line:     int32(inl.Line),
		func_:    int32(codereloc.namemap[inl.Func.Name]),
		parentPc: int32(inl.ParentPC)}
}

func addInlineTree(codereloc *CodeReloc, _func *_func, objsym objSym) (err error) {
	return _addInlineTree(codereloc, _func, objsym)
}
