package main

import (
	"log"
	"net"
	"strconv"
)

type Server struct {
	Ip   string
	Port int
}

func newServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}

	return server
}

func (server *Server) Handler(conn net.Conn) {
	log.Print("连接建立成功！")
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
