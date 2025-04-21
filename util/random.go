package util

import (
	"math/rand"
	"strings"
	"time"
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	numberRunes = []rune("0123456789")
	r           = rand.New(rand.NewSource(time.Now().UnixNano())) // 使用局部 rand 实例
)

// 生成6位数字签到码
func GenerateSignCode() string {
	return RandomNumber(6)
}

// 生成随机字符串，仅使用字母
func RandomString(n int) string {
	return RandomRunes(n, letterRunes)
}

// 生成随机字符串，仅使用数字
func RandomNumber(n int) string {
	return RandomRunes(n, numberRunes)
}

func RandomEmail() string {
	return RandomRunes(6, letterRunes) + "@example.com"
}

// 生成指定字符集的随机字符串
func RandomRunes(n int, charset []rune) string {
	var sb strings.Builder
	k := len(charset)

	for i := 0; i < n; i++ {
		sb.WriteRune(charset[r.Intn(k)])
	}

	return sb.String()
}
