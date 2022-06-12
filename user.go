package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

// 创建一个user对象
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	// 启动监听当前user channel 消息的goroutine
	go user.ListenMessage()
	return user
}

// 监听当前user channel的方法，一旦有消息，就直接发给客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}

func (this *User) OnLine() {
	// 用户上线，将用户加入map中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	// 广播当前用户上线消息
	this.server.BroadCast(this, "已上线")
}
func (this *User) OffLine() {
	// 用户下线，将用户移除map
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()
	// 广播当前用户上线消息
	this.server.BroadCast(this, "已下线")

}
func (this *User) DoMessage(msg string) {

	// 查询当前用户都有哪些
	if msg == "who" {
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ": 在线...\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
		return
	}

	this.server.BroadCast(this, msg)
}

// 给当前用户发送消息
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}
