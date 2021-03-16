// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 227.

// Netcat is a simple read/write client for TCP servers.
package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var done chan struct{} = make(chan struct{})

//!+
func main() {
	conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		fmt.Printf("%T %[1]v \n", conn)
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors 这是个阻塞函数呀,要等conn 输出结束才return
		fmt.Printf("%T %[1]v \n", os.Stdout)
		log.Println("done")
		done <- struct{}{} // signal the main goroutine 示意关闭main goroutine
	}()
	fmt.Println("finish goroutine")
	mustCopy(conn, os.Stdin) //stdin是标准输入,不是你安逸航空个回车啥的就会结束,而是按照EOF这个标识判断,而terminal可以使用C^z模拟此信号
	conn.Close()
	fmt.Println("Connection closes") //我的天,ctrl z就是结束标准输入啊,你好蠢
	<-done                           // wait for background goroutine to finish 整个fo func和这上面的一系列main 代码各自为政,但是遇到unbuffed chan时互相阻塞,等到两边都准备好才会结束,这就是 同步信号量吧


}

//!-

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
