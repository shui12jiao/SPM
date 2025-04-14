package task

import (
	"man/db"
	"man/util"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Scheduler interface {
	Start()
	AddJob(jobName string, tag string, crontab string, function any, parameters ...any)
	AddOnceJob(jobName string, tag string, time time.Time, function any, parameters ...any)
	RemoveByTags(tags ...string)

	CreateReservationTasks(arg CreateReservationTaskArg) error
}

func NewScheduler(config util.Config, store db.Store) Scheduler {
	// 初始化gocron调度器
	scheduler, err := gocron.NewScheduler(
		gocron.WithLocation(time.Local), // 设置时区
		gocron.WithGlobalJobOptions(),   // 全局任务选项
	)
	if err != nil {
		log.Fatal().Err(err).Msg("无法创建gocron调度器")
	}

	// 创建email sender
	email := util.NewEmailSender(
		util.EmailConfig{
			SMTPHost:     "smtp.example.com",
			SMTPPort:     587,
			SenderEmail:  config.EmailSenderAddress,
			SenderName:   config.EmailSenderName,
			SMTPPassword: config.EmailSenderPassword,
		},
	)

	s := &cronScheduler{
		store:     store,
		scheduler: scheduler,
		email:     email,
	}

	// 为scheduler添加服务器定时任务
	s.AddJob("updateRoomSignCode", "room", "0 0 * * *", s.updateRoomSignCode) // 每天0点更新自习室签到码

	return s
}

type cronScheduler struct {
	store     db.Store
	scheduler gocron.Scheduler
	email     util.EmailSender
}

func (s *cronScheduler) Start() {
	// Start是非阻塞的，调用后不会等待任务完成
	s.scheduler.Start()
	log.Info().Msg("定时任务调度器已启动")
}

// 使用cron表达式添加定时任务
func (s *cronScheduler) AddJob(jobName string, tag string, crontab string, function any, parameters ...any) {
	_, err := s.scheduler.NewJob(
		gocron.CronJob(crontab, true),                                                // Cron 表达式
		gocron.NewTask(function, parameters...),                                      // 任务函数和参数
		gocron.WithEventListeners(panicListener(), errorListener(), afterListener()), // 事件监听
		gocron.WithName(jobName),                                                     // 任务名称
		gocron.WithTags(tag),                                                         // 任务标签
		gocron.WithSingletonMode(gocron.LimitModeReschedule),                         // 单例模式
	)
	if err != nil {
		log.Error().Err(err).Str("jobName", jobName).Msg("无法创建定时任务")
		return
	}
	log.Info().Str("jobName", jobName).Str("crontab", crontab).Msg("创建定时任务成功")
}

// 一次性任务，于time执行
func (s *cronScheduler) AddOnceJob(jobName string, tag string, time time.Time, function any, parameters ...any) {
	_, err := s.scheduler.NewJob(
		gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(time)),                      // 一次性任务
		gocron.NewTask(function, parameters...),                                      // 任务函数和参数
		gocron.WithEventListeners(panicListener(), errorListener(), afterListener()), // 事件监听
		gocron.WithName(jobName),                                                     // 任务名称
		gocron.WithTags(tag),                                                         // 任务标签
		gocron.WithSingletonMode(gocron.LimitModeReschedule),                         // 单例模式
	)
	if err != nil {
		log.Error().Err(err).Str("jobName", jobName).Msg("无法创建定时任务")
		return
	}
	log.Info().Str("jobName", jobName).Msg("创建一次性定时任务成功")
}

// 删除指定标签的任务
func (s *cronScheduler) RemoveByTags(tags ...string) {
	// 删除指定标签的任务
	if len(tags) == 0 {
		log.Warn().Msg("没有指定标签，无法删除定时任务")
		return
	}
	s.scheduler.RemoveByTags(tags...)
	log.Info().Strs("tags", tags).Msg("删除定时任务成功")
}

func afterListener() gocron.EventListener {
	return gocron.AfterJobRuns(
		func(jobID uuid.UUID, jobName string) {
			log.Info().Str("jobName", jobName).Msg("定时任务执行完成")
		},
	)
}

func panicListener() gocron.EventListener {
	return gocron.AfterJobRunsWithPanic(
		func(jobID uuid.UUID, jobName string, recoverData any) {
			log.Error().Str("jobName", jobName).Str("recoverData", recoverData.(string)).Msg("定时任务执行异常")
		},
	)
}

func errorListener() gocron.EventListener {
	return gocron.AfterJobRunsWithError(
		func(jobID uuid.UUID, jobName string, err error) {
			log.Error().Str("jobName", jobName).Err(err).Msg("定时任务执行失败")
		},
	)
}
