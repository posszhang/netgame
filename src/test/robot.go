package main

import (
	"base/gnet"
	"base/log"
	"command"
	"fmt"
	"github.com/golang/protobuf/proto"
	"time"
)

const (
	ROBOT_LOGIN   = 1
	ROBOT_GATEWAY = 2
)

type Robot struct {
	gnet.TCPClient

	loginip   string
	loginport int
	account   string
	password  string
	session   string
	//1:login,2:gateway
	status uint32

	msgHandler gnet.MessageHandler
	msgQueue   gnet.MessageQueue

	chClosed chan bool
	inittm   int64
}

func NewRobot(ip string, port int, account string, password string) *Robot {
	robot := &Robot{
		loginip:   ip,
		loginport: port,
		account:   account,
		password:  password,
		chClosed:  make(chan bool, 1),
		inittm:    time.Now().Unix(),
	}

	robot.init()

	return robot
}

func (robot *Robot) init() {

	robot.Derived = robot

	robot.msgHandler.Reg(&command.RetUserVerify{}, robot.onRetUserVerify)
	robot.msgHandler.Reg(&command.RetUserLogin{}, robot.onRetUserLogin)
	robot.msgHandler.Reg(&command.RetGatewayLogin{}, robot.onRetGatewayLogin)
	robot.msgHandler.Reg(&command.TestBroadcastAll{}, robot.onTestBroadcastAll)

}

func (robot *Robot) Run() {

	robot.status = ROBOT_LOGIN
	ret := robot.Connect(robot.loginip, robot.loginport)

	if !ret {
		log.Println("机器人", robot.account, "连接登陆服务器失败")
		return
	}

	robot.inittm = time.Now().Unix()
}

func (robot *Robot) OnConnected() {

	log.Println("connected", robot.status)

	if robot.status == ROBOT_LOGIN {

		snd := new(command.ReqUserVerify)
		robot.SendCmd(snd)
		log.Println("请求登录服验证")

	} else if robot.status == ROBOT_GATEWAY {

		snd := new(command.ReqGatewayLogin)
		snd.Session = robot.session
		robot.SendCmd(snd)
		log.Println("请求网关验证")
	}
}

func (robot *Robot) MsgParse(msg *command.Message) bool {

	robot.msgQueue.Cache(msg)

	return true
}

func (robot *Robot) Do() {

	robot.msgQueue.Do(robot.msgHandler.Process)
}

func (robot *Robot) onRetUserVerify(cmd proto.Message) {

	snd := new(command.ReqUserLogin)
	snd.Loginstr = fmt.Sprint("account=", robot.account, "&password=", robot.password)
	robot.SendCmd(snd)
}

func (robot *Robot) onRetUserLogin(cmd proto.Message) {
	msg := cmd.(*command.RetUserLogin)
	if msg.Retcode != 0 {
		log.Println("机器人", robot.account, "登陆失败,错误码", msg.Retcode)
		return
	}

	log.Println("机器人登陆成功，session=", msg.Session)
	robot.session = msg.Session

	robot.Terminate()
	robot.Join()

	robot.status = ROBOT_GATEWAY
	ret := robot.Connect(msg.Ip, int(msg.Port))
	if !ret {
		log.Println("机器人", robot.account, "登陆网关服务器失败", msg.Ip, ":", msg.Port)
		return
	}

}

func (robot *Robot) onRetGatewayLogin(cmd proto.Message) {

	log.Println("网关验证成功")

	go func() {

		for {
			snd := new(command.TestBroadcastAll)
			snd.Str = "1111111111111111111111111111111111111111111111111111111111111111111111111111111111"
			robot.SendCmd(snd)

			time.Sleep(100 * time.Millisecond)

			flag := false
			select {
			case flag = <-robot.chClosed:
				log.Println("机器人准备退出")
				break
			default:
				break
			}

			if flag {
				break
			}
		}

	}()
}

func (robot *Robot) Close() {
	robot.chClosed <- true
	robot.Terminate()
	robot.Join()
	robot.msgQueue.Final()
}

func (robot *Robot) GetInitSec() uint32 {
	return uint32(time.Now().Unix() - robot.inittm)
}

func (robot *Robot) onTestBroadcastAll(cmd proto.Message) {
	//msg := cmd.(*command.TestBroadcastAll)

	//log.Println("recv:", msg.Str)
}
