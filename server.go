package main

import (	
	"net"
	"fmt"
	"io"
	"log"
	"sync"
)

type Server struct {
	Ip	string
	Port int
	mapLock sync.RWMutex
	// 在线用户列表
	OnlineMap map[string]*User
	// 消息队列
	Message chan string
}

// 创建server
// NewServer creates a new Server instance with the provided IP address and port number. It initializes the OnlineMap for storing user connections and a Message channel for handling messages.
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip: ip,
		Port: port,
		OnlineMap: make(map[string]*User),
		Message: make(chan string),
	}
	return server
}

// 启动server
func (this *Server)Start() {
	// 创建listener
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
    log.Fatal(err)
		return
	}
	defer listener.Close()
	fmt.Println("server is running...")
	fmt.Printf("listening on %s:%d \n", this.Ip, this.Port)

	// 启动监听Message的goroutine
	go this.ListenMessager()

	// 监听
	for {
    	conn, err := listener.Accept()
    	if err != nil {
    	    fmt.Println(err)
    	    continue
    	}
    	go this.handleConn(conn)	// 使用goroutine并发处理连接
	}
}

// 监听消息队列
func (this *Server) ListenMessager() {
	for {
		msg := <- this.Message
		// 将msg发送给全部的在线User
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.conn.Write([]byte(msg))
		}
		this.mapLock.Unlock()
	}
}


// 广播消息
func (this *Server)Broadcast(user *User, msg string){
	sendMsg := user.Name + ":" + msg
	this.Message <- sendMsg
}


func (this *Server)handleConn(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)	// 创建一个缓冲区
	// 创建user
	user := NewUser(conn, this)

	// 将user添加到在线用户列表中
	this.OnlineMap[user.Name] = user

	// 用户上线
	user.Online()

	for {
		// 接收
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
					fmt.Println("Read error:", err)
			}
			break
		}
		// fmt.Printf("%s: %s",user.Name, string(buffer[:n]))
		msg := string(buffer[:n])

		// 发送
		// 广播消息
		this.Broadcast(user, msg)
		// conn.Write([]byte("----------Received your message------------\n"))
	}
}