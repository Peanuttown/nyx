package main

import (
	"log"
	"os"

	"github.com/songgao/water"
)

func main() {
	ifce, err := water.NewTAP("waterDev")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	var b = make([]byte, 1024)
	for {
		n, err := ifce.Read(b)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		log.Println(string(b[:n]))
	}

}
