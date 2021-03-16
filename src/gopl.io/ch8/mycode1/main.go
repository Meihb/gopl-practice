package main

import "fmt"

func main() {
	/*
		channel是iterable的数据类型,此外也是满足go 双参数返回值的数据类型,后者也是前者的实现前提
	*/
	naturals := make(chan int)
	squares := make(chan int)

	// Counter
	go func() {
		for x := 0; x < 5; x++ {
			naturals <- x
		}
		close(naturals)
	}()

	// Squarer
	go func() {
		for x := range naturals { //
			squares <- x * x
		}
		close(squares)
	}()

	// Printer (in main goroutine)
	for x := range squares {
		fmt.Println(x)
	}

	/*
		<-chan 或者chan<-表示单方向的channel,应该不具备 同步信号这个功能了吧
	*/
	

}
