package main

import (
	"log"
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

func newUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	go user.listenMsg()
	return user
}

func (u *User) listenMsg() {
	for {
		msg := <-u.C
		u.conn.Write([]byte(msg + "\n"))
	}
}

func (u *User) Online() {
	// 将用户加入map
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()
	// 广播
	u.server.Broadcast(u, "已上线")
}

func (u *User) Offline() {
	// 将用户加入map
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()
	// 广播
	u.server.Broadcast(u, "已下线")
	log.Print("连接断开成功！")
}

func (u *User) SendMsg(msg string) {
	u.conn.Write([]byte(msg))
}

func (u *User) DoMessage(msg string) {

	// 查询当前用户
	if "who" == msg {
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ": " + "在线...\n"
			u.SendMsg(onlineMsg)
		}
		u.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]
		_, _ok := u.server.OnlineMap[newName]
		if _ok {
			u.SendMsg("当前名称已被占用\n")
		} else {
			u.server.mapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[newName] = u
			u.server.mapLock.Unlock()
			u.Name = newName
			u.SendMsg("你已经更新新的用户名：" + u.Name + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		// 消息格式 to|张三｜你好
		// 获取用户名
		remoteName := strings.Split(msg, "|")[1]
		if "" == remoteName {
			u.SendMsg("消息格式不正确，{to|张三｜你好}")
			return
		}

		// 得到对方user对象
		user, _ok := u.server.OnlineMap[remoteName]
		if !_ok {
			u.SendMsg("该用户名不存在")
			return
		}
		// 根据消息内容，发送消息
		content := strings.Split(msg, "|")[2]
		user.SendMsg("[" + u.Name + "]对您说：[" + content + "]\n")
	} else {
		u.server.Broadcast(u, msg)
	}
}
