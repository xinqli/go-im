package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return nil
	}
	client.conn = conn
	return client
}

var serverIp string
var serverPort int

// ./client -ip xxx
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址（默认127.0.0.1）")
	flag.IntVar(&serverPort, "port", 8888, "设置端口，默认8888")
}

func main() {
	// 命令解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>服务端连接失败>>>>>")
		return
	}
	fmt.Println(">>>>>服务端连接成功>>>>>")
	// 启动客户端的业务
	select {}
}
