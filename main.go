package main

import (
	"fmt"
	"student_go/internal/app"
	"student_go/internal/config"
	"student_go/pkg/log"
)

func main() {
	log.Init()
	r, err := app.SetupApp()
	if err != nil {
		panic(err)
	}
	port := config.Config.Server.Port
	addr := fmt.Sprintf(":%d", port)

	r.Run(addr)
}
