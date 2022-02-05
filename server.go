package main

import (
	"fmt"
	"net"
	"sync"
	"io"
	"time"
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
	isLive := make(chan bool)
    //将用户加入表中
	user := NewUser(conn,this)
	user.Online()
	//接收客户端发送的消息
	go func(){
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Printf("/n cannot read control input:", err)
				return
			}

			//接收传输的message
			msg := string(buf[:n-1])
			user.DoMessage(msg)
			isLive <- true
		}
	}()
	//handler 阻塞
	for {
		select{
			//当前用户是活跃的
			case <-isLive:
				
			case <-time.After(time.Second * 300):
				user.sendMessage("你被踢了")

				close(user.C)
				conn.Close()
				return 
		}
	}

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
