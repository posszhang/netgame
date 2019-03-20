package main

import (
	"base/log"
	"base/util"
	"command"
	"sync"
)

type GatewayTaskManager struct {
	mutex sync.Mutex

	//key account
	taskMap map[string]*GatewayTask

	//timer
	_10_sec *util.Timer
}

func NewGatewayTaskManager() *GatewayTaskManager {
	mgr := &GatewayTaskManager{}

	mgr.init()

	return mgr
}

func (mgr *GatewayTaskManager) init() {

	mgr._10_sec = util.NewTimer(10000)
}

func (mgr *GatewayTaskManager) UniqueAdd(task *GatewayTask) bool {

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	if _, ok := mgr.taskMap[task.GetAccount()]; ok {
		log.Println("账号", task.GetAccount(), "唯一性添加失败，重复登陆")
		return false
	}

	mgr.taskMap[task.GetAccount()] = task

	return true
}

func (mgr *GatewayTaskManager) UniqueRemove(task *GatewayTask) {

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	delete(mgr.taskMap, task.GetAccount())
}

func (mgr *GatewayTaskManager) Count() uint32 {

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	count := uint32(len(mgr.taskMap))

	return count
}

func (mgr *GatewayTaskManager) TimeAction() {

	if mgr._10_sec.Check() {

		snd := new(command.UpdateGatewayOnline)
		snd.Id = uint32(service.GetServerID())
		snd.Online = mgr.Count()

		//广播所有的登陆服务器
		routeManager.BroadcastByType(command.LoginServer, snd)
	}
}
