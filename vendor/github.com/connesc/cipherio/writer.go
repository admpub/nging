package cipherio

import (
	"crypto/cipher"
	"io"
)

type blockWriter struct {
	dst       io.Writer
	blockMode cipher.BlockMode
	padding   Padding
	blockSize int
	buf       []byte // used to store both incomplete and crypted blocks
	err       error
}

// NewBlockWriter wraps the given Writer to add on-the-fly encryption or decryption using the
// given BlockMode.
//
// Data must be aligned to the cipher block size: ErrUnexpectedEOF is returned if Close is called
// in the middle of a block.
//
// This Writer allocates an internal buffer of 1024 blocks, which is freed when an error is
// encountered or when Close is called. Other than that, there is no dynamic allocation.
//
// Close must be called at least once. After that, Close becomes a no-op and Write must not be
// called anymore.
func NewBlockWriter(dst io.Writer, blockMode cipher.BlockMode) io.WriteCloser {
	return NewBlockWriterWithPadding(dst, blockMode, nil)
}

// NewBlockWriterWithPadding is similar to NewBlockWriter, except that Close fills any incomplete
// block with the given padding instead of returning ErrUnexpectedEOF.
func NewBlockWriterWithPadding(dst io.Writer, blockMode cipher.BlockMode, padding Padding) io.WriteCloser {
	blockSize := blockMode.BlockSize()

	return &blockWriter{
		dst:       dst,
		blockMode: blockMode,
		padding:   padding,
		blockSize: blockSize,
		buf:       make([]byte, 0, 1024*blockSize),
		err:       nil,
	}
}

func (w *blockWriter) Write(p []byte) (int, error) {
	count := 0

	// Return the previously saved error, if any.
	if w.err != nil {
		return count, w.err
	}

	// While complete blocks are available, crypt as many as possible in the internal buffer and
	// write the result to the destination writer.
	for len(w.buf)+len(p) >= w.blockSize {
		// Initialize src with remaining bytes.
		src := w.buf
		remaining := len(src)

		// Clear remaining bytes.
		w.buf = w.buf[:0]

		// If src contains an incomplete block then fill it and crypt it inplace.
		if remaining > 0 {
			src = src[:w.blockSize]

			copied := copy(src[remaining:], p)
			p = p[copied:]

			w.blockMode.CryptBlocks(src, src)
		}

		// Otherwise, determine how many complete blocks can be stored in src.
		cryptable := cap(src) - len(src)
		if len(p) < cryptable {
			cryptable = (len(p) / w.blockSize) * w.blockSize
		}

		// If any, crypt them and store the result in src at the same time. This avoids a
		// preliminary copy.
		if cryptable > 0 {
			w.blockMode.CryptBlocks(src[len(src):cap(src)], p[:cryptable])
			p = p[cryptable:]
			src = src[:len(src)+cryptable]
		}

		// Now that src is filled with crypted blocks, write them to the destination writer.
		n, err := w.dst.Write(src)

		// Count written bytes, except those that come from the internal buffer, because they have
		// already been aknowledged by the previous call.
		if n > remaining {
			count += n - remaining
		}

		// If any error is encountered, save it, free the internal buffer and stop immediately.
		if err != nil {
			w.err = err
			w.buf = nil
			return count, err
		}
	}

	// If an incomplete block remains, store it in the internal buffer and consider it as written.
	if len(p) > 0 {
		remaining := len(w.buf)
		w.buf = w.buf[:remaining+len(p)]
		copied := copy(w.buf[remaining:], p)
		count += copied
	}

	return count, nil
}

func (w *blockWriter) Close() error {
	// Return the previously saved error, if any.
	if w.err != nil {
		return w.err
	}

	// Initialize src with remaining bytes.
	src := w.buf
	remaining := len(src)

	// Free the internal buffer.
	w.buf = nil

	// Stop early if the internal buffer does not contain an incomplete block.
	if remaining == 0 {
		return nil
	}

	// Return ErrUnexpectedEOF if no padding is defined.
	if w.padding == nil {
		w.err = io.ErrUnexpectedEOF
		return w.err
	}

	// Fill the incomplete block with padding.
	src = src[:w.blockSize]
	w.padding.Fill(src[remaining:])

	// Crypt the last block inplace.
	w.blockMode.CryptBlocks(src, src)

	// Write the last block to the destination writer.
	_, w.err = w.dst.Write(src)
	return w.err
}
