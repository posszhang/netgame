package gnet

import (
	"base/log"
	"command"
	"github.com/golang/protobuf/proto"
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

	messageHandler MessageHandler
	messageQueue   MessageQueue

	//当前处理消息的来源服务器
	srcid uint32

	Derived IRouteManager
}

func (mgr *RouteManager) InitRouteManager() {
	mgr.messageHandler.Reg(&command.RouteMessage{}, mgr.onRouteMessage)

	mgr.messageHandler.Reg(&command.RouteBroadcastByType{}, mgr.onBroadcastByType)
	mgr.messageHandler.Reg(&command.RouteBroadcastByID{}, mgr.onBroadcastByID)
}

func (mgr *RouteManager) Destory() {
	mgr.messageQueue.Final()
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

	mgr.messageQueue.Do(mgr.process)
	//mgr.messageQueue.Do(mgr.Derived.MsgParse)
}

func (mgr *RouteManager) process(cmd *command.Message) bool {

	if mgr.messageHandler.Process(cmd) {
		return false
	}

	if mgr.Derived == nil {
		return false
	}

	mgr.Derived.MsgParse(cmd)

	return true
}

func (mgr *RouteManager) onRouteMessage(cmd proto.Message) {

	msg := cmd.(*command.RouteMessage)
	mgr.srcid = msg.Srcid

	mgr.Derived.MsgParse(msg.Msg)
}

func (mgr *RouteManager) onBroadcastByType(cmd proto.Message) {

	msg := cmd.(*command.RouteBroadcastByType)
	mgr.srcid = msg.Srcid

	mgr.Derived.MsgParse(msg.Msg)
}

func (mgr *RouteManager) onBroadcastByID(cmd proto.Message) {

	msg := cmd.(*command.RouteBroadcastByID)
	mgr.srcid = msg.Srcid

	mgr.Derived.MsgParse(msg.Msg)

}

func (mgr *RouteManager) PackRouteMessage(cmd proto.Message) *command.Message {

	return PackMessage(cmd)

	msg := new(command.RouteMessage)
	msg.Srcid = mgr.Derived.GetServerInfo().Id
	msg.Destid = 0
	msg.Msg = PackMessage(cmd)

	return PackMessage(msg)
}

func (mgr *RouteManager) BroadcastByType(servertype uint32, msg proto.Message) {

	route := mgr.GetRouteByType(servertype)
	if route == nil {
		log.Println("没有可用的路由服务器")
		return
	}

	snd := new(command.RouteBroadcastByType)
	snd.Srcid = mgr.Derived.GetServerInfo().Id
	snd.Type = servertype
	snd.Msg = mgr.PackRouteMessage(msg)

	route.SendCmd(snd)
}

func (mgr *RouteManager) BroadcastByID(destid uint32, msg proto.Message) {

	servertype := uint32(command.ServerID2Type(int(destid)))
	route := mgr.GetRouteByType(servertype)
	if route == nil {
		log.Println("没有可用的路由服务器")
		return
	}

	snd := new(command.RouteBroadcastByID)
	snd.Srcid = mgr.Derived.GetServerInfo().Id
	snd.Id = destid
	snd.Msg = mgr.PackRouteMessage(msg)

	route.SendCmd(snd)
}

func (mgr *RouteManager) Broadcast(servertype uint32, index uint32, msg proto.Message) {

	destid := uint32(command.GetServerID(int(servertype), int(index)))
	mgr.BroadcastByID(destid, msg)
}

func (mgr *RouteManager) GetRouteByType(servertype uint32) *RouteClient {

	if len(mgr.routeList) == 0 {
		return nil
	}

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	index := servertype % uint32(len(mgr.routeList))

	return mgr.routeList[index]
}

func (mgr *RouteManager) SendCmd(cmd proto.Message) {

	mgr.BroadcastByID(mgr.srcid, cmd)
}

func (mgr *RouteManager) Close() {

	for _, route := range mgr.routeList {
		route.Terminate()
		route.Join()
	}
}
