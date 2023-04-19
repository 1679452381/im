package test

import (
	"crypto/tls"
	"github.com/jordan-wright/email"
	"net/smtp"
	"testing"
)

func TestEmailCode(t *testing.T) {
	e := email.NewEmail()
	e.From = "GET <17700611471@163.com>"
	e.To = []string{"15660589213@163.com"}
	e.Subject = "验证码已发送"
	e.HTML = []byte("您的验证码：<b>123123</b>")
	//err := e.Send("smtp.163.com:465", smtp.PlainAuth("", "15660589213@163.com", "DSBZHQSKFWQVDSVK", "smtp.163.com"))
	//返回EOF 关闭SSL重试
	err := e.SendWithTLS("smtp.163.com:465", smtp.PlainAuth("", "17700611471@163.com", "DSBZHQSKFWQVDSVK", "smtp.163.com"), &tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
	if err != nil {
		t.Fatal(err.Error())
	}
}
