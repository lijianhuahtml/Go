package main

import (
	"fmt"
	"net"
)

//简单来说，一个TCP服务监听某个端口，但可以同时处理多个客户端与该端口的连接。
//这样，多个用户就可以连接到同一个端口上的服务，服务可以依次处理它们的请求。

type client chan<- string // 一个只写的channel，用于发送消息给客户端

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // 所有传入消息都发送到此channel
)

func broadcaster() {
	clients := make(map[client]bool) // 所有连接的客户端
	for {
		select {
		case msg := <-messages:
			fmt.Println("---message:" + msg)
			// 向所有客户端广播消息
			for cli := range clients {
				cli <- msg
			}

		case cli := <-entering:
			clients[cli] = true // 新客户端连接

		case cli := <-leaving:
			delete(clients, cli) // 客户端离开
			close(cli)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string) // 对每个连接都有一个传出消息的channel
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- ch

	for {
		buf := make([]byte, 512)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			// 可以在这里做一些处理，比如关闭连接或者返回
			return
		}
		fmt.Println("Received message from client:", who+" ", string(buf[:n]))
		//messages <- who + ": " + string(buf[:n])
	}

	//leaving <- ch
	//messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // 将传入的消息写入客户端连接
		break
	}
	fmt.Println("-----send----")
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8000") // 创建监听localhost:8000 端口
	if err != nil {
		fmt.Println("Failed to create listener:", err)
		return
	}
	defer listener.Close()

	go broadcaster() // 启动广播器

	for {
		//conn是一个net.Conn类型的对象，代表了与客户端建立的TCP连接。
		//通过这个连接，你可以在服务器端接收来自客户端的消息，并且也可以向客户端发送消息。
		//具体地说，net.Conn接口提供了Read和Write方法，用于从连接中读取数据和向连接中写入数据。
		//因此，你可以使用conn对象来接收和发送消息。
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}

		fmt.Println("Accepted connection from:", conn.RemoteAddr())

		go handleConn(conn) // 处理连接
	}
}
