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

	return true
}

func (task *ServerTask) RecycleConn() bool {

	serverManager.UniqueRemove(task)

	return true
}

func (task *ServerTask) MsgParse(msg *command.Message) bool {
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
