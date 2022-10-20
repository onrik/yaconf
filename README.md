# yaconf
Golang yaml config reader

```golang
package main

import (
	"log"

	"github.com/onrik/yaconf"
)

type Config struct {
	LogFile string `yaml:"log_file" yaconf:"required"`
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
