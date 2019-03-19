package gnet

import (
	"command"
	kcp "github.com/xtaci/kcp-go"
	"math/rand"
	"runtime"
	"time"
)

type UDPNetServiceEvent interface {
	Init() bool
	Final()
	NewUDPTask(conn *kcp.UDPSession, port int)
}

type UDPNetService struct {
	udpServer  *UDPServer
	Event      UDPNetServiceEvent
	terminate  bool
	serverid   int
	servertype int
	port       int
}

func (this *UDPNetService) NewUDPTask(conn *kcp.UDPSession, port int) {

	this.Event.NewUDPTask(conn, port)
}

func (this *UDPNetService) init() (ret bool) {

	runtime.GOMAXPROCS(8)

	rand.Seed(time.Now().Unix())

	this.udpServer = new(UDPServer)

	if this.Event == nil || this.Event.Init() == false {
		return false
	}
	return true
}

func (this *UDPNetService) final() {
	this.Event.Final()
}

func (this *UDPNetService) Bind(name string, ip string, port int) bool {

	err := this.udpServer.bind(name, ip, port)
	if err != nil {
		return false
	}

	this.port = port
	return true
}

func (this *UDPNetService) serviceCallback() (ret bool) {
	conn, _ := this.udpServer.accept()
	if conn != nil {
		this.NewUDPTask(conn, this.port)
	}

	return true
}

func (this *UDPNetService) isTerminate() (ret bool) {
	return this.terminate
}

func (this *UDPNetService) Terminate() {
	if this.udpServer != nil {
		this.udpServer.Terminate()
	}

	this.terminate = true
}

func (this *UDPNetService) Run() {
	if this.init() {
		for !this.isTerminate() {
			if !this.serviceCallback() {
				break
			}
		}
	}

	this.final()
}

func (this *UDPNetService) GetIp() string {
	return this.udpServer.Ip
}

func (this *UDPNetService) GetPort() int {
	return this.udpServer.Port
}

func (this *UDPNetService) GetServerID() int {
	return this.serverid
}

func (this *UDPNetService) SetServerID(id int) {
	this.serverid = id
}

func (this *UDPNetService) SetServerType(servertype int) {
	this.servertype = servertype
}

func (this *UDPNetService) GetServerType() int {
	return this.servertype
}

func (this *UDPNetService) GetServerIndex() int {
	return this.serverid
}

func (this *UDPNetService) GetServerInfo() *command.ServerInfo {
	info := new(command.ServerInfo)
	info.Id = uint32(this.GetServerID())
	info.Type = uint32(this.GetServerType())
	info.Ip = this.GetIp()
	info.Port = uint32(this.GetPort())

	return info
}
