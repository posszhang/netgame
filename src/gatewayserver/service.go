package main

import (
	"base/gnet"
	"base/log"
	"command"
	"net"
)

var timetick *TimeTick
var superclient *SuperClient
var routeManager *RouteManager
var gatewayTaskManager *GatewayTaskManager

type Service struct {
	gnet.NetService
}

func NewService() *Service {
	server := &Service{}
	server.Derived = server

	return server
}

func (server *Service) Init() bool {

	//初始化serverid
	server.SetServerID(command.GetServerID(command.GatewayServer, config.GetInt("server_index")))

	ret := server.Bind("gatewayserver", "", config.GetInt("port"))
	if !ret {
		log.Println("bid port error,service run is error")
		return false
	}

	superclient = NewSuperClient()
	if superclient == nil {
		log.Println("connect superserver is error")
		return false
	}

	routeManager = NewRouteManager()
	gatewayTaskManager = NewGatewayTaskManager()
	timetick = NewTimeTick()

	return true
}

func (server *Service) Final() {

}

func (server *Service) NewTCPTask(conn net.Conn, port int) {

	task := NewGatewayTask()
	task.GoHandler(conn)
}
