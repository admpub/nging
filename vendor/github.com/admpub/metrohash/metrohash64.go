package metrohash

import (
	"encoding/binary"
	"encoding/hex"
	"hash"
)

var (
	metrohash64 *MetroHash64
	_           hash.Hash64 = metrohash64
	_           hash.Hash   = metrohash64
)

// Constants
const (
	k0         = 0xD6D018F5
	k1         = 0xA2AA033B
	k2         = 0x62992FC1
	k3         = 0x30BC5B29
	cBlockSize = 32
)

// MetroHash64 implements the 64-bit variant of the metrohash algorithm.
// MetroHash64 implements hash.Hash and hash.Hash64 interfaces.
type MetroHash64 struct {
	seed    uint64           // Seed for the algorithm
	state   [4]uint64        // Internal state variables
	mem     [cBlockSize]byte // Buffer to store input less than 32 bytes
	memsize uint32           // Number of unprocessed elements in buffer
	len     uint64           // Total length of the input processed (in bytes)
}

// NewSeedMetroHash64 returns an instance of MetroHash64 with the specified seed.
func NewSeedMetroHash64(seed uint64) *MetroHash64 {
	m := &MetroHash64{
		seed: seed,
	}
	m.Reset()
	return m
}

// NewMetroHash64 returns an instance of MetroHash64 with seed set to 0.
func NewMetroHash64() *MetroHash64 {
	return NewSeedMetroHash64(0)
}

// String returns the current value of the hash as a hexadecimal string
func (m *MetroHash64) String() string {
	return hex.EncodeToString(m.Sum(nil))
}

// Uint64 returns the current value of the hash as an uint64
func (m *MetroHash64) Uint64() uint64 {
	return binary.BigEndian.Uint64(m.Sum(nil))
}

//
// Implement hash.Hash interface
//

// Sum appends the current has to b and returns the resulting slice.
// It does not change the underlying hash state.
func (m *MetroHash64) Sum(b []byte) []byte {
	v := m.Sum64()
	return append(b, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32),
		byte(v>>24), byte(v>>16), byte(v>>8), byte(v))

}

// Reset resets the Hash to its initial state.
func (m *MetroHash64) Reset() {
	vseed := (m.seed + k2) * k0
	m.state = [...]uint64{vseed, vseed, vseed, vseed}
	m.memsize = 0
	m.len = 0
}

// Size returns the number of bytes Sum will return.
func (m *MetroHash64) Size() int {
	return 8
}

// BlockSize returns the hash's underlying block size.
func (m *MetroHash64) BlockSize() int {
	return cBlockSize
}

// Write adds more data to the running hash.
// It never returns an error.
func (m *MetroHash64) Write(input []byte) (int, error) {
	l := len(input)
	m.len += uint64(l)

	if m.memsize+uint32(l) < cBlockSize {
		// New data fits into the buffer
		m.memsize += uint32(copy(m.mem[m.memsize:], input))
		return l, nil
	}

	if m.memsize > 0 {
		// new data does not fit into the buffer
		// and some data is still unprocessed from previous update
		n := cBlockSize - m.memsize
		copy(m.mem[m.memsize:], input[:n])

		m.state[0] += binary.LittleEndian.Uint64(m.mem[:8:8]) * k0
		m.state[0] = rotr64_29(m.state[0]) + m.state[2]

		m.state[1] += binary.LittleEndian.Uint64(m.mem[8:16:16]) * k1
		m.state[1] = rotr64_29(m.state[1]) + m.state[3]

		m.state[2] += binary.LittleEndian.Uint64(m.mem[16:24:24]) * k2
		m.state[2] = rotr64_29(m.state[2]) + m.state[0]

		m.state[3] += binary.LittleEndian.Uint64(m.mem[24:32:32]) * k3
		m.state[3] = rotr64_29(m.state[3]) + m.state[1]

		input = input[n:len(input):len(input)]
		m.memsize = 0
	}

	if len(input) >= cBlockSize {
		for len(input) >= cBlockSize {
			m.state[0] += binary.LittleEndian.Uint64(input[:8:8]) * k0
			m.state[0] = rotr64_29(m.state[0]) + m.state[2]

			m.state[1] += binary.LittleEndian.Uint64(input[8:16:16]) * k1
			m.state[1] = rotr64_29(m.state[1]) + m.state[3]

			m.state[2] += binary.LittleEndian.Uint64(input[16:24:24]) * k2
			m.state[2] = rotr64_29(m.state[2]) + m.state[0]

			m.state[3] += binary.LittleEndian.Uint64(input[24:32:32]) * k3
			m.state[3] = rotr64_29(m.state[3]) + m.state[1]

			input = input[cBlockSize:len(input):len(input)]
		}
	}

	if len(input) > 0 {
		m.memsize += uint32(copy(m.mem[m.memsize:], input))
	}

	return l, nil

}

//
// Implement Hash64 interface
//

// Sum64 returns the current hash value
func (m *MetroHash64) Sum64() uint64 {
	v0, v1, v2, v3 := m.state[0], m.state[1], m.state[2], m.state[3]

	if m.len >= cBlockSize {
		v2 ^= rotr64_37((v0+v3)*k0+v1) * k1
		v3 ^= rotr64_37((v1+v2)*k1+v0) * k0
		v0 ^= rotr64_37((v0+v2)*k0+v3) * k1
		v1 ^= rotr64_37((v1+v3)*k1+v2) * k0

		v0 = (m.seed+k2)*k0 + (v0 ^ v1)
	}

	// Process any bytes remaining in the mem
	if m.memsize > 0 {
		in := m.mem[:m.memsize:m.memsize]
		memsize := m.memsize
		if memsize >= 16 {
			v1 = v0 + binary.LittleEndian.Uint64(in[:8:8])*k2
			v1 = rotr64_29(v1) * k3

			v2 = v0 + binary.LittleEndian.Uint64(in[8:16:16])*k2
			v2 = rotr64_29(v2) * k3

			v1 ^= rotr64_21(v1*k0) + v2
			v2 ^= rotr64_21(v2*k3) + v1
			v0 += v2

			in = in[16:memsize:memsize]
			memsize -= 16
		}

		if memsize >= 8 {
			v0 += binary.LittleEndian.Uint64(in[:8:8]) * k3
			v0 ^= rotr64_55(v0) * k1

			in = in[8:memsize:memsize]
			memsize -= 8
		}

		if memsize >= 4 {
			v0 += uint64(binary.LittleEndian.Uint32(in[:4:4])) * k3
			v0 ^= rotr64_26(v0) * k1

			in = in[4:memsize:memsize]
			memsize -= 4
		}

		if memsize >= 2 {
			v0 += uint64(binary.LittleEndian.Uint16(in[:2:2])) * k3
			v0 ^= rotr64_48(v0) * k1

			in = in[2:memsize:memsize]
			memsize -= 2
		}

		if memsize >= 1 {
			v0 += uint64(uint8(in[0])) * k3
			v0 ^= rotr64_37(v0) * k1
		}
	}

	v0 ^= rotr64_28(v0)
	v0 *= k0
	v0 ^= rotr64_29(v0)

	return v0
}
