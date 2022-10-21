package yaconf

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFillDefaultValues(t *testing.T) {
	config := &struct {
		PID      uint   `yaconf:"default=90"`
		LogLevel string `yaconf:"default=debug"`
		DB       struct {
			Host    string        `yaconf:"default=localhost"`
			Port    int           `yaconf:"default=5432"`
			Timeout time.Duration `yaconf:"default=2s"`
			Debug   bool          `yaconf:"default=true"`
		}
	}{}

	err := fillDefaultValues(config)
	require.Nil(t, err)
	require.Equal(t, uint(90), config.PID)
	require.Equal(t, "debug", config.LogLevel)
	require.Equal(t, "localhost", config.DB.Host)
	require.Equal(t, 5432, config.DB.Port)
	require.Equal(t, 2*time.Second, config.DB.Timeout)
	require.Equal(t, true, config.DB.Debug)
}
