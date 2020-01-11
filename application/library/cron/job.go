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

package cron

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/echo/engine"

	"github.com/admpub/log"
	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/library/charset"
	"github.com/admpub/nging/application/library/email"
)

var (
	defaultOuputSize uint64 = 1024 * 200
	mailTpl          *template.Template
	defaultTmpl      = `
	你好 {{.username}}，<br/>

<p>以下是任务执行结果：</p>

<p>
任务ID：{{.task_id}}<br/>
任务名称：{{.task_name}}<br/>       
执行时间：{{.start_time}}<br />
执行耗时：{{.process_time}}秒<br />
执行状态：{{.status}}
</p>
<p>-------------以下是任务执行输出-------------</p>
<p>{{.output}}</p>
<p>
--------------------------------------------<br />
本邮件由系统自动发出，请勿回复<br />
如果要取消邮件通知，请登录到系统进行设置<br />
</p>
`
	cmdPreParams []string
	SYSJobs      = map[string]Jobx{}
	ErrFailure   = errors.New(`Error`)
	//NotRecordPrefixFlag 不记录日志的前缀标识
	NotRecordPrefixFlag = `--/ignore/--`
)

func AddSYSJob(name string, fn RunnerGetter, example string, description string) {
	SYSJobs[name] = Jobx{
		Example:      example,
		Description:  description,
		RunnerGetter: fn,
	}
}

type Runner func(timeout time.Duration) (out string, runingErr string, onRunErr error, isTimeout bool)

type RunnerGetter func(string) Runner

type Jobx struct {
	Example      string //">funcName:param"
	Description  string
	RunnerGetter RunnerGetter
}

func init() {
	if com.IsWindows {
		cmdPreParams = []string{"cmd.exe", "/c"}
		//cmdPreParams = []string{"bash.exe", "-c"}
	} else {
		shell := os.Getenv("SHELL")
		if len(shell) == 0 {
			shell = "/bin/bash"
		}
		cmdPreParams = []string{shell, "-c"}
	}
}

func CmdParams(command string) []string {
	params := append([]string{}, cmdPreParams...)
	params = append(params, command)
	return params
}

func InitialMailTpl() {
	tmpl := DefaultEmailConfig.Template
	if len(tmpl) == 0 {
		tmpl = defaultTmpl
	}
	var err error
	mailTpl, err = template.New("notifyMailTmpl").Parse(tmpl)
	if err != nil {
		panic(err)
	}
}

func MailTpl() *template.Template {
	if mailTpl == nil {
		InitialMailTpl()
	}
	return mailTpl
}

type Job struct {
	id         uint              // 任务ID
	logID      uint64            // 日志记录ID
	name       string            // 任务名称
	task       *dbschema.Task    // 任务对象
	taskLog    *dbschema.TaskLog // 结果日志
	runner     Runner            // 执行函数
	isSYS      bool              // 是否是系统内部功能
	status     int64             // 任务状态，大于0表示正在执行中
	Concurrent bool              // 同一个任务是否允许并行执行
}

func NewJobFromTask(ctx context.Context, task *dbschema.Task) (*Job, error) {
	if task.Id < 1 {
		return nil, fmt.Errorf("Job: missing task.Id")
	}
	var env []string
	task.Env = strings.TrimSpace(task.Env)
	if len(task.Env) > 0 {
		for _, row := range strings.Split(task.Env, "\n") {
			row = strings.TrimSpace(row)
			if len(row) > 0 {
				env = append(env, row)
			}
		}
	}
	cmd := task.Command
	if len(cmd) > 0 && cmd[0] == '>' {
		cmd = cmd[1:]
		cmdInfo := strings.SplitN(cmd, `:`, 2)
		var param string
		switch len(cmdInfo) {
		case 2:
			param = cmdInfo[1]
			fallthrough
		case 1:
			fnName := cmdInfo[0]
			if jobx, ok := SYSJobs[fnName]; ok {
				job := &Job{
					id:         task.Id,
					name:       task.Name,
					task:       task,
					Concurrent: task.Concurrent == 1,
					runner:     jobx.RunnerGetter(param),
					isSYS:      true,
				}
				return job, nil
			}
		}
	}
	if task.GroupId > 0 {
		group := &dbschema.TaskGroup{}
		err := group.Get(nil, `id`, task.GroupId)
		if err != nil {
			return nil, err
		}
		if len(group.CmdPrefix) > 0 {
			cmd = group.CmdPrefix + ` ` + cmd
		}
		if len(group.CmdSuffix) > 0 {
			cmd += ` ` + group.CmdSuffix
		}
	}
	job := NewCommandJob(ctx, task.Id, task.Name, cmd, task.WorkDirectory, env...)
	job.task = task
	job.Concurrent = task.Concurrent == 1
	return job, nil
}

