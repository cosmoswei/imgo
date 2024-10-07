package main

func main() {
	server := newServer("127.0.0.1", 8088)
	server.Run()
}
