package main

func main() {
	server := NewServer("192.168.104.43", 8081)
	server.Start()
}