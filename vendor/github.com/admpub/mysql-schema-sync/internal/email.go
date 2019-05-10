package internal

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

type SendEmail interface {
	SendEmail(string, string)
}

// EmailStruct email conf info
type EmailStruct struct {
	SendMailAble bool   `json:"send_mail"`
	SMTPHost     string `json:"smtp_host"`
	From         string `json:"from"`
	Password     string `json:"password"`
	To           string `json:"to"`
}

const tableStyle = `
<sTyle type='text/css'>
      table {border-collapse: collapse;border-spacing: 0;}
     .tb_1{border:1px solid #cccccc;table-layout:fixed;word-break:break-all;width: 100%;background:#ffffff;margin-bottom:5px}
     .tb_1 caption{text-align: center;background: #F0F4F6;font-weight: bold;padding-top: 5px;height: 25px;border:1px solid #cccccc;border-bottom:none}
     .tb_1 a{margin:0 3px 0 3px}
     .tb_1 tr th,.tb_1 tr td{padding: 3px;border:1px solid #cccccc;line-height:20px}
     .tb_1 thead tr th{font-weight:bold;text-align: center;background:#e3eaee}
     .tb_1 tbody tr th{text-align: right;background:#f0f4f6;padding-right:5px}
     .tb_1 tfoot{color:#cccccc}
     .td_c td{text-align: center}
     .td_r td{text-align: right}
     .t_c{text-align: center !important;}
     .t_r{text-align: right !important;}
     .t_l{text-align: left !important;}
</stYle>
`

// SendMail send mail
func (m *EmailStruct) SendMail(title string, body string) {
	if !m.SendMailAble {
		log.Println("disbale send email")
		return
	}
	if m.SMTPHost == "" || m.From == "" || m.To == "" {
		log.Println("smtp_host,from,to is empty")
		return
	}
	addrInfo := strings.Split(m.SMTPHost, ":")
	if len(addrInfo) != 2 {
		log.Println("smtp_host wrong,eg: host_name:25")
		return
	}
	auth := smtp.PlainAuth("", m.From, m.Password, addrInfo[0])

	_sendTo := strings.Split(m.To, ";")
	var sendTo []string
	for _, _to := range _sendTo {
		_to = strings.TrimSpace(_to)
		if _to != "" && strings.Contains(_to, "@") {
			sendTo = append(sendTo, _to)
		}
	}

	if len(sendTo) < 1 {
		log.Println("mail receiver is empty")
		return
	}

	body = tableStyle + "\n" + body
	body += "<br/><hr style='border:none;border-top:1px solid #ccc'/><center>Powered by <a href='" + AppURL + "'>mysql-schema-sync</a>&nbsp;" + Version + "</center>"

	msgBody := fmt.Sprintf("To: %s\r\nContent-Type: text/html;charset=utf-8\r\nSubject: %s\r\n\r\n%s", strings.Join(sendTo, ";"), title, body)
	err := smtp.SendMail(m.SMTPHost, auth, m.From, sendTo, []byte(msgBody))
	if err == nil {
		log.Println("send mail success")
	} else {
		log.Println("send mail failed,err:", err)
	}
}
