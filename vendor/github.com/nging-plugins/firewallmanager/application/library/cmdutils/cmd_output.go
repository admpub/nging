/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package cmdutils

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os/exec"
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

type readData struct {
	rows    []RowInfo
	hasMore bool
	err     error
}

func LineSeeker(r io.Reader, page, limit uint, parser func(uint64, string) (*RowInfo, error)) (rows []RowInfo, hasMore bool, err error) {
	offset := uint64(com.Offset(page, limit))
	maxNo := offset + uint64(limit)
	var i uint64
	s := bufio.NewScanner(r)
	for s.Scan() {
		if offset > i {
			continue
		}
		if i >= maxNo {
			hasMore = true
			err = ErrCmdForcedExit
			return
		}
		t := s.Text()
		t = strings.TrimSpace(t)
		var rowInfo *RowInfo
		rowInfo, err = parser(i, t)
		if err != nil {
			return
		}
		if rowInfo == nil {
			continue
		}
		rows = append(rows, *rowInfo)
		i++
	}
	return
}

func RecvCmdOutputs(page, limit uint,
	cmdBin string, cmdArgs []string,
	parser func(uint64, string) (*RowInfo, error),
) (rows []RowInfo, hasMore bool, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res := make(chan *readData)
	err = RunCmdWithCallback(ctx, cmdBin, cmdArgs, func(cmd *exec.Cmd) error {
		r, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}
		go func() {
			rd := &readData{}
			rd.rows, rd.hasMore, rd.err = LineSeeker(r, page, limit, parser)
			if rd.err != nil {
				cancel()
			}
			res <- rd
			r.Close()
		}()
		return nil
	})
	rd := <-res
	rows = rd.rows
	hasMore = rd.hasMore
	if rd.err == nil {
		return
	}
	if !errors.Is(rd.err, ErrCmdForcedExit) {
		err = rd.err
		return
	}
	if err != nil && err.Error() == `signal: killed` {
		err = nil
	}
	return
}
