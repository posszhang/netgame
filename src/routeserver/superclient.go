package main

import (
	"base/gnet"
	"command"
)

type SuperClient struct {
	gnet.TCPClient
}

func NewSuperClient() *SuperClient {
	client := &SuperClient{}
	client.Derived = client

	if !client.Connect(config.GetString("super_ip"), config.GetInt("super_port")) {
		return nil
	}

	return client
}

func (client *SuperClient) MsgParse(msg *command.Message) bool {
	return true
}

func (client *SuperClient) OnConnected() {

	msg := new(command.ReqServerVerify)
	msg.Info = service.GetServerInfo()
	client.SendCmd(msg)
}
