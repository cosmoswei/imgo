package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int, name string) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	dial, err := net.Dial("tcp", serverIp+":"+strconv.Itoa(serverPort))
	if err != nil {
		log.Fatal(err)
		return nil
	}
	client.conn = dial
	return client
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "server ip address")
	flag.IntVar(&serverPort, "port", 8088, "server port")
}

func main() {
	flag.Parse()
	log.Print("server ip:", serverIp, " port:", serverPort)
	client := NewClient(serverIp, serverPort, "")
	if nil == client {
		log.Print("连接服务器失败")
		return
	}
	go client.DearResponse()

	log.Print("连接服务器成功")

	client.Run()
}

func (client *Client) menu() bool {
	var flags int
	fmt.Println("1. 公聊模式")
	fmt.Println("2. 私聊模式")
	fmt.Println("3. 更新用户名")
	fmt.Println("4. 用户查询")
	fmt.Println("0. 退出")
	fmt.Scanln(&flags)
	if flags >= 0 && flags <= 4 {
		client.flag = flags
		return true
	} else {
		fmt.Println("输入不合法")
	}
	return false
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}
		switch client.flag {
		case 1:
			fmt.Println("公聊模式")
			client.PublicChat()
			break
		case 2:
			fmt.Println("私聊模式")
			client.PrivateChat()
			break
		case 3:
			fmt.Println("更新用户名")
			client.updateName()
			break
		case 4:
			fmt.Println("用户查询")
			client.PrintOnlineUsers()
			break
		}
	}
}

func (client *Client) updateName() bool {
	fmt.Println("请输入用户名：")
	fmt.Scanln(&client.Name)
	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func (client *Client) PublicChat() {
	var chatMsg string

	fmt.Println("请输入聊天内容，exit 退出")
	fmt.Scanln(&chatMsg)
	for "exit" != chatMsg {
		// 发给服务器
		if len(chatMsg) > 0 {
			senfMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(senfMsg))
			if err != nil {
				log.Fatal(err)

			}
		}
		chatMsg = ""
		fmt.Println("请输入聊天内容，exit 退出")
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) PrintOnlineUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		log.Fatal(err)
		return
	}
}

func (client *Client) PrivateChat() {
	privateUser := ""
	// 查询用户
	client.PrintOnlineUsers()
	// 选择用户
	fmt.Println("请输入选择的聊天对象[用户名]，exit 退出")
	fmt.Scanln(&privateUser)
	// 发送消息
	chatMsg := ""
	for "exit" != privateUser {
		fmt.Println("请输入的消息内容，exit 退出")
		fmt.Scanln(&chatMsg)
		for "exit" != chatMsg {
			if len(chatMsg) > 0 {
				senfMsg := "to|" + privateUser + "|" + chatMsg + "\n"
				_, err := client.conn.Write([]byte(senfMsg))
				if err != nil {
					log.Fatal(err)
					break
				}
			}
			chatMsg = ""
			fmt.Println("请输入的消息内容，exit 退出")
		}
	}
	client.PrintOnlineUsers()
	// 选择用户
	fmt.Println("请输入选择的聊天对象[用户名]，exit 退出")
	fmt.Scanln(&privateUser)
}

func (client *Client) DearResponse() {
	_, err := io.Copy(os.Stdout, client.conn)
	if err != nil {
		return
	}
}
