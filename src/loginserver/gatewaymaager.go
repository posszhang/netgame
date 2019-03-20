package main

import (
	"base/log"
	"command"
	"sync"
)

// 网关结点信息
type GatewayNode struct {
	Info   *command.ServerInfo
	Online uint32
}

type GatewayManager struct {
	gylist map[uint32]*GatewayNode

	mutex sync.Mutex
}

func NewGatewayManager() *GatewayManager {
	mgr := &GatewayManager{
		gylist: make(map[uint32]*GatewayNode),
	}

	return mgr
}

func (mgr *GatewayManager) ResetGateWayList(serverList []*command.ServerInfo) {

	nodeList := make(map[uint32]*GatewayNode)

	for _, server := range serverList {
		node := new(GatewayNode)
		node.Info = server
		node.Online = 0

		if _, ok := mgr.gylist[server.Id]; ok {
			node.Online = mgr.gylist[server.Id].Online
		}

		nodeList[server.Id] = node
	}

	mgr.gylist = nodeList

	log.Println("reset gylist ", mgr.gylist)
}

func (mgr *GatewayManager) Update(id uint32, online uint32) {

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	if _, ok := mgr.gylist[id]; !ok {
		return
	}

	mgr.gylist[id].Online = online
}
