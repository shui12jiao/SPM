package util

import (
	"errors"
	"strings"

	"gopkg.in/gomail.v2"
)

type EmailConfig struct {
	SMTPHost     string // SMTP服务器地址（如smtp.163.com）
	SMTPPort     int    // 端口号（如25/465/587）
	SenderEmail  string // 发件邮箱
	SenderName   string // 发件人别名（如"系统通知"）
	SMTPPassword string // SMTP密码/授权码
}

type EmailSender interface {
	SendEmail(to []string, subject string, body string, cc []string, bcc []string, attachments []string) error
}

type GmailSender struct {
	config EmailConfig
}

func NewEmailSender(config EmailConfig) EmailSender {
	return &GmailSender{config: config}
}

func (es *GmailSender) SendEmail(
	to []string,
	subject string,
	body string,
	cc []string,
	bcc []string,
	attachments []string,
) error {
	// 参数校验
	if len(to) == 0 {
		return errors.New("收件人不能为空")
	}

	// 创建邮件主体
	m := gomail.NewMessage(gomail.SetEncoding(gomail.Base64))
	m.SetHeader("From", m.FormatAddress(es.config.SenderEmail, es.config.SenderName))
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)

	// 内容类型自动判断
	contentType := "text/plain"
	if strings.Contains(body, "</") {
		contentType = "text/html"
	}
	m.SetBody(contentType, body)

	// 添加抄送/密送
	if len(cc) > 0 {
		m.SetHeader("Cc", cc...)
	}
	if len(bcc) > 0 {
		m.SetHeader("Bcc", bcc...)
	}

	// 处理附件
	for _, file := range attachments {
		m.Attach(file)
	}

	// 创建SMTP客户端
	d := gomail.NewDialer(
		es.config.SMTPHost,
		es.config.SMTPPort,
		es.config.SenderEmail,
		es.config.SMTPPassword,
	)

	// 发送邮件
	return d.DialAndSend(m)
}
