package main

import (
	"base/gnet"
	"base/log"
	"command"
	"github.com/golang/protobuf/proto"
)

type ServerTask struct {
	gnet.TCPTask

	msgHandler gnet.MessageHandler

	serverinfo command.ServerInfo
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

	return true
}

func (task *ServerTask) RecycleConn() bool {

	log.Println("删除服务器", task.GetServerInfo())

	serverManager.UniqueRemove(task)

	return true
}

func (task *ServerTask) init() {

	task.msgHandler.Reg(&command.RouteBroadcastByType{}, task.onRouteBroadcastByType)
	task.msgHandler.Reg(&command.RouteBroadcastByID{}, task.onRouteBroadcastByID)
}

func (task *ServerTask) MsgParse(msg *command.Message) bool {

	//task.msgHandler.Reg(&command.RouteBroadcastById, task.onRouteBoradcastByType)
	task.msgHandler.Process(msg)

	return true
}

func (task *ServerTask) GetServerInfo() *command.ServerInfo {
	return &task.serverinfo
}

func (task *ServerTask) GetID() uint32 {
	return task.serverinfo.Id
}

func (task *ServerTask) GetServerType() uint32 {
	return task.serverinfo.Type
}

func (task *ServerTask) onRouteBroadcastByType(cmd proto.Message) {

	msg := cmd.(*command.RouteBroadcastByType)

	log.Println(msg.Msg.Name)

	serverlist := serverManager.GetByType(msg.Type)
	if len(serverlist) == 0 {
		log.Println("boradcast by type error, servers is null", msg.Type)
		serverManager.Print()
		return
	}

	for _, server := range serverlist {
		server.SendCmd(cmd)
	}
}

func (task *ServerTask) onRouteBroadcastByID(cmd proto.Message) {

	msg := cmd.(*command.RouteBroadcastByID)

	log.Println(msg.Msg.Name)

	server := serverManager.GetByID(msg.Id)
	if server == nil {
		log.Println("boradcast by id error, servers is null", msg.Id)
		serverManager.Print()
		return
	}

	server.SendCmd(cmd)
}
