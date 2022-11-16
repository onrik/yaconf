package yaconf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddPrefix(t *testing.T) {
	s := addPrefix("", "bar")
	require.Equal(t, "bar", s)

	s = addPrefix("foo", "bar")
	require.Equal(t, "foo.bar", s)
}

func TestValidate(t *testing.T) {
	type db struct {
		Host string `yaconf:"required"`
	}
	pid := "/var/run/test.pid"
	config := struct {
		LogLevel string  `yaml:"log_level" yaconf:"required,default=info"`
		LogFile  string  `yaml:"log_file" yaconf:"required"`
		PID      *string `yaconf:"required"`
		DB1      db
		DB2      *db
		DB3      *db
		private  string
	}{
		PID: &pid,
		DB3: &db{
			Host: "127.0.0.1",
		},
	}

	errors := validate(&config, "")
	require.Equal(t, 4, len(errors), errors)
	require.Equal(t, "log_level is required", errors[0])
	require.Equal(t, "log_file is required", errors[1])
	require.Equal(t, "DB1.Host is required", errors[2])
	require.Equal(t, "DB2.Host is required", errors[3])
}
