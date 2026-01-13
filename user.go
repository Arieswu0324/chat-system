package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

func (this *User) Online() {
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	//广播当前用户上线消息
	this.server.BroadCast(this, "已上线")
}

func (this *User) Offline() {
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()
	this.server.BroadCast(this, "下线")
}

func (this *User) SendMessage(msg string) {
	this.conn.Write([]byte(msg))
}
func (this *User) DoMessage(msg string) {
	if msg == "who" { //定义通讯规则，如果用户输入who，则表示查询在线用户
		this.server.mapLock.Lock()
		for _, usr := range this.server.OnlineMap {
			onlineMsg := "[" + usr.Addr + "]" + usr.Name + "在线...\n"
			//this.C <- onlineMsg 为什么不能这样
			this.SendMessage(onlineMsg)
		}

		this.server.mapLock.Unlock()

	} else if len(msg) > 7 && msg[:7] == "rename|" { //定义通信协议，如果用户以rename|XXX这种格式输入，则表示要修改用户名
		newName := strings.Split(msg, "|")[1]
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMessage("当前用户名被占用\n")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()

			this.Name = newName
			this.SendMessage("更新用户名成功：" + this.Name + "\n")
		}

	} else {
		this.server.BroadCast(this, msg)
	}
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	//启动监听当前user channel的goroutine
	go user.ListenMessage()
	return user
}

func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
