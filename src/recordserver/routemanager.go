package main

import (
	"base/gnet"
	"base/log"
	"command"
	"github.com/golang/protobuf/proto"
	"gopkg.in/mgo.v2/bson"
)

type RouteManager struct {
	gnet.RouteManager

	msgHandler gnet.MessageHandler
}

func NewRouteManager() *RouteManager {
	mgr := &RouteManager{}

	mgr.init()

	return mgr
}

func (mgr *RouteManager) GetServerInfo() *command.ServerInfo {
	return service.GetServerInfo()
}

func (mgr *RouteManager) MsgParse(msg *command.Message) bool {

	mgr.msgHandler.Process(msg)

	return true
}

func (mgr *RouteManager) init() {

	mgr.Derived = mgr

	mgr.msgHandler.Reg(&command.ReqGateRegUser{}, mgr.onReqGateRegUser)
}

func (mgr *RouteManager) onReqGateRegUser(cmd proto.Message) {

	msg := cmd.(*command.ReqGateRegUser)
	snd := new(command.RetGateRegUser)

	log.Println("网关用户登录查询", msg.Account)

	session := mongo.Clone()
	defer session.Close()

	c := session.DB("").C("userdata")

	userinfo := new(command.GateUserInfo)

	count, err := c.Find(bson.M{"account": msg.Account}).Count()
	if err != nil {
		log.Println("查询帐号数据错误", msg.Account, err)

		snd.Retcode = 1
		mgr.SendCmd(snd)

		return
	}

	// 新帐号
	if count == 0 {

		log.Println("初始化帐号", msg.Account)

	} else {

		err := c.Find(bson.M{"account": msg.Account}).One(&userinfo)
		if err != nil {
			log.Println("查询帐号数据错误", msg.Account, err)

			snd.Retcode = 1
			mgr.SendCmd(snd)
		}
	}

	snd.Retcode = 0
	snd.Userlist = make([]*command.GateUserInfo, 0, 0)
	snd.Userlist = append(snd.Userlist, userinfo)
	mgr.SendCmd(snd)
}
