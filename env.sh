#设置GOPATH目录为当前目录，并写入.bashrc文件
PROJECT_PATH=$(pwd)
echo $GOPATH
echo export GOPATH=$PROJECT_PATH >> ~/.bashrc
source ~/.bashrc

git clone https://github.com/golang/crypto.git $GOPATH/src/github.com/golang/crypto
git clone https://github.com/golang/net.git $GOPATH/src/github.com/golang/net
git clone https://github.com/golang/sys.git $GOPATH/src/github.com/golang/sys

#因特殊原因，设置软链
mkdir $GOPATH/src/golang.org
ln -s $GOPATH/src/github.com/golang/ $GOPATH/src/golang.org/x

if [ ! -d "./src/github.com/golang/protobuf/proto" ]; then
		echo install protobuf
		go get github.com/golang/protobuf/proto
fi

if [ ! -d "./src/github.com/gorilla/websocket" ]; then
		echo install websocket
		go get github.com/gorilla/websocket
fi


if [ ! -d "./src/github.com/xtaci/kcp-go" ]; then
		echo install kcp-go
		go get github.com/xtaci/kcp-go
fi


if [ ! -d "./src/github.com/sirupsen/logrus" ]; then
		echo install logrus
		go get github.com/sirupsen/logrus
fi

exit

