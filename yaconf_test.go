package yaconf

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type testConfig3 struct {
}

func (c *testConfig3) Validate() error {
	return fmt.Errorf("invalid")
}

func TestRead(t *testing.T) {
	file, err := os.CreateTemp("", "")
	require.Nil(t, err)

	defer func() {
		os.Remove(file.Name())
	}()

	_, err = file.Write([]byte(`log_level: info`))
	require.Nil(t, err)

	err = file.Close()
	require.Nil(t, err)

	// Test success
	config := struct {
		LogLevel string `yaml:"log_level"`
	}{}

	err = Read(file.Name(), &config)
	require.Nil(t, err)

	// Test with errors
	config2 := struct {
		LogFile string `yaml:"log_file" yaconf:"required"`
	}{}

	err = Read(file.Name(), &config2)
	require.NotNil(t, err)

	// Test invalid unmarshal
	config3 := struct {
		LogLevel int `yaml:"log_level"`
	}{}
	err = Read(file.Name(), &config3)
	require.NotNil(t, err)

	// Test invalid default valud
	config4 := struct {
		Pid int `yaml:"pid" yaconf:"default=yes"`
	}{}
	err = Read(file.Name(), &config4)
	require.NotNil(t, err)
	require.Equal(t, "yes is invalid value for int", err.Error())

	// Test with custom validator
	err = Read(file.Name(), &testConfig3{})
	require.NotNil(t, err)
	require.Equal(t, "invalid", err.Error())

	// Test invalid file
	err = Read("/tmp/yaconf/test.yml", &testConfig3{})
	require.NotNil(t, err)
}
