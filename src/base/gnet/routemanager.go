package gnet

import (
	"command"
	"sync"
)

/*
	为保证多路由服线程安全
	使用单消息逻辑队列
	上次必须保证tick执行Do
*/
type IRouteManager interface {
	GetServerInfo() *command.ServerInfo
	MsgParse(msg *command.Message) bool
}

type RouteManager struct {
	routeList []*RouteClient
	mutex     sync.Mutex

	messageQueue MessageQueue

	Derived IRouteManager
}

func (mgr *RouteManager) InitRouteList(routeList []*command.ServerInfo) {

	for _, info := range routeList {
		mgr.Add(info)
	}
}

func (mgr *RouteManager) Add(info *command.ServerInfo) {

	client := NewRouteClient(mgr)
	client.Connect(info.Ip, int(info.Port))

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	if mgr.routeList == nil {
		mgr.routeList = make([]*RouteClient, 0, 0)
	}

	mgr.routeList = append(mgr.routeList, client)
}

func (mgr *RouteManager) MsgParse(msg *command.Message) bool {

	mgr.messageQueue.Cache(msg)

	return true
}

//上次必须tick该函数
func (mgr *RouteManager) Do() {

	if mgr.Derived == nil {
		return
	}

	mgr.messageQueue.Do(mgr.Derived.MsgParse)
}
