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

	err = fillDefaultValues(config)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return err
	}

	errors := validate(config, "")
	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, ", "))
	}

	if v, ok := config.(validator); ok {
		return v.Validate()
	}

	return nil
}
