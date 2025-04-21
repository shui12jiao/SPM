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
	}
	server, err := NewServer(config, store, nil)
	require.NoError(t, err)

	return server
}
