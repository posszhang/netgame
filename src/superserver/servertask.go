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

	msgHandler gnet.MessageHandler
}

func NewServerTask() *ServerTask {
	task := &ServerTask{}
	task.Derived = task

	task.init()

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

	if !serverManager.UniqueAdd(task) {
		return false
	}

	log.Println("新增服务器", task.GetServerInfo())

	//服务器新增时，回调处理特殊逻辑
	task.onServerAddCallback()

	snd := new(command.RetServerVerify)
	task.SendCmd(snd)

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
}

func (task *ServerTask) RecycleConn() bool {

	log.Println("删除服务器", task.GetServerInfo())

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

func (task *ServerTask) GetServerInfo() *command.ServerInfo {
	return &task.serverinfo
}

func (task *ServerTask) GetServerType() uint32 {
	return task.serverinfo.Type
}

func (task *ServerTask) init() {
	task.msgHandler.Reg(&command.ReqGatewayList{}, task.onReqGatewayList)
}

func (task *ServerTask) MsgParse(msg *command.Message) bool {

	log.Println(msg)

	task.msgHandler.Process(msg)

	return true
}

func (task *ServerTask) onReqGatewayList(cmd proto.Message) {

	snd := new(command.RetGatewayList)

	serverlist := serverManager.GetByType(command.GatewayServer)
	snd.Serverlist = make([]*command.ServerInfo, 0, len(serverlist))

	for _, server := range serverlist {
		snd.Serverlist = append(snd.Serverlist, server.GetServerInfo())
	}
	task.SendCmd(snd)
}
