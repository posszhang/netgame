package gnet

import (
	"base/log"
	"fmt"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
)

//var upgrader = websocket.Upgrader{} // use default options

type WebSocketServer struct {
	listener  net.Listener
	terminate bool
	conns     chan *websocket.Conn

	Ip   string
	Port int
}

func (this *WebSocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	origin := w.Header().Get("Origin")
	if len(origin) > 0 {
		origin = "http://" + r.Host + "/"
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	c, err := websocket.Upgrade(w, r, nil, 4096, 4096)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	this.conns <- c
}

func (this *WebSocketServer) bind(name string, ip string, port int) (err error) {
	this.conns = make(chan *websocket.Conn)
	ipstr := fmt.Sprintf(":%d", port)
	this.listener, err = net.Listen("tcp", ipstr)

	svr := http.Server{Handler: this}
	go svr.Serve(this.listener)

	this.Ip = ip
	this.Port = int(this.listener.Addr().(*net.TCPAddr).Port)

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return err
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				this.Ip = ipnet.IP.String()
			}

		}
	}

	log.Println("监听端口:", this.Ip, ":", this.Port)

	return nil
}

func (this *WebSocketServer) accept() *websocket.Conn {

	/*
		select {
		case conn := <-this.conns:
			return conn
		default:
			return nil
		}

		return nil
	*/
	conn := <-this.conns
	return conn
}

func (this *WebSocketServer) Terminate() {
	this.listener.Close()
	close(this.conns)
}
