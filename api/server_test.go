package api

import (
	"man/db"
	"man/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		Environment:         util.EnvironmentTest,
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: 3 * time.Minute,
		BusinessConfig: util.BusinessConfig{
			MinReservationDuration:          30 * time.Minute,
			MaxReservationDuration:          4 * time.Hour,
			MinReservationAdvanceDuration:   30 * time.Minute,
			MaxReservationAdvanceDuration:   7 * 24 * time.Hour, // 7天
			CancellableReservationDuration:  30 * time.Minute,
			ReservationRemindBeforeDuration: 10 * time.Minute,
			ReservationRemindAfterDuration:  15 * time.Minute,
			ReservationViolationDuration:    30 * time.Minute,
		},
	}
	server, err := NewServer(config, store, nil)
	require.NoError(t, err)

	return server
}
