package main

import (
	"fmt"
	"log"

	"github.com/xuyun-io/scalemetric/route"
)

func main() {
	route := route.InitRoute()
	if err := route.Run(); err != nil {
		log.Panic(fmt.Errorf("run failed, %v, exit", err.Error()))
	}
}
