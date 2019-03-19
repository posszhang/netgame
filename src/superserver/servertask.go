package main

import (
	"base/gnet"
	"base/log"
	"command"
	"github.com/golang/protobuf/proto"
)

type ServerTask struct {
	gnet.TCPTask
	serverinfo command.ServerInfo
}

func NewServerTask() *ServerTask {
	task := &ServerTask{}
	task.Derived = task

	return task
}

func (task *ServerTask) VerifyConn(msg *command.Message) bool {

	cmd := new(command.ReqServerVerify)
	if err := proto.Unmarshal(msg.Data, cmd); err != nil {
		return false
	}

	if cmd == nil || cmd.Info == nil {
		return false
	}

	task.serverinfo = *cmd.Info

	log.Println("verify conn", task.GetServerInfo())

	if !serverManager.UniqueAdd(task) {
		return false
	}

	//服务器新增时，回调处理特殊逻辑
	task.onServerAddCallback()
	return true
}

func (task *ServerTask) onServerAddCallback() {

	srvtp := task.serverinfo.Type

	//不是路由服务器，则广播路由服务器
	if srvtp != command.RouteServer {
		serverManager.NotifyRouteServerInit(task)
	} else {
		//是路由服务器代表主动新增
		serverManager.NotifyRouteServerAdd(task)
	}

	//登陆服务器新增，主动刷新网关
	if srvtp == command.GatewayServer {
		serverManager.NotifyGate2Login(task)
	}

}

func (task *ServerTask) RecycleConn() bool {

	serverManager.UniqueRemove(task)

	task.onServerRemoveCallback()

	return true
}

func (task *ServerTask) onServerRemoveCallback() {

	//srvtp := task.serverinfo.Type
}

func (task *ServerTask) GetID() uint32 {
	return task.serverinfo.Id
}

func (task *ServerTask) GetServerInfo() command.ServerInfo {
	return task.serverinfo
}

func (task *ServerTask) GetServerType() uint32 {
	return task.serverinfo.Type
}

func (task *ServerTask) MsgParse(msg *command.Message) bool {
	return true
}
