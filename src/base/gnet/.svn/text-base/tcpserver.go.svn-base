package gnet

import (
	"base/log"
	"fmt"
	"net"
	"time"
)

/*
这是一个可以绑定多个端口的TCPServer
*/

type TCPServer struct {
	listener  *net.TCPListener
	terminate bool

	Ip   string
	Port int
}

func (this *TCPServer) bind(name string, ip string, port int) (err error) {

	ipstr := fmt.Sprintf(":%d", port)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", ipstr)
	if err != nil {
		log.Println("[服务] 解析失败 ", ipstr)
		return err
	}

	this.listener, err = net.ListenTCP("tcp", tcpAddr)
	if nil != err {
		return err
	}

	this.Ip = ip
	this.Port = int(this.listener.Addr().(*net.TCPAddr).Port)

	log.Println(name, "监听端口:", this.Ip, ":", this.Port)

	return err
}

func (this *TCPServer) accept() (*net.TCPConn, error) {

	this.listener.SetDeadline(time.Now().Add(time.Second * 1))
	conn, err := this.listener.AcceptTCP()

	if err != nil {
		return nil, err
	}

	conn.SetKeepAlive(true)
	conn.SetKeepAlivePeriod(1 * time.Minute)
	conn.SetNoDelay(true)
	conn.SetWriteBuffer(128 * 1024)
	conn.SetReadBuffer(128 * 1024)

	return conn, nil
}

func (this *TCPServer) Terminate() {

	this.terminate = true

	if this.listener != nil {
		this.listener.Close()
		fmt.Println("关闭监听端口")
	}

}

func (this *TCPServer) IsTerminate() bool {
	return this.terminate
}
