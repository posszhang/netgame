package main

import (
	"fmt"
)

type RobotManager struct {
	robotList []*Robot
}

func NewRobotManager() *RobotManager {
	mgr := &RobotManager{}

	return mgr
}

func (mgr *RobotManager) Do() {

	for _, robot := range mgr.robotList {
		robot.Do()
	}
}

func (mgr *RobotManager) Add(num int, ip string, port int) {
	for i := 0; i != num; i++ {
		robot := NewRobot(ip, port, fmt.Sprint(i), "123456")
		robot.Run()

		mgr.robotList = append(mgr.robotList, robot)
	}
}
