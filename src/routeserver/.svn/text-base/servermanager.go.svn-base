package main

import (
	"base/log"
	//	"command"
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
