package main

import (
	"flag"
	"os"

	"github.com/Peanuttown/nyx/client"
	"github.com/Peanuttown/nyx/pkg/cfg"
	"github.com/Peanuttown/tzzGoUtil/encoding"
	"github.com/Peanuttown/tzzGoUtil/log"
	"gopkg.in/yaml.v2"
)

func main() {
	logger := log.NewLogger()
	var cfgFile string
	flag.StringVar(&cfgFile, "config", "", "config file path")
	flag.Parse()
	config := &cfg.ClientCfg{}
	if len(cfgFile) != 0 {
		err := encoding.UnMarshalByFile(cfgFile, config, yaml.Unmarshal)
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}
	}
	cfg.EndowDefault(config)
	err := config.Check()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	} 
	err = client.NewClient().Run(config)
	if err != nil{
		logger.Error(err)
	}
 }
