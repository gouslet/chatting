/*
 * File: /cmd/tcp/server.go                                                    *
 * Project: chatting                                                           *
 * Created At: Thursday, 2022/07/14 , 12:01:13                                 *
 * Author: elchn                                                               *
 * -----                                                                       *
 * Last Modified: Friday, 2022/07/15 , 13:14:01                                *
 * Modified By: elchn                                                          *
 * -----                                                                       *
 * HISTORY:                                                                    *
 * Date      	By	Comments                                                   *
 * ----------	---	---------------------------------------------------------  *
 */
package main

import (
	"bufio"
	"fmt"
	"gouslet/chatting/logic"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

var (
	enteringChannel = make(chan *logic.User) // 新用户到来，通过该channel登记
	leavingChannel  = make(chan *logic.User) // 用户离开，通过该channel登记
	messageChannel  = make(chan string)      // 用户离开，通过该channel登记
)

func main() {
	listener, err := net.Listen("tcp", ":2022")
	if err != nil {
		panic(err)
	}

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(conn)
	}
}

// broadcaster 用于记录聊天室用户，并进行消息广播
func broadcaster() {
	users := make(map[*logic.User]struct{})
	for {
		select {
		case user := <-enteringChannel:
			users[user] = struct{}{} // 新用户进入
		case user := <-leavingChannel:
			delete(users, user)        // 用户离开
			close(user.MessageChannel) // 关闭channel，避免goroutine泄漏
		case msg := <-messageChannel:
			// 给所有在线用户发送消息
			for user := range users {
				user.MessageChannel <- msg
			}
		}
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	// 1. 新用户进来，构建该用户的实例
	user := &logic.User{
		ID:             GenUserID(),
		Addr:           conn.RemoteAddr().String(),
		EnterAt:        time.Now(),
		MessageChannel: make(chan string, 8),
	}

	// 2. 由于当前是在一个新的goroutine中进行读操作的，所以需要开一个goroutine用于写操作
	go sendMessage(conn, user.MessageChannel)

	// 3. 给当前用户发送欢迎信息，向所有用户告知新用户的到来
	user.MessageChannel <- "Welcome, " + user.String()
	messageChannel <- "user: `" + strconv.Itoa(user.ID) + "` has entered"

	// 4. 记录到全局用户列表中，避免用锁
	enteringChannel <- user

	// 5. 循环读取用户输入
	input := bufio.NewScanner(conn)
	for input.Scan() {
		messageChannel <- strconv.Itoa(user.ID) + ": " + input.Text()
	}

	if err := input.Err(); err != nil {
		log.Println("读取错误：", err)
	}
	// 6. 用户离开
	leavingChannel <- user
	messageChannel <- "user: `" + strconv.Itoa(user.ID) + "` has left"

}

var (
	globalID int
	idLocker sync.Mutex
)

func GenUserID() int {
	idLocker.Lock()
	defer idLocker.Unlock()

	globalID++

	return globalID
}

func sendMessage(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}
