package main

import (
	"net"
	"strings"
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
		for name, _ := range this.server.OnlineMap {
			onlineMsg := name+" : " + "在线\n"
			this.sendMessage(onlineMsg)
		}
		this.server.mapLock.Unlock()
		endMessage := "================\n" 
		this.sendMessage(endMessage)
	}else if len(msg) > 7 && msg[:7]=="rename|"{	
		userName := strings.Split(msg,"|")[1]
		if _, ok := this.server.OnlineMap[userName]; ok {
			this.sendMessage("用户名已经存在\n")
		}else{
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[userName] = this
			this.server.mapLock.Unlock()
			this.Name = userName
			this.sendMessage("用户名更改成功\n")
		}
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