package main

import (
	"fmt"
	"net"
	"sync"
	"io"
)

type Server struct{
	Ip string
	Port int
	//添加用户在线列表
	OnlineMap map[string]*User
	mapLock sync.RWMutex
	//消息广播的channel
	Message chan string
}


//创建server对象
func NewServer(ip string, port int) *Server {
	return &Server{
		Ip : ip,
		Port : port,
		OnlineMap : make(map[string]*User),
		Message : make(chan string),
	}
}

//监听广播,创建为一个goroutine来模拟为服务
func (this *Server) listenMessage() {
	for {
		msg := <-this.Message
		//将onlineMap锁上
		this.mapLock.Lock()
		for _, user := range this.OnlineMap {
			user.C <- msg
		}
		this.mapLock.Unlock() 
	}
}

//Message接收信息,进行广播
func (this *Server)  BroadCast(user *User,msg string) {
	sendMsg := user.Name + " : " + msg
	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
    //将用户加入表中
	user := NewUser(conn)
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()
	//将消息发送到message中,进行广播
	this.BroadCast(user,"上线了")
	//接收客户端发送的消息
	go func(){
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				this.BroadCast(user, "下线")
				return
			}

			if err != nil && err != io.EOF {
				fmt.Printf("/n cannot read control input:", err)
				return
			}

			//接收传输的message
			msg := string(buf[:n-1])
			fmt.Printf(msg)
			this.BroadCast(user, msg)
		}
	}()
}

//启动服务
func (this *Server) start() {
	//socket listen
	listener, err := net.Listen("tcp",fmt.Sprintf("%s:%d",this.Ip,this.Port))
	if err != nil {
		fmt.Printf("net.Listen err:",err)
		return 
	}
	go this.listenMessage()
	//accpet
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("liseter.Accpet err:",err)
			continue
		}
		//do handler,异步非阻塞
		go this.Handler(conn)
	}

	defer listener.Close()
	//close listen socket
}
