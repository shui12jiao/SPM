package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendEmail(t *testing.T) {
	// 当Short()为true时，跳过测试
	if testing.Short() {
		t.Skip("short模式下跳过测试")
	}

	config := LoadConfig()

	sender := NewEmailSender(config.EmailConfig)
	subject := "测试邮件"
	content := `
	<h1>自习室管理系统</h1>
	<p>这是一封测试邮件</p>
	<p>请勿回复</p>
	`
	to := []string{"1873978303@qq.com"}
	attachFiles := []string{}
	err := sender.SendEmail(to, subject, content, nil, nil, attachFiles)
	require.NoError(t, err)
}
