package internal

import (
	"fmt"
	"html"
	"log"
	"os"
	"time"
)

type Statics struct {
	timer     *myTimer
	tables    []*tableStatics
	Config    *Config
	failedNum int
	result    *string
	diff      bool
}

type tableStatics struct {
	timer       *myTimer
	table       string
	alter       *TableAlterData
	alterRet    error
	schemaAfter string
}

func newStatics(cfg *Config) *Statics {
	return &Statics{
		timer:     newMyTimer(),
		tables:    make([]*tableStatics, 0),
		Config:    cfg,
		failedNum: -1,
		diff:      true,
	}
}

func (s *Statics) newTableStatics(table string, sd *TableAlterData) *tableStatics {
	ts := &tableStatics{
		timer: newMyTimer(),
		table: table,
		alter: sd,
	}
	if sd.Type != alterTypeNo {
		s.tables = append(s.tables, ts)
	}
	return ts
}

func (s *Statics) String() string {
	if s.result != nil {
		return *s.result
	}
	result := s.toHTML()
	s.result = &result
	return result
}

func (s *Statics) ChangeTableNum() int {
	return len(s.tables)
}

func (s *Statics) ChangeTables() []string {
	tables := make([]string, len(s.tables))
	for k, t := range s.tables {
		tables[k] = t.table
	}
	return tables
}

func (s *Statics) FailedNum() int {
	return s.alterFailedNum()
}

func (s *Statics) StartTime() time.Time {
	return s.timer.start
}

func (s *Statics) Diff(on bool) *Statics {
	s.diff = on
	return s
}

func (s *Statics) EndTime() time.Time {
	return s.timer.end
}

func (s *Statics) Elapsed() time.Duration {
	return s.timer.end.Sub(s.timer.start)
}

func (s *Statics) toHTML() string {
	cfg := s.Config
	hostName, _ := os.Hostname()
	code := "<h2>Info</h2>\n<pre>"
	code += "  from : " + dsnSort(cfg.SourceDSN) + "\n"
	code += "    to : " + dsnSort(cfg.DestDSN) + "\n"
	code += " alter : " + fmt.Sprintf("%d", len(s.tables)) + " tables\n"
	code += "<font color=green>  sync : " + fmt.Sprintf("%t", s.Config.Sync) + "</font>\n"
	if s.Config.Sync {
		fn := s.alterFailedNum()
		code += "<font color=red>failed : " + fmt.Sprintf("%d", fn) + "</font>\n"
	}
	code += "\n"
	code += "  host : " + hostName + "\n"
	code += " start : " + s.timer.start.Format(timeFormatStd) + "\n"
	code += "   end : " + s.timer.end.Format(timeFormatStd) + "\n"
	code += "  used : " + s.timer.usedSecond() + "\n"

	code += "</pre>\n"
	code += "<h2>Result</h2>\n"
	code += "<h3>Tables</h3>\n"
	code += `<table class='table table-bordered tb_1'>
		<thead>
			<tr>
			<th width="60px">no</th>
			<th>table</th>
			<th>alter result</th>
			<th>used</th>
			</tr>
		</thead><tbody>
		`
	for idx, tb := range s.tables {
		code += "<tr>"
		code += "<td>" + fmt.Sprintf("%d", idx+1) + "</td>\n"
		code += "<td><a href='#detail_" + tb.table + "'>" + tb.table + "</a></td>\n"
		code += "<td>"
		if s.Config.Sync {
			if tb.alterRet == nil {
				code += "<font color=green>success</font>"
			} else {
				code += "<font color=red>failed," + html.EscapeString(tb.alterRet.Error()) + "</font>"
			}
		} else {
			code += "not sync"
		}
		code += "</td>\n"
		code += "<td>" + tb.timer.usedSecond() + "</td>\n"
		code += "</tr>\n"
	}
	code += "</tbody></table>\n<h3>Sqls</h3>\n<pre>"
	for _, tb := range s.tables {
		code += "<a name='detail_" + tb.table + "'></a>"
		code += html.EscapeString(tb.alter.String()) + "\n\n"
	}
	code += "</pre>\n\n"
	if !s.diff {
		return code
	}
	code += "<h3>Detail</h3>\n"
	code += `<table class='table table-bordered tb_1'>
		<thead>
			<tr>
			<th width="40px">no</th>
			<th width="80px">table</th>
			<th>source</th>
			<th>destination</th>
			</tr>
		</thead><tbody>
		`
	for idx, tb := range s.tables {
		code += "<tr>"
		code += "<th rowspan=2>" + fmt.Sprintf("%d", idx+1) + "</th>\n"
		code += "<td rowspan=2>" + tb.table + "<br/><br/>"
		if s.Config.Sync {
			if tb.alterRet == nil {
				code += "<font color=green>success</font>"
			} else {
				code += "<font color=red>failed," + tb.alterRet.Error() + "</font>"
			}
		} else {
			code += "(no sync)"
		}
		code += "</td>\n"
		code += "<td valign=top><b>source schema:</b><br/>" + htmlPre(tb.alter.SchemaDiff.Source.SchemaRaw) + "</td>\n"
		code += "<td valign=top><b>dest schema:</b><br/>" + htmlPre(tb.alter.SchemaDiff.Dest.SchemaRaw) + "</td>\n"
		code += "</tr>\n"

		code += "<tr>\n"
		code += "<td valign=top><b>alter:</b><br/>" + htmlPre(tb.alter.SQL) + "</td>\n"
		code += "<td valign=top><b>alter after:</b><br/>" + htmlPre(tb.schemaAfter) + "</td>\n"
		code += "</tr>\n"
	}
	code += "</tbody></table>\n"
	return code
}

func (s *Statics) alterFailedNum() int {
	if s.failedNum > -1 {
		return s.failedNum
	}
	n := 0
	for _, tb := range s.tables {
		if tb.alterRet != nil {
			n++
		}
	}
	s.failedNum = n
	return n
}

func (s *Statics) sendMailNotice() {
	cfg := s.Config
	if cfg.Email == nil {
		log.Println("mail conf is not set,skip send mail")
		return
	}
	alterTotal := len(s.tables)
	if alterTotal < 1 {
		log.Println("no table change,skip send mail")
		return
	}
	title := "[mysql_schema_sync] " + fmt.Sprintf("%d", alterTotal) + " tables change [" + dsnSort(cfg.DestDSN) + "]"
	var body string

	if !s.Config.Sync {
		title += "[preview]"
		body += "<font color=red>this is preview,all sql never execute!</font>\n"
	} else {
		fn := s.alterFailedNum()
		if fn > 0 {
			title += " [failed=" + fmt.Sprintf("%d", fn) + "]"
		}
	}
	body += s.toHTML()
	cfg.Email.SendMail(title, body)
}
