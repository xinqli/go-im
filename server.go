package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// 创建一个server的接口
func NewServer(ip string, port int) *Server {
	return &Server{Ip: ip, Port: port}
}

func (this *Server) handler(conn net.Conn) {
	// 当前连接的业务
	fmt.Println("连接创建成功")
}

// 启动服务器接口
func (this *Server) Start() {

	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen error: ", err)
	}
	// close listen socket
	defer listener.Close()
	for {
		// accpet
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept error: ", err)
			continue
		}
		// do handler
		go this.handler(conn)
	}

}
