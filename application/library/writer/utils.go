package writer

import (
	"io"
	"os"
)

func GetOutAndErr(writer ...io.Writer) (wOut, wErr io.Writer) {
	length := len(writer)
	wOut, wErr = os.Stdout, os.Stderr
	if length > 0 {
		if writer[0] != nil {
			wOut = writer[0]
		}
		if length > 1 {
			if writer[1] != nil {
				wErr = writer[1]
			}
		}
	}
	return
}
