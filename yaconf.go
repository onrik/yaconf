package yaconf

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Read config from file
func Read(filename string, config interface{}) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = FillDefaultValues(config)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return err
	}

	err = Validate(config)
	if err != nil {
		return err
	}

	if v, ok := config.(validator); ok {
		return v.Validate()
	}

	return nil
}

func Validate(config interface{}) error {
	errors := validate(config, "")
	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, ", "))
	}

	return nil
}
