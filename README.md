# 网络游戏服务器脚手架 netgame
多进程分布式游戏服务器框架

## 开始


```bash
#1
git clone https://github.com/posszhang/netgame.git
#2
cd netgame
#3
./env.sh
```

执行`env.sh`脚本，将会自动获取依赖库，以及设置GOPATH

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

<br/>

# 执行更新...
