package util

import (
	"errors"
	"net"

	"go.uber.org/zap"
)

func NewTcpAddr() *net.TCPAddr {
	//把字符串解析为tcp端点
	addr, _ := net.ResolveTCPAddr("tcp", "192.168.199.128:0")
	//tcp协议中,如果端口为0则在listen,dial这些函数中会默认给它分配一个空闲的端口
	lis, _ := net.ListenTCP("tcp", addr)
	defer lis.Close()
	addr, ok := lis.Addr().(*net.TCPAddr)
	if ok {
		return addr
	}
	err := errors.New("内部服务器错误,无法产生空闲的端口")
	zap.S().Errorw(err.Error())
	panic(err)
}