func NewOutputWriter(sizes ...uint64) OutputWriter {
	var size uint64
	if len(sizes) > 0 {
		size = sizes[0]
	}
	if size == 0 {
		size = defaultOuputSize
	}
	return NewCmdRec(size)
}

func NewCommandJob(ctx context.Context, id uint, name string, command string, dir string, env ...string) *Job {
	job := &Job{
		id:   id,
		name: name,
	}
	job.runner = func(timeout time.Duration) (string, string, error, bool) {
		if ctx == nil {
			ctx = context.Background()
		}
		bufOut := NewCmdRec(defaultOuputSize)
		bufErr := NewCmdRec(defaultOuputSize)
		params := CmdParams(command)
		cmd := exec.Command(params[0], params[1:]...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), env...)
		cmd.Stdout = bufOut
		cmd.Stderr = bufErr
		cmd.Start()
		err, isTimeout := runCmdWithTimeout(cmd, timeout, ctx)
		if com.IsWindows {
			bOut, e := charset.Convert(`gbk`, `utf-8`, bufOut.Bytes())
			if e != nil {
				log.Error(e)
			}
			bErr, e := charset.Convert(`gbk`, `utf-8`, bufErr.Bytes())
			if e != nil {
				log.Error(e)
			}
			return engine.Bytes2str(bOut), engine.Bytes2str(bErr), err, isTimeout
		}

		return bufOut.String(), bufErr.String(), err, isTimeout
	}
	return job
}

func (j *Job) Status() int64 {
	return atomic.LoadInt64(&j.status)
}

func (j *Job) Name() string {
	return j.name
}

func (j *Job) Id() uint {
	return j.id
}

func (j *Job) LogID() uint64 {
	return j.logID
}

func (j *Job) LogData() *dbschema.TaskLog {
	return j.taskLog
}

func (j *Job) addAndReturningLog() *Job {
	// 插入日志
	_, err := j.taskLog.Add()
	if err != nil {
		log.Error("Job: 日志写入失败: ", err)
	}
	j.logID = j.taskLog.Id
	return j
}

func (j *Job) sendEmail(elapsed float64, t time.Time, err error, cmdOut string, isTimeout bool, timeout time.Duration) {
	user := new(dbschema.User)
	uerr := user.Get(nil, `id`, j.task.Uid)
	if uerr != nil {
		return
	}

	var title string

	data := make(map[string]interface{})
	data["task_id"] = j.task.Id
	data["username"] = user.Username
	data["task_name"] = j.task.Name
	data["start_time"] = t.Format("2006-01-02 15:04:05")
	data["process_time"] = elapsed / 1000
	data["output"] = cmdOut

	if isTimeout {
		title = fmt.Sprintf("任务执行结果通知 #%d: %s", j.task.Id, "超时")
		data["status"] = fmt.Sprintf("超时（%d秒）", int(timeout/time.Second))
	} else if err != nil {
		title = fmt.Sprintf("任务执行结果通知 #%d: %s", j.task.Id, "失败")
		data["status"] = "失败（" + err.Error() + "）"
	} else {
		title = fmt.Sprintf("任务执行结果通知 #%d: %s", j.task.Id, "成功")
		data["status"] = "成功"
	}

	content := new(bytes.Buffer)
	MailTpl().Execute(content, data)
	var ccList []string
	if len(j.task.NotifyEmail) > 0 {
		ccList = strings.Split(j.task.NotifyEmail, "\n")
		for index, email := range ccList {
			email = strings.TrimSpace(email)
			if len(email) == 0 {
				continue
			}
			ccList[index] = email
		}
	}
	if err = SendMail(user.Email, user.Username, title, content.Bytes(), ccList...); err != nil {
		log.Error(err)
	}
}

