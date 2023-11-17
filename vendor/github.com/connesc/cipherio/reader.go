package cipherio

import (
	"crypto/cipher"
	"io"
)

type blockReader struct {
	src       io.Reader
	blockMode cipher.BlockMode
	padding   Padding
	blockSize int
	buf       []byte // used to store remaining bytes (before or after crypting)
	crypted   int    // if > 0, then buf contains remaining crypted bytes
	err       error
}

// NewBlockReader wraps the given Reader to add on-the-fly encryption or decryption using the
// given BlockMode.
//
// Data must be aligned to the cipher block size: ErrUnexpectedEOF is returned if EOF is reached in
// the middle of a block.
//
// This Reader avoids buffering and copies as much as possible. A call to Read leads to at most
// one Read from the wrapped Reader. Unless the destination buffer is smaller than BlockSize,
// (en|de)cryption happens inplace within it.
//
// There is no dynamic allocation: an internal buffer of BlockSize bytes is used to store both
// incomplete blocks (not yet (en|de)crypted) and partially read blocks (already (en|de)crypted).
//
// The wrapped Reader is guaranteed to never be consumed beyond the last requested block. This
// means that it is safe to stop reading from this Reader at a block boundary and then resume
// reading from the wrapped Reader for another purpose.
func NewBlockReader(src io.Reader, blockMode cipher.BlockMode) io.Reader {
	return NewBlockReaderWithPadding(src, blockMode, nil)
}

// NewBlockReaderWithPadding is similar to NewBlockReader, except that any incomplete block is
// filled with the given padding instead of returning ErrUnexpectedEOF.
func NewBlockReaderWithPadding(src io.Reader, blockMode cipher.BlockMode, padding Padding) io.Reader {
	blockSize := blockMode.BlockSize()

	return &blockReader{
		src:       src,
		blockMode: blockMode,
		padding:   padding,
		blockSize: blockSize,
		buf:       make([]byte, 0, blockSize),
		crypted:   0,
		err:       nil,
	}
}

func (r *blockReader) readCryptedBuf(p []byte) int {
	n := copy(p, r.buf[r.blockSize-r.crypted:])
	r.crypted -= n
	return n
}

func (r *blockReader) Read(p []byte) (int, error) {
	count := 0

	// Read previously crypted bytes, even if an error has already been encountered. Stop early if
	// the crypted buffer cannot be entirely consumed.
	if r.crypted > 0 {
		n := r.readCryptedBuf(p)
		p = p[n:]
		count += n
		if r.crypted > 0 {
			return count, nil
		}
		r.buf = r.buf[:0]
	}
	// At this point, the internal buffer cannot contain crypted bytes anymore.

	// Return the previously saved error, if any.
	if r.err != nil {
		return count, r.err
	}

	// Stop early if there is no more space in the destination buffer.
	if len(p) == 0 {
		return count, nil
	}

	// If the destination buffer is smaller than BlockSize, then use the internal buffer.
	if len(p) < r.blockSize {
		// The internal buffer may already contain some bytes, try to fill the rest with a single
		// Read.
		n, err := r.src.Read(r.buf[len(r.buf):r.blockSize])
		r.buf = r.buf[:len(r.buf)+n]

		// Apply padding if EOF is reached in the middle of a block.
		if err == io.EOF && len(r.buf) < r.blockSize && r.padding != nil {
			r.padding.Fill(r.buf[len(r.buf):r.blockSize])
			r.buf = r.buf[:r.blockSize]
		}

		// Crypt the buffered block if complete, then fill the destination buffer with the first
		// crypted bytes.
		if len(r.buf) == r.blockSize {
			r.blockMode.CryptBlocks(r.buf, r.buf)
			r.crypted = r.blockSize
			count += r.readCryptedBuf(p)
		}

		// Save any encountered error.
		r.err = err

		if r.crypted > 0 {
			// Hide any error until crypted bytes have been entirely consumed.
			err = nil
		} else if err == io.EOF && len(r.buf) > 0 {
			// If EOF is reached in the middle of a block, convert it to ErrUnexpectedEOF.
			err = io.ErrUnexpectedEOF
			r.err = err
		}
		return count, err
	}
	// Otherwise, use the destination buffer.

	// Initialize the destination buffer with buffered bytes, then try to fill the rest with a
	// single Read.
	copy(p, r.buf)
	n, err := r.src.Read(p[len(r.buf):])
	available := len(r.buf) + n
	exceeding := available % r.blockSize
	cryptable := available - exceeding

	// Crypt all complete blocks.
	if cryptable > 0 {
		r.blockMode.CryptBlocks(p[:cryptable], p[:cryptable])
		p = p[cryptable:]
		count += cryptable
	}

	// Store exceeding bytes to the internal buffer.
	r.buf = r.buf[:exceeding]
	copy(r.buf, p)
	// At this point, both the destination and the internal buffers contain the exceeding bytes.

	// Save any encountered error.
	r.err = err

	// Handle EOF when encountered in the middle of a block.
	if err == io.EOF && exceeding > 0 {
		if r.padding == nil {
			// If no padding is defined, convert EOF to ErrUnexpectedEOF.
			err = io.ErrUnexpectedEOF
			r.err = err

		} else if len(p) < r.blockSize {
			// If padding does not fit the destination buffer, then use the internal buffer.
			r.padding.Fill(r.buf[exceeding:r.blockSize])
			r.buf = r.buf[:r.blockSize]

			// Crypt the padded block, then fill the rest of the destination buffer with the first
			// crypted bytes.
			r.blockMode.CryptBlocks(r.buf, r.buf)
			r.crypted = r.blockSize
			count += r.readCryptedBuf(p)

			// Hide any error until crypted bytes have been entirely consumed.
			if r.crypted > 0 {
				err = nil
			}

		} else {
			// Otherwise, apply padding to the destination buffer and crypt the padded block.
			r.padding.Fill(p[exceeding:r.blockSize])
			r.buf = r.buf[:0]
			r.blockMode.CryptBlocks(p[:r.blockSize], p[:r.blockSize])
			count += r.blockSize
		}
	}

	return count, err
}
