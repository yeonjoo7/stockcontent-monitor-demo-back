package config

import (
	"encoding/json"
	"os"
)

const (
	localConfigFilePath = "./config.local.json"
	configFilePath      = "./config.json"
)

var (
	_default configImpl
	Default  Config = &_default
)

func init() {
	var file *os.File
	var err error

	if Default.IsDebug() {
		file, err = os.Open(localConfigFilePath)
	} else {
		file, err = os.Open(configFilePath)
	}

	if err != nil {
		panic(err)
	}

	err = json.NewDecoder(file).Decode(&_default)
	if err != nil {
		panic(err)
	}

}
