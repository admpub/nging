package nftables

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/webx-top/com"
)

type RowInfo struct {
	RowNo  uint64
	Handle *uint64
	Row    string
}

func (r RowInfo) HasHandleID() bool {
	return r.Handle != nil
}

func (r RowInfo) GetHandleID() uint64 {
	if r.Handle == nil {
		return 0
	}
	return *r.Handle
}

func (r RowInfo) String() string {
	return r.Row
}

func ListPage(r io.Reader, page, limit uint) (rows []RowInfo, hasMore bool, err error) {
	offset := uint64(com.Offset(page, limit))
	maxNo := offset + uint64(limit)
	s := bufio.NewScanner(r)
	var i uint64
	for s.Scan() {
		if offset > i {
			continue
		}
		if i >= maxNo {
			hasMore = true
			return
		}
		t := s.Text()
		t = strings.TrimSpace(t)
		if strings.HasSuffix(t, `{`) || t == `}` {
			continue
		}
		var rowInfo *RowInfo
		parts := strings.SplitN(t, `# handle `, 2)
		if len(parts) == 2 {
			parts[0] = strings.TrimSpace(parts[0])
			if strings.HasSuffix(parts[0], `{`) {
				continue
			}
			var handleID uint64
			handleID, err = strconv.ParseUint(parts[1], 10, 64)
			if err != nil {
				return
			}
			rowInfo = &RowInfo{
				RowNo:  i,
				Row:    parts[0],
				Handle: &handleID,
			}
		} else {
			rowInfo = &RowInfo{
				RowNo: i,
				Row:   t,
			}
		}
		rows = append(rows, *rowInfo)
		i++
	}
	return
}
