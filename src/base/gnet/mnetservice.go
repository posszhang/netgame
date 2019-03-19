package gnet

import (
	"base/log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

type IMNetService interface {
	Init() bool
	Final()
	NewTCPTask(conn net.Conn, srcPort int)
}

type MNetService struct {
	tcpServer []*TCPServer
	terminate bool

	Derived IMNetService
}

func (this *MNetService) NewTCPTask(conn net.Conn, srcPort int) {

	this.Derived.NewTCPTask(conn, srcPort)
}

func (this *MNetService) init() (ret bool) {

	runtime.GOMAXPROCS(runtime.NumCPU())
	//初始化时间种子
	rand.Seed(time.Now().Unix())

	if this.Derived == nil || this.Derived.Init() == false {
		return false
	}
	return true
}

func (this *MNetService) final() {
	this.Derived.Final()
}

func (this *MNetService) Bind(name string, ip string, port int) bool {

	server := new(TCPServer)

	err := server.bind(name, ip, port)
	if err != nil {
		log.Println("监听端口", port, "失败")
		return false
	}

	this.tcpServer = append(this.tcpServer, server)

	go func() {

		log.Println("启动监听协程", port)
		for {
			if this.isTerminate() {
				return
			}
			log.Println("监控用户连接", port)
			conn, _ := server.accept()
			if conn != nil {
				this.NewTCPTask(conn, port)
			}
		}
		log.Println("释放监听用户连接:", port)
	}()

	return true
}

func (this *MNetService) serviceCallback() (ret bool) {

	time.Sleep(time.Millisecond * time.Duration(10))

	return true
}

func (this *MNetService) isTerminate() (ret bool) {
	return this.terminate
}

func (this *MNetService) Terminate() {

	for i := 0; i != len(this.tcpServer); i++ {
		server := this.tcpServer[i]
		if server != nil {
			server.Terminate()
		}
	}

	this.terminate = true
}

func (this *MNetService) Run() {

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
