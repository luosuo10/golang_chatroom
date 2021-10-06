package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	var conn net.Conn
	var err error
	if conn, err = net.Dial("tcp", ":1234"); err != nil {
		panic(err)
	}

	// 新开的 goroutine 通过一个 channel 来和 main goroutine 通讯；
	done := make(chan struct{})
	go func() {
		// 通过 io.Copy 来操作 IO，包括从标准输入读取数据写入 TCP 连接中，以及从 TCP 连接中读取数据写入标准输出；
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors
		fmt.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()

	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		fmt.Println(err)
	}
}
