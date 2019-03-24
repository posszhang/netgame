package main

import (
	"sync"
)

type GateUserManager struct {
	accountMap map[string]*GateUser
	idMap      map[uint64]*GateUser

	mutex sync.Mutex
}

func NewGateUserManager() *GateUserManager {
	mgr := &GateUserManager{
		accountMap: make(map[string]*GateUser),
		idMap:      make(map[uint64]*GateUser),
	}

	mgr.init()

	return mgr
}

func (mgr *GateUserManager) init() {

}

func (this *GateUserManager) Insert(user *GateUser) bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	_, ok := this.accountMap[user.GetAccount()]
	if ok {
		return false
	}
	this.accountMap[user.GetAccount()] = user

	_, ok = this.idMap[user.GetID()]
	if ok {
		delete(this.accountMap, user.GetAccount())
		return false
	}
	this.idMap[user.GetID()] = user

	return true
}

func (this *GateUserManager) Erase(user *GateUser) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	_, ok := this.accountMap[user.GetAccount()]
	if ok {
		delete(this.accountMap, user.GetAccount())
	}

	_, ok = this.idMap[user.GetID()]
	if ok {
		delete(this.idMap, user.GetID())
	}
}

func (this *GateUserManager) GetByID(id uint64) *GateUser {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	user, ok := this.idMap[id]
	if !ok {
		return nil
	}

	return user
}

func (this *GateUserManager) GetByAccount(account string) *GateUser {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	user, ok := this.accountMap[account]
	if !ok {
		return nil
	}

	return user
}

func (this *GateUserManager) GetNum() uint32 {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	return uint32(len(this.accountMap))
}
