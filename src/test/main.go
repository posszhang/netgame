package main

import (
	"base/log"
	"flag"
	"time"
)

func main() {

	log.NewLog("")

	ip := ""
	port := 0
	num := 0

	flag.StringVar(&ip, "ip", "127.0.0.1", "目标登陆服务器IP")
	flag.IntVar(&port, "port", 8010, "目标登陆服务器端口")
	flag.IntVar(&num, "num", 1, "机器人数量")
	flag.Parse()

	robotManager := NewRobotManager()
	robotManager.Add(num, ip, port)

	for {
		robotManager.Do()
		time.Sleep(1000)
	}
}
