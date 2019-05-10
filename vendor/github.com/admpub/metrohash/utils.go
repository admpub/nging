package metrohash

// Rotate-right (asm RORQ) functions
// General form : rotr(v, k) := (v >> k) | (v << (64-k))
// Implemented to force go compiler to generate rorq/rolq statements

func rotr64_21(v uint64) uint64 {
	return (v >> 21) | (v << (64 - 21))
}

func rotr64_26(v uint64) uint64 {
	return (v >> 26) | (v << (64 - 26))
}

func rotr64_28(v uint64) uint64 {
	return (v >> 28) | (v << (64 - 28))
}

func rotr64_29(v uint64) uint64 {
	return (v >> 29) | (v << (64 - 29))
}

func rotr64_37(v uint64) uint64 {
	return (v >> 37) | (v << (64 - 37))
}

func rotr64_48(v uint64) uint64 {
	return (v >> 48) | (v << (64 - 48))
}

func rotr64_55(v uint64) uint64 {
	return (v >> 55) | (v << (64 - 55))
}
