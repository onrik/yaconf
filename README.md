# yaconf

[![Tests](https://github.com/onrik/yaconf/workflows/Tests/badge.svg)](https://github.com/onrik/yaconf/actions)
[![Coverage Status](https://coveralls.io/repos/github/onrik/yaconf/badge.svg?branch=main)](https://coveralls.io/github/onrik/yaconf?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/onrik/yaconf)](https://goreportcard.com/report/github.com/onrik/yaconf)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/onrik/yaconf)](https://pkg.go.dev/github.com/onrik/yaconf)

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
