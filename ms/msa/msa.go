package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	libmsa "github.com/kamasamikon/miego/ms/libmsa"
)

func main() {
	go libmsa.RunService()
	go libmsa.RegisterLoop()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("MSA exiting ...")
}
