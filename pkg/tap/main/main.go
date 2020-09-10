package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/Peanuttown/nyx/pkg/tap"
)

func main() {
	_, err := tap.NewTap("testDev")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	notify := make(chan os.Signal)
	signal.Notify(notify, os.Interrupt)
	<-notify
}
