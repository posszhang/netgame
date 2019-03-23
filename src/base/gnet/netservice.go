package gnet

import (
	"base/log"
	"command"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

type INetService interface {
	Init() bool
	Final()
	NewTCPTask(conn net.Conn, port int)
}

type NetService struct {
	tcpServer *TCPServer
	terminate bool
	serverid  int
	port      int

	Derived INetService
}

func (this *NetService) NewTCPTask(conn net.Conn, port int) {

	this.Derived.NewTCPTask(conn, port)
}

func (this *NetService) init() (ret bool) {

	runtime.GOMAXPROCS(runtime.NumCPU())
	//初始化时间种子
	rand.Seed(time.Now().Unix())

	this.tcpServer = new(TCPServer)

	if this.Derived == nil || this.Derived.Init() == false {
		return false
	}
	return true
}

func (this *NetService) final() {
	this.Derived.Final()
}

func (this *NetService) Bind(name string, ip string, port int) bool {

	err := this.tcpServer.bind(name, ip, port)
	if err != nil {
		return false
	}

	this.port = port
	return true
}

func (this *NetService) serviceCallback() (ret bool) {
	conn, _ := this.tcpServer.accept()
	if conn != nil {
		this.NewTCPTask(conn, this.port)
	}

	return true
}

func (this *NetService) isTerminate() (ret bool) {
	return this.terminate
}

func (this *NetService) Terminate() {
	if this.tcpServer != nil {
		this.tcpServer.Terminate()
	}

	this.terminate = true
}

func (this *NetService) Run() {

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGPIPE, syscall.SIGHUP)
	go func() {
		for sig := range ch {

			//忽略，对应关闭的conn多次write
			if sig == syscall.SIGPIPE {
				continue
			}

			this.Terminate()
			log.Println("[服务] 收到信号 ", sig)
			break
		}
	}()

	if !this.init() {
		return
	}

	for !this.isTerminate() {
		if !this.serviceCallback() {
			break
		}
	}

	this.final()
}

func (this *NetService) GetIp() string {
	return this.tcpServer.Ip
}

func (this *NetService) GetPort() int {
	return this.tcpServer.Port
}

func (this *NetService) GetServerID() int {
	return this.serverid
}

func (this *NetService) SetServerID(id int) {
	this.serverid = id
}

func (this *NetService) GetServerType() int {
	return command.ServerID2Type(this.serverid)
}

func (this *NetService) GetServerIndex() int {
	return command.ServerID2Index(this.serverid)
}

func (this *NetService) GetServerInfo() *command.ServerInfo {
	info := new(command.ServerInfo)
	info.Id = uint32(this.GetServerID())
	info.Index = uint32(this.GetServerIndex())
	info.Type = uint32(this.GetServerType())
	info.Ip = this.GetIp()
	info.Port = uint32(this.GetPort())

	return info
}
