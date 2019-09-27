package imgparse

import (
    "encoding/binary"
    "io"
)

func readbuf(r io.Reader, length int) ([]byte, error) {
    buf    := make([]byte, length)
    offset := 0

    for offset < length {
        read, err := r.Read(buf[offset:])
        if err != nil {
            return buf, err
        }

        offset += read
    }

    return buf, nil
}

func read32be(in []byte) uint32 {
    return binary.BigEndian.Uint32(in)
}

func read24le(in []byte) int {
    return int(binary.LittleEndian.Uint32(in) & 0xFFFFFF)
}

func read16le(in []byte) int {
    return int(binary.LittleEndian.Uint16(in))
}

func read16be(in []byte) int {
    return int(binary.BigEndian.Uint16(in))
}
