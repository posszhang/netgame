package main

import (
	"base/log"
	"command"
	"github.com/golang/protobuf/proto"
	"sync"
)

type ServerManager struct {
	mutex     sync.Mutex
	serverMap map[uint32]*ServerTask
}

func NewServerManager() *ServerManager {
	mgr := &ServerManager{
		serverMap: make(map[uint32]*ServerTask),
	}

	return mgr
}

func (mgr *ServerManager) GetByID(id uint32) *ServerTask {

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	if _, ok := mgr.serverMap[id]; !ok {
		return nil
	}

	return mgr.serverMap[id]
}

func (mgr *ServerManager) GetByType(tp uint32) []*ServerTask {

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	serverlist := make([]*ServerTask, 0, 0)
	for _, task := range mgr.serverMap {

		if task.GetServerType() != tp {
			continue
		}
		serverlist = append(serverlist, task)
	}

	return serverlist
}

func (mgr *ServerManager) UniqueAdd(task *ServerTask) bool {

	if mgr.GetByID(task.GetID()) != nil {
		log.Println(task.GetServerInfo(), "重复添加服务器")
		return false
	}

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	mgr.serverMap[task.GetID()] = task

	return true
}

func (mgr *ServerManager) UniqueRemove(task *ServerTask) {

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	delete(mgr.serverMap, task.GetID())
}

// 发送消息给全部server，排除指定类型服务器
func (mgr *ServerManager) SendCmdToAllExceptionType(msg proto.Message, tp uint32) {

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	for _, task := range mgr.serverMap {

		if task.GetServerType() == tp {
			continue
		}

		task.SendCmd(msg)
	}
}

//通知服务器初始化路由列表
func (mgr *ServerManager) NotifyRouteServerInit(task *ServerTask) {

	msg := new(command.NotifyRouteServerInit)

	serverlist := mgr.GetByType(command.RouteServer)
	msg.Serverlist = make([]*command.ServerInfo, 0, len(serverlist))

	for _, server := range serverlist {
		info := server.GetServerInfo()
		msg.Serverlist = append(msg.Serverlist, info)
	}

	task.SendCmd(msg)
}

//通知所有服务器新增路由服
func (mgr *ServerManager) NotifyRouteServerAdd(task *ServerTask) {

	msg := new(command.NotifyRouteServerAdd)
	info := task.GetServerInfo()

	msg.Info = info

	mgr.SendCmdToAllExceptionType(msg, command.RouteServer)
}

func (mgr *ServerManager) CheckCareAndNotify(task *ServerTask) {

	for _, server := range serverMap {

		servertype := task.GetServerType()
		if !server.IsCare(servertype) {
			continue
		}

		//关心服务器
		mgr.NotifyServerCare
	}
}

// 获取我关心的服务器列表
func (mgr *ServerManager) GetCareList(task *ServerTask) []*ServerTask {

	serverlist := make([]*ServerTask, 0, 0)

	for _, server := range serverMap {

		if !task.IsCare(server.GetServerType()) {
			continue
		}

		serverlist = append(serverlist, server)
	}

	return serverlist
}
