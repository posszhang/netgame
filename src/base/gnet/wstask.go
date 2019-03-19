package gnet

import (
	"bytes"
	"command"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"net"
	//	"runtime/debug"
	"base/log"
	"base/util"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type WSTaskEvent interface {
	VerifyConn(msg *command.Message) bool
	RecycleConn() bool
	MsgParse(msg *command.Message) bool
}

type WSTask struct {
	Event WSTaskEvent
	conn  *websocket.Conn

	BufferSnd bool
	sndQueue  chan []byte

	terminate bool
	mutex     sync.Mutex

	//超时验证
	verify        bool
	verifyTimeout *time.Timer
}

func (this *WSTask) GoHandler(conn *websocket.Conn) {

	this.conn = conn
	this.sndQueue = make(chan []byte, 16)

	// 协程读
	go this.reader()
	// 协程写
	go this.writer()

	this.verifyTimeout = time.AfterFunc(time.Second*time.Duration(10), func() {
		if this.IsTerminate() {
			return
		}
		this.Terminate()
		this.verifyTimeout = nil
		log.Println("task验证超时,关闭连接")
	})
}

func (this *WSTask) reader() {

	defer func() {
		log.Println("WSTask render 退出")
	}()

	bytesBuffer := bytes.NewBuffer([]byte{})
	header := make([]byte, 4)
	buf := make([]byte, 0)

	size := 0
	readsize := 0

	var timeout time.Duration = 60

	for {

		//60秒读超时，发送心跳
		var err error

		this.conn.SetReadDeadline(time.Now().Add(timeout * time.Second))

		//读取数据
		if _, buf, err = this.conn.ReadMessage(); err == nil {
			binary.Write(bytesBuffer, binary.BigEndian, buf)

			//拆包
			for {

				//读取包大小
				if size == 0 {
					if bytesBuffer.Len() < 4 {
						break
					}

					readsize, err = bytesBuffer.Read(header)
					if readsize != 4 || err != nil {
						break
					}

					size = util.Byte2Int(header)
					log.Println("包头:", header, ",", size, ",", bytesBuffer.Len())
				}

				if bytesBuffer.Len() < size {
					break
				}

				cmd := make([]byte, size)
				readsize, err = bytesBuffer.Read(cmd)
				if readsize != size || err != nil {
					break
				}
				size = 0
				this.doCmd(cmd)
			}

			continue
		}

		//是超时错误,检查心跳
		//websocket的超时将会导致链接直接失效，所以要结束
		if e, ok := err.(net.Error); ok && e.Timeout() {
			log.Println("读取数据超时", err.Error())
		}

		// 异常结束
		this.Terminate()
		break
	}

	if this.sndQueue != nil {
		close(this.sndQueue)
	}

	//验证成功的需要通知释放
	if this.isVerify() {
		this.Event.RecycleConn()
	}

}

func (this *WSTask) writer() {

	for {

		//如果关闭了，可以读不能写
		buf := <-this.sndQueue

		if buf == nil {
			break
		}

		this.SendCmd_NoBuf(buf)
	}

}

func (this *WSTask) doCmd(buf []byte) {

	msg := new(command.Message)
	err := proto.Unmarshal(buf, msg)
	if err != nil {
		log.Println("解包失败:", string(buf), err)
		return
	}

	//心跳包，处理
	if msg.Type == 0 {
		return
	}

	// 没有验证,则第一个包是验证包
	if !this.isVerify() {
		if this.Event != nil && this.Event.VerifyConn(msg) {

			if this.verifyTimeout != nil {
				this.verifyTimeout.Stop()
			}

			// 设置验证成功
			this.setVerify()
		} else {
			log.Println("tcptask验证失败，关闭连接")
			this.Terminate()
		}
		return
	}

	if this.Event != nil {
		this.Event.MsgParse(msg)
	}
}

func (this *WSTask) SendCmd_NoBuf(buf []byte) (ret bool) {

	if this.IsTerminate() {
		return false
	}

	err := this.conn.WriteMessage(websocket.BinaryMessage, buf)
	if err != nil {
		this.Terminate()
		return false
	}

	return true
}

func (this *WSTask) SendCmd_NoPack(typeid uint32, name string, data []byte) bool {

	if this.IsTerminate() {
		return false
	}

	msg := new(command.Message)
	msg.Name = name
	msg.Type = typeid
	msg.Data = data

	d := make([][]byte, 2)
	d[1], _ = proto.Marshal(msg)
	d[0] = util.Int2Byte(len(d[1]))

	f := []byte("")
	g := bytes.Join(d, f)

	if this.BufferSnd {

		this.mutex.Lock()
		defer this.mutex.Unlock()

		if !this.terminate {
			this.sndQueue <- g
		}
		return true
	}

	return this.SendCmd_NoBuf(g)
}

func (this *WSTask) SendCmd(cmd proto.Message) (ret bool) {

	if this.IsTerminate() {
		return false
	}

	msg := new(command.Message)
	msg.Name = proto.MessageName(cmd)
	msg.Type = util.BKDRHash(msg.Name)
	msg.Data, _ = proto.Marshal(cmd)

	d := make([][]byte, 2)
	d[1], _ = proto.Marshal(msg)
	d[0] = util.Int2Byte(len(d[1]))

	f := []byte("")
	g := bytes.Join(d, f)

	if this.BufferSnd {

		this.mutex.Lock()
		defer this.mutex.Unlock()

		if !this.terminate {
			this.sndQueue <- g
		}
		return
	}

	this.SendCmd_NoBuf(g)

	return true
}

func (this *WSTask) isVerify() bool {
	return this.verify
}

func (this *WSTask) setVerify() {
	this.verify = true
}

func (this *WSTask) IsTerminate() bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	return this.terminate
}

func (this *WSTask) Terminate() {

	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.terminate {
		return
	}

	this.conn.Close()
	this.terminate = true
}

func (this *WSTask) sendHeartBeat() {
	msg := new(command.Message)
	msg.Type = uint32(0)
	msg.Data = make([]byte, 0)
	this.SendCmd(msg)
}
