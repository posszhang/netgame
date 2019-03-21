package gnet

import (
	"base/log"
	"base/util"
	"bytes"
	"command"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"net"
	"sync"
	"time"
)

type ITCPClient interface {
	MsgParse(msg *command.Message) bool
	OnConnected()
}

type TCPClient struct {
	Derived ITCPClient
	conn    net.Conn

	sndQueue chan []byte
	sndCache bool

	wg        sync.WaitGroup
	mutex     sync.Mutex
	terminate bool

	ip     string
	port   int
	reconn bool
}

func (this *TCPClient) Connect(ip string, port int) (ret bool) {

	this.ip = ip
	this.port = port

	this.reconn = false
	this.terminate = true

	if !this.connect() {
		this.reconnect()

		return true
	}

	this.terminate = false
	this.GoHandler()

	if this.Derived != nil {
		this.Derived.OnConnected()
	}

	return true
}

func (this *TCPClient) connect() bool {

	s := fmt.Sprintf("%s:%d", this.ip, this.port)
	var err error
	this.conn, err = net.Dial("tcp", s)
	if err != nil {
		log.Println("connect error:%s", err.Error())

		return false
	}

	return true
}

func (this *TCPClient) reconnect() bool {

	log.Println("启动重连(", this.ip, ":", this.port, ")")

	go func() {
		defer func() {
			log.Println("重连成功(", this.ip, ":", this.port, ")")
		}()

		for i := 1; ; i++ {
			time.Sleep(10 * time.Second)
			log.Println("第", i, "次重连:", this.ip, ":", this.port)
			if this.connect() {
				break
			}
		}

		this.mutex.Lock()
		this.terminate = false
		this.mutex.Unlock()

		this.GoHandler()

		if this.Derived != nil {
			this.Derived.OnConnected()
		}

	}()

	return true
}

func (this *TCPClient) GoHandler() {

	if this.conn == nil {
		return
	}

	this.sndCache = true
	this.sndQueue = make(chan []byte, 1280)

	go this.reader()
	go this.writer()
}

func (this *TCPClient) reader() {

	this.wg.Add(1)
	defer this.wg.Done()

	header := make([]byte, 4)

	for {

		//读取数据
		if _, err := io.ReadFull(this.conn, header); err != nil {
			this.TerminateReconn()
			break
		}

		size := binary.BigEndian.Uint32(header)
		if size > MAX_SOCKETDATASIZE {
			break
		}

		data := make([]byte, size)
		if _, err := io.ReadFull(this.conn, data); err != nil {
			this.TerminateReconn()
			break
		}

		this.doCmd(data)
	}

	/*
		this.mutex.Lock()
		defer this.mutex.Unlock()

		if this.reconn {
			this.reconnect()
		}
		this.terminate = true
		this.reconn = false

		log.Println("reader defer")
	*/
	this.Terminate()
}

func (this *TCPClient) writer() {

	this.wg.Add(1)
	defer this.wg.Done()

	if this.sndQueue == nil {
		log.Println("直接退出")
		return
	}

	for {

		buf := <-this.sndQueue

		if buf == nil {
			break
		}

		//this.SendCmd_NoBuf(buf)

		_, err := this.conn.Write(buf)
		if err != nil {
			break
		}
	}

	log.Println("writer defer")

	time.Sleep(100 * time.Millisecond)
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.conn.(*net.TCPConn).SetLinger(0)
	this.conn.Close()
	close(this.sndQueue)
	this.terminate = true

	if this.reconn {
		this.reconnect()
	}
}

func (this *TCPClient) SendCmd_NoBuf(buf []byte) (ret bool) {

	packetLen := len(buf)
	current := 0

	for current < packetLen {
		n, err := this.conn.Write(buf[current:packetLen])

		if err != nil {
			this.TerminateReconn()
			return false
		}

		current += n
	}

	return true
}

func (this *TCPClient) SendCmd(cmd proto.Message) (ret bool) {

	msg := new(command.Message)
	msg.Name = proto.MessageName(cmd)
	msg.Type = uint32(util.BKDRHash(msg.Name))
	msg.Index = 0
	msg.Data, _ = proto.Marshal(cmd)

	return this.SendCmd_NoPack(msg)

}

func (this *TCPClient) SendCmd_NoPack(msg *command.Message) bool {

	if this.IsTerminate() {
		return false
	}

	d := make([][]byte, 2)
	d[1], _ = proto.Marshal(msg)
	d[0] = util.Int2Byte(len(d[1]))

	f := []byte("")
	g := bytes.Join(d, f)

	if this.sndCache {

		this.mutex.Lock()
		defer this.mutex.Unlock()

		if !this.terminate {

			if len(this.sndQueue) == cap(this.sndQueue) {
				log.Println("close conn: channel full")
			}

			this.sndQueue <- g
		}
		return true
	}

	this.SendCmd_NoBuf(g)
	return true
}

func (this *TCPClient) Terminate() {

	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.terminate {
		return
	}

	log.Println("主动断开，不需要重连")
	this.reconn = false
	this.terminate = true
	this.sndQueue <- nil
}

func (this *TCPClient) IsTerminate() bool {

	this.mutex.Lock()
	defer this.mutex.Unlock()

	return this.terminate
}

func (this *TCPClient) Join() {

	log.Println("进入TCPClient Join")

	/*
		this.mutex.Lock()


		//等待数据全部发送
		if this.sndQueue != nil {

			for {
				select {
				case <-this.sndQueue:
					continue
				default:
					break
				}
			}

			close(this.sndQueue)
			this.sndQueue = nil
		}

		this.mutex.Unlock()
	*/
	this.Terminate()

	log.Println("等待协程退出")
	//等待退出
	this.wg.Wait()
	log.Println("协程全部退出")
	time.Sleep(100 * time.Millisecond)
}

func (this *TCPClient) doCmd(buf []byte) {

	msg := new(command.Message)
	err := proto.Unmarshal(buf, msg)
	if err != nil {
		return
	}

	if len(msg.Name) != 0 && msg.Name == "command.PingMsg" {
		msg.Type = CMD_PING
	}

	//心跳包，处理
	if msg.Type == CMD_PING {

		this.pong()
		return
	}

	if this.Derived != nil {
		this.Derived.MsgParse(msg)
	}
}

func (this *TCPClient) TerminateReconn() {

	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.terminate {
		return
	}

	log.Println("系统断开，需要重连")
	this.reconn = true
	this.terminate = true
	this.sndQueue <- nil
}

func (this *TCPClient) pong() {

	msg := new(command.PongMsg)
	this.SendCmd(msg)
}
