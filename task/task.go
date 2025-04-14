package task

import (
	"context"
	"database/sql"
	"man/util"

	"github.com/rs/zerolog/log"
)

// 更新自习室签到码
func (s *cronScheduler) updateRoomSignCode() error {
	// 获取所有自习室数量
	num, err := s.store.CountRoom(context.Background(), sql.NullBool{Bool: true, Valid: true})
	if err != nil {
		log.Error().Err(err).Msg("获取自习室数量失败")
		return err
	}
	// 更新签到码
	var codes []string
	for i := 0; i < int(num); i++ {
		codes = append(codes, util.GenerateSignCode())
	}

	// 更新数据库
	err = s.store.UpdateAllRoomCode(context.Background(), codes)
	if err != nil {
		log.Error().Err(err).Msg("更新签到码失败")
		return err
	}
	log.Info().Strs("codes", codes).Int64("num", num).Msg("更新签到码成功")
	return nil
}
