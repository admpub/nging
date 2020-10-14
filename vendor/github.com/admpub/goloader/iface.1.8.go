// +build go1.8
// +build !go1.10,!go1.11,!go1.12,!go1.13,!go1.14,!go1.15,!go1.16

package goloader

import (
	"unsafe"
)

// layout of Itab known to compilers
// allocated in non-garbage-collected memory
// Needs to be in sync with
// ../cmd/compile/internal/gc/reflect.go:/^func.dumptypestructs.
type itab struct {
	inter  *interfacetype
	_type  *_type
	link   *itab
	bad    int32
	inhash int32      // has this itab been added to hash?
	fun    [1]uintptr // variable sized
}

// See: src/runtime/iface.go
const hashSize = 1009

//go:linkname hash runtime.hash
var hash [hashSize]*itab

//go:linkname ifaceLock runtime.ifaceLock
var ifaceLock mutex

//go:linkname itabhash runtime.itabhash
func itabhash(inter *interfacetype, typ *_type) uint32

//go:linkname additab runtime.additab
func additab(m *itab, locked, canfail bool)

func additabs(module *moduledata) {
	lock(&ifaceLock)
	for _, itab := range module.itablinks {
		if itab.inhash == 0 {
			methods := itab._type.uncommon().methods()
			for k := 0; k < len(methods); k++ {
				for m := 0; m < len(itab.inter.mhdr); m++ {
					if itab.inter.typ.nameOff(itab.inter.mhdr[m].name).name() ==
						itab._type.nameOff(methods[k].name).name() {
						itype := uintptr(unsafe.Pointer(itab.inter.typ.typeOff(itab.inter.mhdr[m].ityp)))
						module.typemap[methods[k].mtyp] = itype
					}
				}
			}
			additab(itab, true, false)
		}
	}
	unlock(&ifaceLock)
}

func removeitabs(module *moduledata) bool {
	lock(&ifaceLock)
	defer unlock(&ifaceLock)

	//the itab alloc by runtime.persistentalloc, can't free
	for index, h := range &hash {
		last := h
		for m := h; m != nil; m = m.link {
			uintptrm := uintptr(unsafe.Pointer(m))
			inter := uintptr(unsafe.Pointer(m.inter))
			_type := uintptr(unsafe.Pointer(m._type))
			if (inter >= module.types && inter <= module.etypes) || (_type >= module.types && _type <= module.etypes) ||
				(uintptrm >= module.types && uintptrm <= module.etypes) {
				if m == h {
					hash[index] = m.link
				} else {
					last.link = m.link
				}
			}
			last = m
		}
	}
	return true
}
