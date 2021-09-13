package common

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/abh/errorutil"
)

func JSONBytesParseError(err error, jsonBytes []byte) error {
	var offset int64 = -1
	switch rErr := err.(type) {
	case *json.UnmarshalTypeError:
		offset = rErr.Offset
	case *json.SyntaxError:
		offset = rErr.Offset
	}
	if offset > -1 {
		byteReader := bytes.NewReader(jsonBytes)
		line, col, highlight := errorutil.HighlightBytePosition(byteReader, offset)
		extra := fmt.Sprintf(":\nError at line %d, column %d (offset %d):\n%s",
			line, col, offset, highlight)
		return fmt.Errorf("error parsing json object%s\n%w",
			extra, err)
	}
	return fmt.Errorf("%w: %s", err, string(jsonBytes))
}
