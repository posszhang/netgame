package main

import (
	"base/gnet"
	"base/log"
	"net"
)

var serverManager *ServerManager

type Service struct {
	gnet.NetService
}

func NewService() *Service {
	server := &Service{}
	server.Derived = server

	return server
}

func (server *Service) Init() bool {

	ret := server.Bind("test", "", config.GetInt("port"))
	if !ret {
		log.Println("bid port error,service run is error")
		return false
	}

	serverManager = NewServerManager()

	return true
}

func (server *Service) Final() {

}

func (server *Service) NewTCPTask(conn net.Conn, port int) {

	log.Println("new tcp task conn", port)

	task := NewServerTask()
	task.GoHandler(conn)
}
