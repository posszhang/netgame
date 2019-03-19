package gnet

import (
	"base/log"
	"github.com/gorilla/websocket"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

type IWSNetService interface {
	Init() bool
	Final()
	NewTCPTask(conn net.Conn, srcPort int)
	NewWSTask(conn *websocket.Conn, srcPort int)
}

type WSNetService struct {
	tcpServer *TCPServer
	wsServer  *WebSocketServer
	terminate bool
	serverid  int

	Derived IWSNetService
}

func (this *WSNetService) init() (ret bool) {

	runtime.GOMAXPROCS(runtime.NumCPU())
	//初始化时间种子
	rand.Seed(time.Now().Unix())

	if this.Derived == nil || this.Derived.Init() == false {
		return false
	}
	return true
}

func (this *WSNetService) final() {
	this.Derived.Final()
}

func (this *WSNetService) Bind(name string, ip string, port int, ws bool) bool {

	if ws {
		this.wsServer = new(WebSocketServer)
		err := this.wsServer.bind(name, ip, port)
		if err != nil {
			log.Println("监听端口", port, "失败")
			return false
		}

		go func() {
			log.Println("启动监听协程", port)
			for {
				if this.isTerminate() {
					return
				}
				conn := this.wsServer.accept()
				if conn != nil {
					this.Derived.NewWSTask(conn, port)
				}
			}
			log.Println("释放监听用户连接:", port)
		}()

		return true
	}

	this.tcpServer = new(TCPServer)
	err := this.tcpServer.bind(name, ip, port)
	if err != nil {
		log.Println("监听端口", port, "失败")
		return false
	}

	go func() {

		log.Println("启动监听协程", port)
		for {
			if this.isTerminate() {
				return
			}
			conn, _ := this.tcpServer.accept()
			if conn != nil {
				this.Derived.NewTCPTask(conn, port)
			}
		}
		log.Println("释放监听用户连接:", port)
	}()

	return true
}

func (this *WSNetService) serviceCallback() (ret bool) {

	time.Sleep(time.Millisecond * time.Duration(10))

	return true
}

func (this *WSNetService) isTerminate() (ret bool) {
	return this.terminate
}

func (this *WSNetService) Terminate() {

	if this.tcpServer != nil {
		this.tcpServer.Terminate()
		this.tcpServer = nil
	}

	if this.wsServer != nil {
		this.wsServer.Terminate()
		this.wsServer = nil
	}

	this.terminate = true
}

func (this *WSNetService) Run() {

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGPIPE, syscall.SIGHUP)
	go func() {
		for sig := range ch {
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

func (this *WSNetService) GetIp() string {
	return this.wsServer.Ip
}

func (this *WSNetService) GetPort() int {
	return this.wsServer.Port
}

func (this *WSNetService) GetServerID() int {
	return this.serverid
}

func (this *WSNetService) SetServerID(id int) {
	this.serverid = id
}
