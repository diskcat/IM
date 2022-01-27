package main

import (
	"fmt"
	"net"
)

type Server struct{
	Ip string
	Port int
}

//创建server对象
func NewServer(ip string, port int) *Server {
	return &Server{
		Ip : ip,
		Port : port,
	}
}

func (this *Server) Handler(conn net.Conn) {
    fmt.Printf("连接建立成功\n")
}

//启动服务
func (this *Server) start() {
	//socket listen
	listener, err := net.Listen("tcp",fmt.Sprintf("%s:%d",this.Ip,this.Port))
	if err != nil {
		fmt.Printf("net.Listen err:",err)
		return 
	}
	//accpet
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("liseter.Accpet err:",err)
			continue
		}
		//do handler
		go this.Handler(conn)
	}

	defer listener.Close()
	//close listen socket
}
