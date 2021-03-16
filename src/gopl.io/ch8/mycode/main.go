package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

var finish chan bool = make(chan bool)

func main() {
	fmt.Println("Chapter 8:Goroutines and Channels")
	/*
		communicating sequential processes (CSP)顺序通信进程
	*/

	//所有的goroutine在 main goroutine return 之后统一结束

	/*
		ch := make(chan int) // ch has type 'chan int'
		和map类似，channel也对应一个make创建的底层数据结构的引用。当我们复制一个channel或用于函数参数传递时，我们只是拷贝了一个channel引用，因此调用者和被调用者将引用同一个channel对象。
		和其它的引用类型一样，channel的零值也是nil。
		两个相同类型的channel可以使用==运算符比较。如果两个channel引用的是相同的对象，那么比较的结果为真。一个channel也可以和nil进行比较。
	*/

	listener, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		log.Fatal("err")
	}

	for {
		fmt.Println("initiate conn")
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn)
		true := <-finish //等待finish命令
		if true {
			io.Copy(os.Stdout, conn)
		}

	}
}
func handleConn(c net.Conn) {
	defer c.Close()
	for i := 0; i < 3; {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second / 10)
		i++
	}
	finish <- true

}
