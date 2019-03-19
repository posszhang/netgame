package gnet

import (
	"base/log"
	"command"
	"fmt"
	"sync"
)

type MessageQueue struct {
	chMsg chan *command.Message

	mutex sync.Mutex
}

func (this *MessageQueue) Final() {

	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.chMsg != nil {
		close(this.chMsg)
	}
}

func (this *MessageQueue) Cache(msg *command.Message) {

	log.Println("message queue cache", msg)

	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.chMsg == nil {
		this.chMsg = make(chan *command.Message, 10240)
	}

	this.chMsg <- msg
}

func (this *MessageQueue) Do(fun func(*command.Message) bool) {

	this.mutex.Lock()
	defer this.mutex.Unlock()

	i := 0

	defer func() {
		if i > 100 {
			fmt.Println("本次执行消息数:", i, len(this.chMsg), cap(this.chMsg))
		}
	}()

	for {
		select {
		case msg := <-this.chMsg:
			if msg == nil {
				return
			}

			log.Println("message queue do", msg)
			i++
			fun(msg)
			continue
		default:
			return
		}
	}

	return
}