func (j *Job) Run() {
	if !j.Concurrent && atomic.LoadInt64(&j.status) > 0 {
		log.Debugf("任务[ %d. %s ]上一次执行尚未结束，本次被忽略。", j.id, j.name)
		return
	}

	defer func() {
		if err := recover(); err != nil {
			log.Error(err, "\n", string(debug.Stack()))
		}
	}()

	if workPool != nil {
		workPool <- true
		defer func() {
			if workPool == nil {
				return
			}
			<-workPool
		}()
	}

	log.Debugf("开始执行任务: %d", j.id)

	atomic.StoreInt64(&j.status, atomic.LoadInt64(&j.status)+1)
	defer func() {
		atomic.StoreInt64(&j.status, atomic.LoadInt64(&j.status)-1)
	}()

	t := time.Now()
	timeout := time.Duration(time.Hour * 24)
	if j.task.Timeout > 0 {
		timeout = time.Second * time.Duration(j.task.Timeout)
	}

	cmdOut, cmdErr, err, isTimeout := j.runner(timeout)

	elapsed := time.Now().Sub(t).Seconds()

	tl := new(dbschema.TaskLog)
	tl.TaskId = j.id
	tl.Output = cmdOut
	tl.Error = cmdErr
	tl.Created = uint(t.Unix())
	tl.Elapsed = uint(elapsed)
	if isTimeout {
		tl.Status = `timeout`
		tl.Error = fmt.Sprintf("任务执行超过 %d 秒\n----------------------\n%s\n", int(timeout/time.Second), tl.Error)
	} else if err != nil {
		tl.Status = `failure`
		tl.Error = err.Error() + ":" + tl.Error
	} else {
		tl.Status = `success`
	}

	j.taskLog = tl

	if j.task.ClosedLog == `N` && !strings.HasPrefix(cmdOut, NotRecordPrefixFlag) && !strings.HasPrefix(cmdErr, NotRecordPrefixFlag) {
		j.addAndReturningLog()
	}

	// 更新上次执行时间
	j.task.PrevTime = uint(t.Unix())
	j.task.ExecuteTimes++
	setErr := j.task.SetFields(nil, map[string]interface{}{
		`prev_time`:     j.task.PrevTime,
		`execute_times`: j.task.ExecuteTimes,
	}, `id`, j.task.Id)
	if setErr != nil {
		log.Error(setErr)
	}

	// 发送邮件通知
	if (j.task.EnableNotify == 1 && err != nil) || j.task.EnableNotify == 2 {
		j.sendEmail(elapsed, t, err, cmdOut, isTimeout, timeout)
	}
}

func SendMail(toEmail string, toUsername string, title string, content []byte, ccList ...string) error {
	return SendMailWithID(0, toEmail, toUsername, title, content, ccList...)
}

func SendMailWithID(id uint64, toEmail string, toUsername string, title string, content []byte, ccList ...string) error {
	if len(toEmail) < 1 {
		//收信人邮箱地址不正确
		return ErrIncorrectRecipient
	}
	conf := &email.Config{
		ID:         id,
		Engine:     DefaultEmailConfig.Engine,
		SMTP:       DefaultSMTPConfig,
		From:       DefaultEmailConfig.Sender,
		ToAddress:  toEmail,
		ToUsername: toUsername,
		Subject:    title,
		Content:    content,
		CcAddress:  ccList,
		Timeout:    DefaultEmailConfig.Timeout,
	}
	return email.SendMail(conf)
}
