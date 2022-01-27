package main

import (
	"fmt"
	"net"
	"sync"
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
	msg := <-this.Message
	for {
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
		//do handler
		go this.Handler(conn)
		//启动监听Message的goroutine
	}

	defer listener.Close()
	//close listen socket
}
