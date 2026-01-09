package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// 构造函数
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}

	return server
}

// 回调接口，当前链接的业务
func (this *Server) Handler(conn net.Conn) {
	fmt.Println("链接建立成功")

}

// 服务器启动的接口
func (this *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen error: ", err)
		return
	}

	defer listener.Close()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept error: ", err)
			continue
		}

		//do handler
		go this.Handler(conn)
	}

	//close listen socket

}
