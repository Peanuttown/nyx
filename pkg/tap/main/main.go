package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/Peanuttown/nyx/pkg/tap"
)

func main() {
	var devName = "testDev"
	tapDev, err := tap.NewTap(devName)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	err = tap.SetIP(devName, "172.16.2.151/16")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	var b = make([]byte, 1024)
	for {
		n, err := tapDev.Read(b)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(b[:n]))
	}
	notify := make(chan os.Signal)
	signal.Notify(notify, os.Interrupt)
	<-notify
}
