package task

import (
	"man/db/mockdb"
	"man/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewScheduler(t *testing.T) {
	store := mockdb.NewMockStore(t)
	config := util.Config{
		Environment: "test",
		EmailConfig: util.EmailConfig{
			SenderName:   "Test",
			SenderEmail:  "test@example.com",
			SMTPHost:     "smtp.example.com",
			SMTPPort:     587,
			SMTPPassword: "pass",
		},
	}

	scheduler := NewScheduler(config, store)
	require.NotNil(t, scheduler, "调度器不应为空")
}

func TestAddJob(t *testing.T) {
	store := mockdb.NewMockStore(t)
	config := util.Config{
		Environment: "test",
		EmailConfig: util.EmailConfig{
			SenderName:   "Test",
			SenderEmail:  "test@example.com",
			SMTPHost:     "smtp.example.com",
			SMTPPort:     587,
			SMTPPassword: "pass",
		},
	}

	scheduler := NewScheduler(config, store)

	// 测试添加任务
	testFunc := func() {
		// 仅用于测试，不需要实际执行
	}

	scheduler.AddJob("testJob", "test", "* * * * *", testFunc)

	// 启动调度器
	scheduler.Start()

	// 测试添加任务不会panic
	require.NotPanics(t, func() {
		scheduler.AddJob("anotherJob", "test", "* * * * *", testFunc)
	})
}

func TestAddOnceJob(t *testing.T) {
	store := mockdb.NewMockStore(t)
	config := util.Config{
		Environment: "test",
		EmailConfig: util.EmailConfig{
			SenderName:   "Test",
			SenderEmail:  "test@example.com",
			SMTPHost:     "smtp.example.com",
			SMTPPort:     587,
			SMTPPassword: "pass",
		},
	}

	scheduler := NewScheduler(config, store)

	// 测试添加一次性任务
	testFunc := func() {
		// 仅用于测试，不需要实际执行
	}

	executionTime := time.Now().Add(1 * time.Hour)
	scheduler.AddOnceJob("testOnceJob", "test", executionTime, testFunc)

	// 测试添加任务不会panic
	require.NotPanics(t, func() {
		scheduler.AddOnceJob("anotherOnceJob", "test", executionTime, testFunc)
	})
}

func TestRemoveByTags(t *testing.T) {
	store := mockdb.NewMockStore(t)
	config := util.Config{
		Environment: "test",
		EmailConfig: util.EmailConfig{
			SenderName:   "Test",
			SenderEmail:  "test@example.com",
			SMTPHost:     "smtp.example.com",
			SMTPPort:     587,
			SMTPPassword: "pass",
		},
	}

	scheduler := NewScheduler(config, store)

	// 添加一个测试任务
	testFunc := func() {
		// 仅用于测试，不需要实际执行
	}

	scheduler.AddJob("testTagJob", "testTag", "* * * * *", testFunc)

	// 删除任务
	scheduler.RemoveByTags("testTag")

	// 测试删除任务不会panic
	require.NotPanics(t, func() {
		scheduler.RemoveByTags("nonExistentTag")
	})
}
