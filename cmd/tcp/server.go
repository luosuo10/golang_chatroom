package main

import (
	"bufio"
	"fmt"
	"golang_chatroom/cmd/tcp/utils"
	"net"
)


func main() {
	var listener net.Listener
	var err error
	if listener, err = net.Listen("tcp", ":1234"/*没有指定IP时 默认绑定到所有该主机的IP上*/); err != nil {
		panic(err)
	}
	go charRoomWorker()

	for {
		var conn net.Conn
		if conn, err = listener.Accept(); err != nil {
			fmt.Println(err)
			continue
		}
		go handleConn(conn)
	}
}

var (
	userInCh = make(chan *utils.User)
	userOutCh = make(chan *utils.User)
	MsgBufSize = 8
	globalMsgCh = make(chan utils.Message, MsgBufSize)
)


func charRoomWorker() {
	users := make(map[*utils.User]struct{})
	for {
		select {
		case user := <-userInCh:
			users[user] = struct{}{}
			fmt.Printf("add user %+v\n", user)
		case user := <-userOutCh:
			delete(users, user)
			fmt.Printf("delete user %+v\n", user)
		case msg := <-globalMsgCh:
			// 给所有在线用户发送消息
			for user := range users {
				user.InfoSelfByMsg(msg)
			}
		}
	}
}


func handleConn(conn net.Conn) {
	defer func() {
		fmt.Printf("conn = %v close\n", conn)
		conn.Close()
	}()
	utils.CreateIDGenerator()
	user := addInUser(conn)
	listenUserMsg(conn, user)
	delOutUser(user)  /*delete pointer not value*/
}

func addInUser(conn net.Conn) (user *utils.User) {
	user = utils.NewUser(conn.RemoteAddr().String(), MsgBufSize)
	userInCh <- user
	infoOthers(user.GetID(), fmt.Sprintf("user[%v] has enter", user.GetIDStr()))

	go user.SendMsg(conn)
	user.InfoSelfByIDAndContent(user.GetID(), fmt.Sprintf("Welcome, user[%v]", user.String()))
	return user
}

func delOutUser(user *utils.User) {
	fmt.Printf("[delOutUser] user = %p\n", user)
	user.CloseMsgCh()
	userOutCh <- user
	infoOthers(user.GetID(), fmt.Sprintf("user[%v] has left", user.GetIDStr()))
}

func listenUserMsg(conn net.Conn, user *utils.User) {
	input := bufio.NewScanner(conn)
	for input.Scan() {
		infoOthers(user.GetID(), fmt.Sprintf("user[%v] : %v", user.GetIDStr(), input.Text()))
	}
	if err := input.Err(); err != nil {
		fmt.Println("read err = ", err)
	}
}

func infoOthers(userID int, content string) {
	msg := utils.NewMsg(userID, content)
	globalMsgCh <- msg
}
