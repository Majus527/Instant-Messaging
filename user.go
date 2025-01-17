package main

import(
	"net"
	"io"
	"fmt"
)

type User struct {
	Name string
	Addr string
	C chan string
	conn net.Conn
}


func NewUser(conn net.Conn) *User {
	// 获取远程地址
	remoteAddr := conn.RemoteAddr().String()

	// 让用户输入用户名
	buffer := make([]byte, 1024)
	conn.Write([]byte("input your name: "))
	n, err := conn.Read(buffer)
	if err != nil {
			if err != io.EOF {
					fmt.Println("Read error:", err)
			}
	}
	name := string(buffer[:n])
	name = name[:len(name)-1]

	// 初始化用户
	user := &User {
		Name: name,
		Addr: remoteAddr,
		conn: conn,
	}

	return user
}

// 用户上线业务
func (this *User) Online() {
	
}

//用户下线业务
func (this *User) Offline() {
	
}

// 用户处理消息的业务
func (this *User) DoMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}

// 监听当前user channel的方法，一旦有消息，直接发送给对端客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}