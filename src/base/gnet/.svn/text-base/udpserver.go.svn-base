package gnet

import (
	"fmt"
	kcp "github.com/xtaci/kcp-go"
	//	"strings"
	//	"time"
	"base/log"
)

/*
	这是一个可以绑定多个端口的UDPServer
*/
type UDPServer struct {
	listener  *kcp.Listener
	terminate bool

	Ip   string
	Port int
}

func (this *UDPServer) bind(name string, ip string, port int) (err error) {

	//ipstr := fmt.Sprintf(":%d", port)

	this.listener, err = kcp.ListenWithOptions(fmt.Sprint(":", port), nil, 0, 0)
	if nil != err {
		return err
	}

	log.Println("listening on:", this.listener.Addr())

	//this.Ip = Global["ifname"] //GetIp(Global["ifname"])
	this.Ip = ip
	this.Port = port

	log.Println("监听端口:", this.Ip, ":", this.Port)

	return err
}

func (this *UDPServer) accept() (*kcp.UDPSession, error) {

	this.listener.SetReadBuffer(4 * 1024 * 1024)
	this.listener.SetWriteBuffer(4 * 1024 * 1024)
	this.listener.SetDSCP(46)
	//this.listener.SetDeadline(time.Now().Add(time.Second * 1))
	conn, err := this.listener.AcceptKCP()

	if err != nil {
		return nil, err
	}

	//conn.SetKeepAlive(true)
	//conn.SetKeepAlivePeriod(1 * time.Minute)
	conn.SetNoDelay(0, 40, 0, 0)
	conn.SetWriteBuffer(128 * 1024)
	conn.SetReadBuffer(128 * 1024)

	return conn, nil
}

func (this *UDPServer) Terminate() {

	this.terminate = true

	if this.listener != nil {
		this.listener.Close()
		fmt.Println("关闭监听端口")
	}

}

func (this *UDPServer) IsTerminate() bool {
	return this.terminate
}
