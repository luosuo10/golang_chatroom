package utils

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type User struct {
	id             int
	addr           string
	enterAt        time.Time
	msgCh chan Message
}


func (u *User) String() string {
	return fmt.Sprintf("%v, UID = %v, Enter At %v",
		u.addr, strconv.Itoa(u.id), u.enterAt.Format("2006-01-02 15:04:05"))
}

func NewUser(connAddr string, chBufSize int) *User {
	return &User{
		id: GenUserID(),
		addr: connAddr,
		enterAt: time.Now(),
		msgCh: make(chan Message, chBufSize),
	}
}

func (u *User) InfoSelfByIDAndContent(userID int, content string) {
	msg := NewMsg(userID, content)
	if msg.GetUserID() == u.id {
		return
	}
	u.msgCh <- msg
}

func (u *User) InfoSelfByMsg(msg Message) {
	//fmt.Printf("[InfoSelfByMsg], user = %+v, msg = %+v ", u, msg)
	if msg.GetUserID() == u.id {
		//fmt.Printf("skip\n")
		return
	}
	//fmt.Printf("\n")
	u.msgCh <- msg
}

func (u *User) SendMsg(conn net.Conn) {
	for msg := range u.msgCh {
		fmt.Fprintln(conn, msg)
	}
}

func (u *User) CloseMsgCh() {
	close(u.msgCh)
}

func (u *User) GetID() int {
	return u.id
}


func (u *User) GetIDStr() string {
	return strconv.Itoa(u.id)
}