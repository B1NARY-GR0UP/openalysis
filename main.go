package main

import (
	"fmt"
	"github.com/B1NARY-GR0UP/openalysis/config"
)

func main() {
	fmt.Println(config.DefaultConfig)
	//cmd.Execute()
	if err := config.DefaultConfig.ReadInConfig("./default.yaml"); err != nil {
		panic(err.Error())
	}
	fmt.Println(config.DefaultConfig)
}
