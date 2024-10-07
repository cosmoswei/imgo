package main

import (
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int
	// online user map
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	Message   chan string
}

func newServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

func (server *Server) Broadcast(user *User, msg string) {
	sengMsg := "[" + user.Addr + "]" + user.Name + ": " + msg
	server.Message <- sengMsg
}

// 监听message 广播
func (server *Server) listenMessage() {
	for {
		msg := <-server.Message
		server.mapLock.RLock()
		for _, cli := range server.OnlineMap {
			cli.C <- msg
		}
		server.mapLock.RUnlock()
	}
}

func (server *Server) Handler(conn net.Conn) {
	log.Print("[" + conn.RemoteAddr().String() + "] 连接建立成功！")

	// 将用户加入map
	user := newUser(conn, server)

	// 上线
	user.Online()

	isLive := make(chan bool)

	// 接受客户端的消息
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil && err != io.EOF {
				log.Println("conn.Read err", err)
				return
			}
			if n == 0 {
				user.Offline()
				return
			}
			msg := string(buf[:n-1])
			user.DoMessage(msg)
			isLive <- true
		}
	}()

	for {
		select {
		case <-isLive:
		case <-time.After(time.Second * 300):
			user.SendMsg("你被踢了！")
			user.Offline()
			close(user.C)
			conn.Close()
			return
		}
	}

	select {}
}

func (server *Server) Run() {
	// socket listen
	listen, err := net.Listen("tcp", server.Ip+":"+strconv.Itoa(server.Port))
	if err != nil {
		log.Print("net.Listen err ", err)
		return
	}

	log.Print("server run success")

	//close
	defer listen.Close()

	go server.listenMessage()

	for {
		// accept
		conn, err := listen.Accept()
		if err != nil {
			log.Print("listen.Accept err ", err)
			continue
		}
		// handler
		go server.Handler(conn)
	}

}
