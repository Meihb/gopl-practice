package main

import (
	"fmt"
	"time"
)

func main() {
	var timeout chan bool
	fmt.Printf("%T %[1]v %v\n", timeout, timeout == nil)
	// timeout = make(chan bool, 1)
	// fmt.Printf("%T %[1]v %v\n", timeout, timeout == nil)
	// ch2 := make(chan int)
	// close(timeout)//所以close之后的chan并不是nil,所谓 文中所说到底为何
	// fmt.Printf("%T %[1]v %v\n", timeout, timeout == nil)

	go func() {
		// time.Sleep(1 * time.Second) // sleep 3 seconds
		// timeout <- true
		// close(timeout)
		fmt.Printf("go :%T %[1]v %v\n", timeout, timeout == nil)
	}()
	time.Sleep(4 * time.Second)
	fmt.Printf("%T %[1]v %v\n", timeout, timeout == nil)
	select {
	// case <-ch2:
	case <-timeout: //如果这个通信channel 是nil,则此case 分支被无视
		fmt.Printf("timeout:%T %[1]v %v\n", timeout, timeout == nil)
		fmt.Println("timeout!")
	default:
		fmt.Printf("default :%T %[1]v %v\n", timeout, timeout == nil)
		fmt.Println("default case is running")
	}

}
