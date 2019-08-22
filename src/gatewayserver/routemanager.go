package main

import (
	"base/gnet"
	"base/log"
	"command"
	"github.com/golang/protobuf/proto"
)

type RouteManager struct {
	gnet.RouteManager

	msgHandler gnet.MessageHandler
}

func NewRouteManager() *RouteManager {
	mgr := &RouteManager{}
	mgr.Derived = mgr

	mgr.init()

	return mgr
}

func (mgr *RouteManager) GetServerInfo() *command.ServerInfo {
	return service.GetServerInfo()
}

func (mgr *RouteManager) MsgParse(msg *command.Message) bool {

	log.Println("route manager:", msg)

	mgr.msgHandler.Process(msg)

	return true
}

func (mgr *RouteManager) init() {

	mgr.msgHandler.Reg(&command.RetGateRegUser{}, mgr.onRetGateRegUser)
}

func (mgr *RouteManager) onRetGateRegUser(cmd proto.Message) {

	msg := cmd.(*command.RetGateRegUser)

	log.Println("数据库返回用户数据:", msg)
}
