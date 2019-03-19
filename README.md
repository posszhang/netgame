# 网络游戏服务器脚手架 netgame
## 多进程分布式游戏服务器框架

使用本项目务必正确设置GOPATH
```bash
export GOPATH=/home/posszhang/gamenet
```

<br/>

然后执行`env.sh`脚本，里面包含项目中依赖库和各种其他相关的设置
```Bash
go get github.com/golang/protobuf/proto
go get github.com/gorilla/websocket
go get github.com/xtaci/kcp-go
go get github.com/sirupsen/logrus

#因特殊原因，走github下载
git clone https://github.com/golang/crypto.git $GOPATH/src/github.com/golang/crypto
git clone https://github.com/golang/net.git $GOPATH/src/github.com/golang/net
git clone https://github.com/golang/sys.git $GOPATH/src/github.com/golang/sys
git clone https://github.com/golang/text.git $GOPATH/src/github.com/golang/text
git clone https://github.com/golang/lint.git $GOPATH/src/github.com/golang/lint
git clone https://github.com/golang/tools.git $GOPATH/src/github.com/golang/tools

mkdir $GOPATH/src/golang.org
#因特殊原因，设置软链
ln -s $GOPATH/src/github.com/golang/ $GOPATH/src/golang.org/x
```
<br/>

进入src目录执行`make`即可编译，当然你也可以根据自己的实际情况修改Makefile文件

<br/>

## 服务器架构拆解
* superserver `服务器管理器，所有的服务器均要连接次服务`
* routeserver `路由服务器，主要做消息路由，解耦服务器关系` 
* loginserver `登陆服务器，多点，用来做渠道sdk登陆等其他功能` 
* gatewayserver `网络服务器，多点，代理用户之后所有的数据通信` 
* sessionserver `会话服务器，暂单点，可以做成多点，目前没这必要，尽量少逻辑，保存用户状态` 
* logicserver `逻辑服务器，rpg中又称为sceneserver, 场景服务器` 
* recordserver `数据库服务器，可多点，滚服游单点即可，oneworld可设计多点，0用来作为数据库结点分配器，1-N实际用来处理用户数据`

##执行更新...
