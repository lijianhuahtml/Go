package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	localAddr := &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 20001} // 设置本地IP地址和端口号
	dialer := &net.Dialer{LocalAddr: localAddr}                        // 创建Dialer结构体并设置LocalAddr字段
	conn, err := dialer.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	defer conn.Close()

	go listen(conn)

	inputReader := bufio.NewReader(os.Stdin)
	for {
		input, _ := inputReader.ReadString('\n') // 读取用户输入
		inputInfo := strings.Trim(input, "\r\n")
		if strings.ToUpper(inputInfo) == "Q" { // 如果输入q就退出
			return
		}
		_, err := conn.Write([]byte(inputInfo + "\n"))
		if err != nil {
			return
		}
		fmt.Println("发送给服务器信息：" + inputInfo)
	}
}

func listen(conn net.Conn) {
	fmt.Println("----------------")
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading message:", err.Error())
			return
		}
		fmt.Println("Received message from server:", message)
	}
}
