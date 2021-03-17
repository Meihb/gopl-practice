package main

import (
	"fmt"
)

func main() {
	fmt.Println("基于共享变量的变法  竞争条件  不要使用共享数据来通信；使用通信来共享数据")

	/*
		第二种避免数据竞争的方法是，避免从多个goroutine访问变量。这也是前一章中大多数程序所采用的方法。例如前面的并发web爬虫(§8.6)的main goroutine
		是唯一一个能够访问seen map的goroutine，而聊天服务器(§8.10)中的broadcaster goroutine是唯一一个能够访问clients map的goroutine。这些变量都
		被限定在了一个单独的goroutine中。

		有意思啊,甚至semaphore何尝不是 第二种思路的延伸呢,当然缩短了该goroutine的占用时间,更像是协程了,效率理应更快
	*/
}
