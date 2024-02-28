package config

import (
	"fmt"
	"testing"
)

var path = "../default.yaml"

func TestReadInConfig(t *testing.T) {
	fmt.Println(DefaultConfig)
	err := DefaultConfig.ReadInConfig(path)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(DefaultConfig)
}
