# yaconf
Golang yaml config reader

```golang
package main

import (
	"log"
	"os"

	"github.com/onrik/yaconf"
)

type Config struct {
	LogFile  string `yaml:"log_file" yaconf:"required"`
	LogLevel string `yaml:"log_level" yaconf:"default=info"`
}

func (c *Config) Validate() error {
	_, err := os.Stat(c.LogFile)
	return err
}

func main() {
	config := Config{}
	err := yaconf.Read("config.yml", &config)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("%+v\n", config)
}

```
