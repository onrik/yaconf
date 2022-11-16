package yaconf

import (
	"reflect"
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

	// Test empty
	config2 := &struct {
		Count int `yaconf:"default"`
	}{}
	err = fillDefaultValues(&config2)
	require.Nil(t, err)

	// Test invalid int
	config3 := &struct {
		Count int `yaconf:"default=yes"`
	}{}
	err = fillDefaultValues(&config3)
	require.NotNil(t, err)
	require.Equal(t, "yes is invalid value for int", err.Error())

	// Test invalid uint
	config4 := &struct {
		Count uint `yaconf:"default=a"`
	}{}
	err = fillDefaultValues(&config4)
	require.NotNil(t, err)
	require.Equal(t, "a is invalid value for uint", err.Error())

	// Test invalid duration
	config5 := &struct {
		Timeout time.Duration `yaconf:"default=22"`
	}{}
	err = fillDefaultValues(&config5)
	require.NotNil(t, err)
	require.Equal(t, "22 is invalid value for time.Duration", err.Error())

	// Test invalid bool
	config6 := &struct {
		Debug bool `yaconf:"default=22"`
	}{}
	err = fillDefaultValues(&config6)
	require.NotNil(t, err)
	require.Equal(t, "22 is invalid value for bool", err.Error())

}

func TestSetDefaultValue(t *testing.T) {
	type db struct {
		Host string `yaconf:"default=localhost"`
	}
	config := &struct {
		LogLevel string `yaconf:"default=info"`
		LogFile  string
		private  string `yaconf:"default=test"`
		DB       *db
	}{
		DB: &db{},
	}

	err := setDefaultValue(reflect.ValueOf(config))
	require.Nil(t, err)
	require.Empty(t, config.private)
	require.Equal(t, "localhost", config.DB.Host)

	// Test empty

	// Test invalid bool
	config2 := struct {
		Debug bool `yaconf:"default=22"`
	}{}
	err = setDefaultValue(reflect.ValueOf(config2))
	require.NotNil(t, err)
}

func TestIsDuration(t *testing.T) {
	require.True(t, isDuration(reflect.ValueOf(time.Duration(1000))))
	require.False(t, isDuration(reflect.ValueOf(int64(22))))
}

func TestIsInt(t *testing.T) {
	require.True(t, isInt(reflect.ValueOf(int8(1))))
	require.True(t, isInt(reflect.ValueOf(int16(2))))
	require.True(t, isInt(reflect.ValueOf(int32(3))))
	require.True(t, isInt(reflect.ValueOf(int64(4))))
	require.True(t, isInt(reflect.ValueOf(int8(5))))

	require.False(t, isInt(reflect.ValueOf(byte('6'))))
	require.False(t, isInt(reflect.ValueOf("22")))
}

func TestIsUint(t *testing.T) {
	require.True(t, isUint(reflect.ValueOf(uint8(1))))
	require.True(t, isUint(reflect.ValueOf(uint16(2))))
	require.True(t, isUint(reflect.ValueOf(uint32(3))))
	require.True(t, isUint(reflect.ValueOf(uint64(4))))
	require.True(t, isUint(reflect.ValueOf(uint8(5))))

	require.False(t, isUint(reflect.ValueOf("22")))
}
