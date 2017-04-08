package cron

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"
	"time"

	"runtime"

	"github.com/admpub/log"
	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/library/charset"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/library/email"
	"github.com/webx-top/echo/engine"
)

var (
	isWin            bool
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
)

func init() {
	isWin = runtime.GOOS == `windows`
	if isWin {
		cmdPreParams = []string{"cmd.exe", "/c"}
		//CmdPreParams = []string{"bash.exe", "-c"}
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
	tmpl := config.DefaultConfig.Cron.Template
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
	id         uint                                              // 任务ID
	logId      uint                                              // 日志记录ID
	name       string                                            // 任务名称
	task       *dbschema.Task                                    // 任务对象
	runFunc    func(time.Duration) (string, string, error, bool) // 执行函数
	status     int                                               // 任务状态，大于0表示正在执行中
	Concurrent bool                                              // 同一个任务是否允许并行执行
}

func NewJobFromTask(task *dbschema.Task) (*Job, error) {
	if task.Id < 1 {
		return nil, fmt.Errorf("ToJob: 缺少id")
	}
	job := NewCommandJob(task.Id, task.Name, task.Command)
	job.task = task
	job.Concurrent = task.Concurrent == 1
	return job, nil
}

func NewCommandJob(id uint, name string, command string) *Job {
	job := &Job{
		id:   id,
		name: name,
	}
	job.runFunc = func(timeout time.Duration) (string, string, error, bool) {
		bufOut := NewCmdRec(defaultOuputSize)
		bufErr := NewCmdRec(defaultOuputSize)
		params := CmdParams(command)
		cmd := exec.Command(params[0], params[1:]...)
		cmd.Stdout = bufOut
		cmd.Stderr = bufErr
		cmd.Start()
		err, isTimeout := runCmdWithTimeout(cmd, timeout)
		if isWin {
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

func (j *Job) Status() int {
	return j.status
}

func (j *Job) Name() string {
	return j.name
}

func (j *Job) Id() uint {
	return j.id
}

func (j *Job) LogId() uint {
	return j.logId
}

func (j *Job) Run() {
	if !j.Concurrent && j.status > 0 {
		log.Warnf("任务[%d]上一次执行尚未结束，本次被忽略。", j.id)
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
			<-workPool
		}()
	}

	log.Debugf("开始执行任务: %d", j.id)

	j.status++
	defer func() {
		j.status--
	}()

	t := time.Now()
	timeout := time.Duration(time.Hour * 24)
	if j.task.Timeout > 0 {
		timeout = time.Second * time.Duration(j.task.Timeout)
	}

	cmdOut, cmdErr, err, isTimeout := j.runFunc(timeout)

	ut := time.Now().Sub(t).Seconds()

	// 插入日志
	tl := new(dbschema.TaskLog)
	tl.TaskId = j.id
	tl.Output = cmdOut
	tl.Error = cmdErr
	tl.Created = uint(t.Unix())
	tl.Elapsed = uint(ut)

	if isTimeout {
		tl.Status = `timeout`
		tl.Error = fmt.Sprintf("任务执行超过 %d 秒\n----------------------\n%s\n", int(timeout/time.Second), cmdErr)
	} else if err != nil {
		tl.Status = `failure`
		tl.Error = err.Error() + ":" + cmdErr
	} else {
		tl.Status = `success`
	}
	pk, err2 := tl.Param().SetSend(tl).Insert()
	if err2 != nil {
		log.Error("日志写入失败: ", err2)
	}
	switch id := pk.(type) {
	case uint:
		tl.Id = id
	case int64:
		tl.Id = uint(id)
	default:
		log.Fatalf("无法获取写入日志的主键ID，返回的数据类型为: %T", pk)
	}
	j.logId = tl.Id

	// 更新上次执行时间
	j.task.PrevTime = uint(t.Unix())
	j.task.ExecuteTimes++
	err2 = j.task.Param().SetArgs(`id`, j.task.Id).SetSend(map[string]interface{}{
		`prev_time`:     j.task.PrevTime,
		`execute_times`: j.task.ExecuteTimes,
	}).Update()
	if err2 != nil {
		log.Error(err2)
	}

	// 发送邮件通知
	if (j.task.EnableNotify == 1 && err != nil) || j.task.EnableNotify == 2 {
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
		data["process_time"] = float64(ut) / 1000
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
				ccList[index] = strings.TrimSpace(email)
			}
		}
		if err = SendMail(user.Email, user.Username, title, content.Bytes(), ccList...); err != nil {
			log.Error(err)
		}
	}
}

func SendMail(toEmail string, toUsername string, title string, content []byte, ccList ...string) error {
	conf := &email.Config{
		SMTP:       config.DefaultConfig.Email.SMTPConfig,
		From:       config.DefaultConfig.Email.From,
		ToAddress:  toEmail,
		ToUsername: toUsername,
		Subject:    title,
		Content:    content,
		CcAddress:  ccList,
	}
	return email.SendMail(conf)
}
