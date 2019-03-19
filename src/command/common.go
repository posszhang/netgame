package command

import (
	"regexp"
	"strconv"
)

/*
 服务器的类型定义
*/
const (
	SuperServer   = 1
	LoginServer   = 2
	GatewayServer = 3
	SessionServer = 4
	RecordServer  = 5
	RouteServer   = 6
)

func ServerID2Index(id int) int {
	return int(id % 10000)
}

func ServerID2Type(id int) int {
	return int(id / 10000)
}

func GetServerID(tp int, index int) int {
	id := tp*10000 + index
	return id
}

func GetIndexFromFilename(filename string) int {
	reg := regexp.MustCompile(`\d`)
	ids := reg.FindAllString(filename, -1)

	if ids == nil || len(ids) == 0 {
		return 0
	}

	id, _ := strconv.Atoi(ids[0])

	return id
}
