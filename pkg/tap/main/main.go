package main

import "flag"


func main() {
	var cfgFile string
	flag.StringVar(&cfgFile,"config","/etc/nyx.yaml","config file")
	flag.Parse()
	var devName = "nyxTap"
}
