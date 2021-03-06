package main

import (
	"net"
)

//创建用户结构体
type User struct {
	Name string
	Addr string
	C chan string
	conn net.Conn
	server *Server
}

//创建一个用户的API
func NewUser(conn net.Conn,server *Server) *User {

	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name : userAddr,
		Addr : userAddr,
		C : make(chan string),
		conn : conn,
		server : server,
	}

	//启动监听是否接受到消息
	go user.listenMessage()

	return user

} 

//给当前用户发送消息
func (this *User)sendMessage(msg string)  {
	this.conn.Write([]byte(msg))
}

func (this *User) Online(){
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	this.server.BroadCast(this,"上线")
}

func (this *User) Offline(){
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap,this.Name)
	this.server.mapLock.Unlock()
	this.server.BroadCast(this,"下线")
}

func (this *User) DoMessage(msg string)  {
	if msg == "who" {
		startMessage := "在线的用户如下：\n================\n" 
		this.sendMessage(startMessage)
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := user.Name+" : " + "在线\n"
			this.sendMessage(onlineMsg)
		}
		this.server.mapLock.Unlock()
		endMessage := "================\n" 
		this.sendMessage(endMessage)
	}else {
		this.server.BroadCast(this,msg)
	}
}

//监听当前chan的消息，一旦有消息就发送给user
func (this *User) listenMessage() {

	for {
		//msg := <-this.C
		this.conn.Write([]byte(<-this.C + "\n"))
	}

}