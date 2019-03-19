package main

import (
	"base/gnet"
	"log"
	"net"
)

type MyServer struct {
	gnet.NetService
}

func NewMyServer() *MyServer {
	server := &MyServer{}
	server.Derived = server

	return server
}

func (server *MyServer) Init() bool {

	ret := server.Bind("test", "", 0)
	if !ret {
		log.Println("bid port error,service run is error")
		return false
	}

	return true
}

func (server *MyServer) Final() {

}

func (server *MyServer) NewTCPTask(conn net.Conn, port int) {

	log.Println("new tcp task conn", port)

	task := NewMyTask()
	task.GoHandler(conn)
}
