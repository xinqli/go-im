package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	// 在线用户mao
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息广播channel
	Message chan string
}

// 创建一个server的接口
func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

func (this *Server) handler(conn net.Conn) {
	// 当前连接的业务
	fmt.Println("连接创建成功")
	// 用户上线，将用户加入map中
	user := NewUser(conn)
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()
	// 广播当前用户上线消息
	this.BroadCast(user, "已上线")

	// 接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)

			if n == 0 {
				this.BroadCast(user, user.Name+"下线")
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("conn read err:", err)
				return
			}
			// 提取用户消息 去除\n
			msg := string(buf[:n-1])
			// 将msg广播
			this.BroadCast(user, msg)
		}
	}()

	// 当前阻塞
	select {}
}

// 监听message广播消息channel，一旦有消息就发给全部在线user
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message
		// 将msg发送给全部的user
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// 广播消息
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
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
	go this.ListenMessage()
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
