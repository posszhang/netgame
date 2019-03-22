package main

import (
	"base/log"
	"base/util"
	"sync"
)

type LoginTaskManager struct {
	mutex sync.Mutex

	//key account
	taskMap map[uint32]*LoginTask

	//timer
	_10_sec *util.Timer
}

func NewLoginTaskManager() *LoginTaskManager {
	mgr := &LoginTaskManager{
		taskMap: make(map[uint32]*LoginTask),
	}

	mgr.init()

	return mgr
}

func (mgr *LoginTaskManager) init() {

	mgr._10_sec = util.NewTimer(10000)
}

func (mgr *LoginTaskManager) UniqueAdd(task *LoginTask) bool {

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	if _, ok := mgr.taskMap[task.GetTempID()]; ok {
		log.Println("账号", task.GetTempID(), "唯一性添加失败，重复登陆")
		return false
	}

	mgr.taskMap[task.GetTempID()] = task

	return true
}

func (mgr *LoginTaskManager) UniqueRemove(task *LoginTask) {

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	delete(mgr.taskMap, task.GetTempID())
}

func (mgr *LoginTaskManager) Count() uint32 {

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	count := uint32(len(mgr.taskMap))

	return count
}

func (mgr *LoginTaskManager) TimeAction() {

	if mgr._10_sec.Check() {

	}
}
