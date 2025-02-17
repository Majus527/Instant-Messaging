package main

func main() {
	// Start server
	server := NewServer("127.0.0.1", 8081)
	server.Start()
}