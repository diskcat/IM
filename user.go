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
}

//创建一个用户的API
func NewUser(conn net.Conn) *User {

	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name : userAddr,
		Addr : userAddr,
		C : make(chan string),
		conn : conn,
	}

	//启动监听是否接受到消息
	go user.listenMessage()

	return user

} 

//监听当前chan的消息，一旦有消息就发送给user
func (this *User) listenMessage() {

	for {
		//msg := <-this.C
		this.conn.Write([]byte(<-this.C + "\n"))
	}

}