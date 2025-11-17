package yaconf

import (
	"errors"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Read config from file
func Read(filename string, config any) error {
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

func Validate(config any) error {
	errs := validate(config, "")
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}

	return nil
}
