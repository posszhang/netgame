package gnet

import (
	"base/log"
	"base/util"
	"bytes"
	"command"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	kcp "github.com/xtaci/kcp-go"
	"io"
	"net"
	"runtime/debug"
	"sync"
	"time"
)

type UDPTaskEvent interface {
	VerifyConn(msg *command.Message) bool
	RecycleConn() bool
	MsgParse(msg *command.Message) bool
}

type UDPTask struct {
	Event UDPTaskEvent
	conn  *kcp.UDPSession

	sndQueue chan []byte

	terminate bool
	mutex     sync.Mutex

	//超时验证
	verify        bool
	verifyTimeout *time.Timer

	pingCount uint32
}

func (this *UDPTask) GoHandler(conn *kcp.UDPSession) {

	this.conn = conn
	this.pingCount = 0
	this.sndQueue = make(chan []byte, 1280)

	// 协程读
	go this.reader()
	// 协程写
	go this.writer()

	this.mutex.Lock()

	this.verifyTimeout = time.AfterFunc(time.Second*time.Duration(10), func() {
		if this.IsTerminate() {
			return
		}
		this.Terminate()

		this.mutex.Lock()
		defer this.mutex.Unlock()

		this.verifyTimeout.Stop()
		this.verifyTimeout = nil
		log.Println("task验证超时,关闭连接")
	})
	this.mutex.Unlock()
}

func (this *UDPTask) reader() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("[异常] ", err, "\n", string(debug.Stack()))
		}
	}()

	header := make([]byte, 4)
	var timeout time.Duration = 30

	for {

		//60秒读超时，发送心跳
		var err error
		this.conn.SetReadDeadline(time.Now().Add(timeout * time.Second))

		//读取数据
		if _, err = io.ReadFull(this.conn, header); err == nil {

			size := binary.BigEndian.Uint32(header)
			if size > MAX_SOCKETDATASIZE {
				log.Println("包体过大", size)
				break
			}

			data := make([]byte, size)
			if _, err = io.ReadFull(this.conn, data); err == nil {

				this.doCmd(data)
				continue
			}
		}

		//是超时错误,检查心跳
		if e, ok := err.(net.Error); ok && e.Timeout() {

			if this.pingCount <= 1 {
				this.ping()
				continue
			}

			log.Println("读取数据超时,心跳也无返回", err.Error())
		}

		break
	}

	// 异常结束
	this.Terminate()

	//验证成功的需要通知释放
	if this.isVerify() {
		this.Event.RecycleConn()
	}

}

func (this *UDPTask) writer() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("[异常] ", err, "\n", string(debug.Stack()))
		}
	}()

	for {

		//如果关闭了，可以读不能写
		buf := <-this.sndQueue

		if buf == nil {
			break
		}

		//this.sendCmd_NoBuf(buf)

		_, err := this.conn.Write(buf)
		if err != nil {
			break
		}
	}

	time.Sleep(100 * time.Millisecond)
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.conn.Close()
	close(this.sndQueue)
	this.terminate = true

}

func (this *UDPTask) doCmd(buf []byte) {

	msg := new(command.Message)

	err := proto.Unmarshal(buf, msg)
	if err != nil {
		log.Println("udptask 读取数据解析失败", buf, len(buf))
		//return
	}

	if len(msg.Name) != 0 && msg.Name == "command.PongMsg" {
		msg.Type = CMD_PONG
	}

	//心跳包，处理
	if msg.Type == CMD_PONG {

		//收到pong消息，重置发送ping数量
		log.Println(this.conn.RemoteAddr(), "收到pong包")
		this.pingCount = 0
		return
	}

	// 没有验证,则第一个包是验证包
	if !this.isVerify() {
		if this.Event != nil && this.Event.VerifyConn(msg) {

			log.Println("udptask验证成功")

			this.mutex.Lock()
			if this.verifyTimeout != nil {
				this.verifyTimeout.Stop()
			}
			this.mutex.Unlock()

			// 设置验证成功
			this.setVerify()
		} else {
			log.Println("udptask验证失败，关闭连接")
			this.Terminate()
		}
		return
	}

	if this.Event != nil {
		this.Event.MsgParse(msg)
	}
}

/*
func (this *UDPTask) sendCmd_NoBuf(buf []byte) (ret bool) {

	if this.IsTerminate() {
		return false
	}

	packetLen := len(buf)
	current := 0

	for current < packetLen {
		n, err := this.conn.Write(buf[current:packetLen])
		if err != nil {
			return false
		}

		current += n
	}

	return true
}
*/

func (this *UDPTask) SendCmd(cmd proto.Message) (ret bool) {

	name := proto.MessageName(cmd)
	data, _ := proto.Marshal(cmd)
	return this.SendCmd_NoPack(uint32(util.BKDRHash(name)), name, 0, data)
}

func (this *UDPTask) SendCmd_NoPack(typeid uint32, name string, index uint32, data []byte) bool {

	msg := new(command.Message)
	msg.Name = name
	msg.Type = typeid
	msg.Index = index
	msg.Data = data

	d := make([][]byte, 2)
	d[1], _ = proto.Marshal(msg)
	d[0] = util.Int2Byte(len(d[1]))

	f := []byte("")
	g := bytes.Join(d, f)

	if len(this.sndQueue) == cap(this.sndQueue)-1 {
		log.Println("close conn: channel full")
		this.Terminate()
		return true
	}

	this.mutex.Lock()

	if !this.terminate {

		//发送缓存
		this.sndQueue <- g
	}

	this.mutex.Unlock()

	return true
	//return this.sendCmd_NoBuf(g)
}

func (this *UDPTask) isVerify() bool {
	return this.verify
}

func (this *UDPTask) setVerify() {
	this.verify = true
}

func (this *UDPTask) IsTerminate() bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	return this.terminate
}

func (this *UDPTask) Terminate() {

	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.terminate {
		return
	}

	this.terminate = true
	this.sndQueue <- nil
}

func (this *UDPTask) ping() {

	log.Println(this.conn.RemoteAddr(), "发送ping", this.pingCount+1)
	msg := new(command.PingMsg)
	this.SendCmd(msg)

	this.pingCount++
}

func (this *UDPTask) GetConn() *kcp.UDPSession {
	return this.conn
}
