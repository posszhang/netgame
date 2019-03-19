package main

import (
	"base/log"
	"base/util"
	"flag"
)

var config *util.Config

func main() {

	var logfile string

	flag.StringVar(&logfile, "logfile", "", "日志文件路径")
	flag.Parse()

	log.NewLog(logfile)

	var err error = nil
	config, err = util.NewConfig("config.json", "superserver")
	if err != nil {
		log.Errorln("读取config.json失败")
		return
	}

	server := NewService()
	server.Run()

}
