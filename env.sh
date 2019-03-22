PROJECT_PATH=$(pwd)
export GOPATH=$PROJECT_PATH

go get github.com/golang/protobuf/proto
go get github.com/gorilla/websocket
go get github.com/xtaci/kcp-go
go get github.com/sirupsen/logrus

git clone https://github.com/golang/crypto.git $GOPATH/src/github.com/golang/crypto
git clone https://github.com/golang/net.git $GOPATH/src/github.com/golang/net
git clone https://github.com/golang/sys.git $GOPATH/src/github.com/golang/sys
git clone https://github.com/golang/text.git $GOPATH/src/github.com/golang/text
git clone https://github.com/golang/lint.git $GOPATH/src/github.com/golang/lint
git clone https://github.com/golang/tools.git $GOPATH/src/github.com/golang/tools

#因特殊原因，设置软链
mkdir $GOPATH/src/golang.org
ln -s $GOPATH/src/github.com/golang/ $GOPATH/src/golang.org/x
