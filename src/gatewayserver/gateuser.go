package main

import (
	"command"
	"log"
)

type GateUser struct {
	account string
	userid  uint64

	//角色列表
	rolelist []*command.GateUserInfo

	task *GatewayTask
}

func NewGateUser(task *GatewayTask) *GateUser {
	user := &GateUser{
		task:    task,
		account: task.GetAccount(),
	}

	user.init()

	return user
}

func (user *GateUser) init() {

}

//通知大厅初始化玩家信息
func (user *GateUser) RegFromRecord(msg *command.RetGateRegUser) bool {

	//注册失败
	if msg.Retcode != 0 {
		return false
	}

	user.rolelist = msg.Userlist

	snd := new(command.NotifyChardata)
	snd.Rolelist = user.rolelist

	return true
}

//通知大厅玩家退出
func (user *GateUser) Unreg() bool {

	log.Println("网关用户:", user.account, ",", user.userid, "退出，通知卸载")

	gateUserManager.Erase(user)

	/*
		msg := new(command.SessionUnregUser)
		msg.Account = user.account
		msg.Userid = user.userid

		routeClientManager.Broadcast(GatewayServerInstance().GetServerType(), GatewayServerInstance().GetServerIndex(), command.SERVER_TYPE_FUNCTIONSERVER, 0, msg)
	*/

	return true
}

func (user *GateUser) GetAccount() string {
	return user.account
}

func (user *GateUser) GetID() uint64 {
	return user.userid
}
