package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net dial error", err)
		return nil
	}

	client.conn = conn
	return client
}

var serverIp string
var serverPort int

func init() {

	flag.StringVar(&serverIp, "serverIp", "127.0.0.1", "server ip")
	flag.IntVar(&serverPort, "serverPort", 8888, "server port")

}

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1. 公聊模式")
	fmt.Println("2. 私聊模式")
	fmt.Println("3. 更新用户名")
	fmt.Println("0. 已退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("请输入合法范围的数字")
		return false
	}
}

func (client *Client) UpdateName() bool {
	fmt.Println("请输入用户名：")
	fmt.Scanln(&client.Name) //这里如果不用地址接收会bug

	msg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("conn write error:", err)
		return false
	}

	return true
}

func (client *Client) PublicChat() {

	var chatMsg string

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn write error:", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println("请输入聊天内容，exit退出")
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) SelectUser() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn write error:", err)
		return
	}
}

func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string
	client.SelectUser()

	fmt.Println("请输入聊天对象用户名, exit退出")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println("输入聊天内容, exit 退出")
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn write error:", err)
					break
				}
			}
			chatMsg = ""
			fmt.Println("输入聊天内容, exit 退出")
			fmt.Scanln(&chatMsg)
		}

		client.SelectUser()
		fmt.Println("请输入聊天对象用户名, exit退出")
		fmt.Scanln(&remoteName)
	}

}

func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}

		switch client.flag {
		case 1:
			client.PublicChat()
			break
		case 2:
			client.PrivateChat()
			break
		case 3:
			client.UpdateName()
			break
		case 0:
			fmt.Println("退出")
			break

		}
	}

}

func main() {
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("服务器连接失败")
		return
	}

	// 单独异步处理server消息
	go client.DealResponse()

	fmt.Println("服务器连接成功")
	//为什么go 不能开在这里

	//阻塞
	client.Run()
}
