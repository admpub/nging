package goloader

import "fmt"

func dumpPCData(b []byte, prefix string) {
	fmt.Println(prefix, b)
	var pc uintptr
	val := int32(-1)
	var ok bool
	b, ok = step(b, &pc, &val, true)
	for {
		if !ok || len(b) <= 0 {
			fmt.Println(prefix, "step end")
			break
		}
		fmt.Println(prefix, "pc:", pc, "val:", val)
		b, ok = step(b, &pc, &val, false)
	}
}

func dumpStackMap(f interface{}) {
	finfo := findfunc(getFunctionPtr(f))
	fmt.Println(funcname(finfo))
	stkmap := (*stackmap)(funcdata(finfo, _FUNCDATA_LocalsPointerMaps))
	fmt.Printf("%v %p\n", stkmap, stkmap)
}
