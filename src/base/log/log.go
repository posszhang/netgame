package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Log struct {
	log   *logrus.Logger
	mutex sync.Mutex

	filename string
	lasttime int64
	pfile    *os.File
}

var log *Log = &Log{
	log: logrus.New(),
}

func NewLog(filename string) {

	if len(filename) == 0 {
		filename = fmt.Sprint("./log/", filepath.Base(os.Args[0]), ".log")
	}

	log = &Log{
		log:      logrus.New(),
		filename: filename,
	}

	log.init()
}

func (log *Log) init() {

	file := log.GetLogFile(log.filename)
	if file == nil {
		return
	}

	log.pfile = file
	log.lasttime = time.Now().Unix()

	log.log.Out = log.pfile

	go log.TimeAction()
}

func (log *Log) CheckHour() bool {

	now := time.Now().Hour()
	last := time.Unix(log.lasttime, 0).Hour()

	return now != last
}

func (log *Log) CheckMinute() bool {

	now := time.Now().Minute()
	last := time.Unix(log.lasttime, 0).Minute()

	return now != last
}

func (log *Log) SwapFile() {

	log.mutex.Lock()
	defer log.mutex.Unlock()

	lasttime := time.Unix(log.lasttime, 0)
	sfile := fmt.Sprintf("%s_%d%02d%d-%02d", log.filename, lasttime.Year(), lasttime.Month(), lasttime.Day(), lasttime.Hour())
	os.Rename(log.filename, sfile)

	file := log.GetLogFile(log.filename)
	if file == nil {
		return
	}

	log.log.Out = file

	log.lasttime = time.Now().Unix()
	log.pfile.Close()
	log.pfile = file
}

func (log *Log) TimeAction() {

	for {
		time.Sleep(1 * time.Second)

		if !log.CheckHour() {
			continue
		}

		log.SwapFile()
	}
}

func (log *Log) GetLogFile(filename string) *os.File {

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return nil
	}

	return file
}

func Debugln(args ...interface{}) {

	log.mutex.Lock()
	defer log.mutex.Unlock()

	log.log.Debugln(args...)
}

func Infoln(args ...interface{}) {
	log.mutex.Lock()
	defer log.mutex.Unlock()

	log.log.Infoln(args...)
}

func Println(args ...interface{}) {
	log.mutex.Lock()
	defer log.mutex.Unlock()

	log.log.Println(args...)
}

func Warnln(args ...interface{}) {
	log.mutex.Lock()
	defer log.mutex.Unlock()

	log.log.Warnln(args...)
}

func Warningln(args ...interface{}) {
	log.mutex.Lock()
	defer log.mutex.Unlock()

	log.log.Warningln(args...)
}

func Errorln(args ...interface{}) {
	log.mutex.Lock()
	defer log.mutex.Unlock()

	log.log.Errorln(args...)
}

func Fatalln(args ...interface{}) {
	log.mutex.Lock()
	defer log.mutex.Unlock()

	log.log.Fatalln(args...)
}

func Panicln(args ...interface{}) {
	log.mutex.Lock()
	defer log.mutex.Unlock()

	log.log.Panicln(args...)
}
