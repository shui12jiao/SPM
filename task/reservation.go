package task

import (
	"context"
	"fmt"
	"man/db"
	"man/util"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type CreateReservationTaskArg struct {
	//  预约ID
	ReservationID uuid.UUID
	// 预约用户ID
	UserID int32

	//   预约开始时间
	ReservationTime time.Time
	//   预约开始前提醒时间
	RemindBeforeDuration time.Duration
	//   预约开始后提醒时间
	RemindAfterDuration time.Duration
	//   违约处理时间
	ViolationDuration time.Duration
}

func (s *cronScheduler) CreateReservationTasks(arg CreateReservationTaskArg) error {
	// 链式处理预约创建任务
	// 任务列表：
	// 在预约开始前15min发送邮件提醒
	// 在预约开始后10min，检查是否签到，未签到则发送邮件提醒
	// 在预约开始后15min，检查是否签到，未签到则取消预约并进行违约处理

	// 获取用户信息
	user, err := s.store.GetUser(context.Background(), arg.UserID)
	if err != nil {
		return err
	}

	// 1. 第一次提醒 在预约开始前
	s.AddOnceJob(
		"预约开始前提醒",
		arg.ReservationID.String(),
		arg.ReservationTime.Add(-arg.RemindBeforeDuration),
		func() error {
			// 发送邮件
			return s.email.SendEmail(
				[]string{user.Email},
				"自习室签到提醒",
				fmt.Sprintf("你的自习室预约将在%s开始，请及时签到。", arg.ReservationTime.Format(time.TimeOnly)),
				[]string{},
				[]string{},
				[]string{},
			)
		},
	)

	// 2. 第二次提醒 在预约开始后10min，检查是否签到，未签到则发送邮件提醒
	s.AddOnceJob(
		"预约开始后提醒",
		arg.ReservationID.String(),
		arg.ReservationTime.Add(-arg.RemindBeforeDuration),
		func() error {
			// 获取预约信息
			// 闭包参数reservationID
			reservation, err := s.store.GetReservation(context.Background(), arg.ReservationID)
			if err != nil {
				return err
			}

			// 检查是否签到
			if reservation.Status != util.ReservationStatusReserved {
				// 取消违约处理
				// 注意这不会影响到该task的执行
				s.RemoveByTags(arg.ReservationID.String())

				// 取消和已签到的预约不需要发送提醒
				return nil
			}

			// 发送邮件
			return s.email.SendEmail(
				[]string{user.Email},
				"自习室签到提醒",
				fmt.Sprintf("你的自习室预约已于%s开始，请及时签到！超时%.0f分钟未签到将视为违约。",
					arg.ReservationTime.Format(time.TimeOnly),
					arg.ViolationDuration.Minutes(),
				),
				[]string{},
				[]string{},
				[]string{},
			)
		},
	)

	// 3. 第三次提醒 在预约开始后15min，检查是否签到，未签到则取消预约并进行违约处理
	s.AddOnceJob(
		"违约处理",
		arg.ReservationID.String(),
		arg.ReservationTime.Add(arg.ViolationDuration),
		func() error {
			// 获取预约信息
			// 闭包参数reservationID
			reservation, err := s.store.GetReservation(context.Background(), arg.ReservationID)
			if err != nil {
				return err
			}

			// 检查是否签到
			if reservation.Status != util.ReservationStatusReserved {
				return nil
			}

			// 设置为违约状态，创建违约记录
			ctx := context.Background()
			_, err = s.store.UpdateReservationStatus(ctx, db.UpdateReservationStatusParams{
				ID:     arg.ReservationID,
				Status: util.ReservationStatusViolated,
			})
			if err != nil {
				log.Error().Err(err).Str("reservationID", arg.ReservationID.String()).Msg("更新预约状态失败")
				return err
			}

			_, err = s.store.CreateViolation(ctx, db.CreateViolationParams{
				ReservationID: arg.ReservationID,
				UserID:        arg.UserID,
			})
			if err != nil {
				log.Error().Err(err).Str("reservationID", arg.ReservationID.String()).Msg("创建违约记录失败")
				return err
			}

			return nil
		},
	)

	log.Info().Str("reservationID", arg.ReservationID.String()).Msg("创建预约任务成功")
	return nil
}
