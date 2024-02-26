package config

import (
	"fmt"
	"testing"
)

var path = "../default.yaml"

var defaultConfig = &Config{}

func TestReadInConfig(t *testing.T) {
	fmt.Println(defaultConfig)
	err := defaultConfig.ReadInConfig(path)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(defaultConfig)
}
