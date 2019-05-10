/*
 * Copyright (c) 2015, Robert Bieber
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 *
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 *
 * 2. Redistributions in binary form must reproduce the above
 *    copyright notice, this list of conditions and the following
 *    disclaimer in the documentation and/or other materials provided
 *    with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
 * FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE
 * COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
 * INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
 * STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED
 * OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package decode

import (
	"errors"
)

// #include <zbar.h>
import "C"

var NoMemoryError = errors.New("zbar: Out of memory")
var InternalError = errors.New("zbar: Internal library error")
var UnsupportedError = errors.New("zbar: Unsupported request")
var InvalidError = errors.New("zbar: Invalid request")
var SystemError = errors.New("zbar: System error")
var LockingError = errors.New("zbar: Locking error")
var BusyError = errors.New("zbar: All resources busy")
var XDisplayError = errors.New("zbar: X11 display error")
var XProtoError = errors.New("zbar: X11 protocol error")
var ClosedError = errors.New("zbar: Output window is closed")
var WinAPIError = errors.New("zbar: Windows system error")

var zbarCodeToError = map[int]error{
	C.ZBAR_ERR_NOMEM:       NoMemoryError,
	C.ZBAR_ERR_INTERNAL:    InternalError,
	C.ZBAR_ERR_UNSUPPORTED: UnsupportedError,
	C.ZBAR_ERR_INVALID:     InvalidError,
	C.ZBAR_ERR_SYSTEM:      SystemError,
	C.ZBAR_ERR_LOCKING:     LockingError,
	C.ZBAR_ERR_BUSY:        BusyError,
	C.ZBAR_ERR_XDISPLAY:    XDisplayError,
	C.ZBAR_ERR_XPROTO:      XProtoError,
	C.ZBAR_ERR_CLOSED:      ClosedError,
	C.ZBAR_ERR_WINAPI:      WinAPIError,
}

func errorCodeToError(errorCode int) error {
	if errorCode == 0 {
		return nil
	} else if errorCode >= C.ZBAR_ERR_NUM {
		return errors.New("zbar: Unknown error code from zbar library")
	} else {
		return zbarCodeToError[errorCode]
	}
}
