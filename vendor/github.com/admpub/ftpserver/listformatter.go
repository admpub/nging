package server

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type listFormatter []FileInfo

// Short returns a string that lists the collection of files by name only,
// one per line
func (formatter listFormatter) Short() []byte {
	var buf bytes.Buffer
	for _, file := range formatter {
		fmt.Fprintf(&buf, "%s\r\n", file.Name())
	}
	fmt.Fprintf(&buf, "\r\n")
	return buf.Bytes()
}

// Detailed returns a string that lists the collection of files with extra
// detail, one per line
func (formatter listFormatter) Detailed() []byte {
	var buf bytes.Buffer
	for _, file := range formatter {
		fmt.Fprintf(&buf, file.Mode().String())
		fmt.Fprintf(&buf, " 1 %s %s ", file.Owner(), file.Group())
		fmt.Fprintf(&buf, lpad(strconv.Itoa(int(file.Size())), 12))
		fmt.Fprintf(&buf, file.ModTime().Format(" Jan _2 15:04 "))
		fmt.Fprintf(&buf, "%s\r\n", file.Name())
	}
	fmt.Fprintf(&buf, "\r\n")
	return buf.Bytes()
}

func lpad(input string, length int) (result string) {
	if len(input) < length {
		result = strings.Repeat(" ", length-len(input)) + input
	} else if len(input) == length {
		result = input
	} else {
		result = input[0:length]
	}
	return
}
