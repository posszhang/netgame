package main

import (
	"base/log"
	"base/util"
	"command"
	"flag"
)

var config *util.Config
var service *Service

func main() {

	var logfile string
	var port int

	flag.StringVar(&logfile, "logfile", "", "日志文件路径")
	flag.IntVar(&port, "port", 0, "监听端口,0表示自动分配")
	flag.Parse()

	log.NewLog(logfile)
	log.Println("服务器启动")

	var err error = nil
	config, err = util.NewConfig("config.json", "routeserver")
	if err != nil {
		log.Errorln("读取config.json失败")
		return
	}
	//记录服务器序号
	config.Set("server_index", command.GetIndexFromFilename(logfile))
	config.Set("port", port)

	service = NewService()
	service.Run()

	log.Println("服务器关闭")
}
