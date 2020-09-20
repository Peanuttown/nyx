package main

import (
	"log"
	"os"

	"github.com/Peanuttown/tzzGoUtil/net"
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
		icmp,err := net.DecodeICMPV4(b[:n])
		if err != nil{
			log.Println(err)
			continue
		}
		log.Println(icmp.Contents)
	}

}
