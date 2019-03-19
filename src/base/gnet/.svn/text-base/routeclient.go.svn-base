package gnet

import (
	"command"
)

type RouteClient struct {
	TCPClient

	owner *RouteManager
}

func NewRouteClient(mgr *RouteManager) *RouteClient {
	client := &RouteClient{
		owner: mgr,
	}

	client.Derived = client

	return client
}

func (client *RouteClient) MsgParse(msg *command.Message) bool {

	return client.owner.MsgParse(msg)
}

func (client *RouteClient) OnConnected() {

	if client.owner.Derived == nil {
		return
	}

	msg := new(command.ReqServerVerify)
	msg.Info = client.owner.Derived.GetServerInfo()
	client.SendCmd(msg)
}
