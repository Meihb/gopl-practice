package main

import (
	"fmt"
	"image"
	"sync"
	"time"
)

func main() {
	fmt.Println("基于共享变量的变法  竞争条件(race condition)  不要使用共享数据来通信；使用通信来共享数据")

	/*
		第二种避免数据竞争的方法是，避免从多个goroutine访问变量。这也是前一章中大多数程序所采用的方法。例如前面的并发web爬虫(§8.6)的main goroutine
		是唯一一个能够访问seen map的goroutine，而聊天服务器(§8.10)中的broadcaster goroutine是唯一一个能够访问clients map的goroutine。这些变量都
		被限定在了一个单独的goroutine中。

		有意思啊,甚至semaphore何尝不是 第二种思路的延伸呢,当然缩短了该goroutine的占用时间,更像是协程了,效率理应更快

		文中三种防止数据竞争的办法都是有法可循的,第一种避免 发生数据竞争
		1. 一开始就初始化且不会发生竞争读写 lazy initialization
		2. 读写都在同一个goroutine中实现,不要再共享数据中通信,在通信中共享数据
		3. 第二种方法是在太占用资源了,没有达到并发的最大利用效率,故派生了读写锁这类,只在锁时符合第二条,读写在本goroutine中
		希望 能看到理论在方法中主键延伸的步骤
	*/

	/*
		var x, y int
		go func() {
			x = 1                   // A1
			fmt.Print("y:", y, " ") // A2
		}()
		go func() {
			y = 1                   // B1
			fmt.Print("x:", x, " ") // B2
		}()

		以上两个go routine会有什么结果打印出来呢
		A1,A2,B1,B2 y:0 x:1
		A1,B1,A2,B2 y:1,x:1
		A1,B1,B2,A2 x:1 y:1
		B1,B2,A1,A2 x:0 y:1
		B1,A1,A2,B1 y:1 x:1
		B1,A1,B2,A2 x:1 y:1

		其实这些都有一个问题,为什么你一定要线性思维考虑问题,咋就不能是A1/B1/B2 A2 这样呢,结果就是  y:0 x:0或者类似的x:0 y:0
		所以当程序不再在同一个goroutine中之时,永远要跳出线性思维来考虑,只有在同一个goroutine的代码才会严格执行线性顺序

		matters:所有并发的问题都可以用一致的、简单的既定的模式来规避。所以可能的话，将变量限定在goroutine内部；如果是多个goroutine都需要访问的变量，使用互斥条件来访问。

	*/

	/*
		var mu sync.RWMutex // guards icons
		var icons map[string]image.Image
		// Concurrency-safe.
		func Icon(name string) image.Image {
		    mu.RLock()
		    if icons != nil {
		        icon := icons[name]
		        mu.RUnlock()
		        return icon
		    }
		    mu.RUnlock()

		    // acquire an exclusive lock
		    mu.Lock()
		    if icons == nil { // NOTE: must recheck for nil
		        loadIcons()
		    }
		    icon := icons[name]
		    mu.Unlock()
		    return icon
		}

		这一段代码只有这个 must recheck for nil 最值得关注
	*/

	/*
		测试一下,channel 在关闭之后,其读写操作会迅速结束阻塞 这个情况是否属实

		所以channel 通道实现一对一通信非常方便,而channel本身更能实现一对多的notify,妙啊
		对于一个
	*/

	var ch chan int = make(chan int)
	go func() {
		time.Sleep(3 * time.Second)
		close(ch)
	}()

	// loop:
	for i, start := 0, time.Now(); i < 10; i++ {
		fmt.Println("time past:", time.Since(start))
		select {
		// case r, ok := <-ch:
		// 	if !ok {
		// 		fmt.Println("looks like the channel has been closed")
		// 		break loop
		// 	}
		// 	fmt.Println(r, ok)
		case r := <-ch: //close之后居然会返回0,看来获取的是 channel内值的零值,所以如果你预测到你的channel有关闭的可能,最好就报错双返回值的方式获取,手动break

			fmt.Println("NO!", r)
		default:
			fmt.Println("明月照沟渠")
		}
		time.Sleep(500 * time.Millisecond)
	}

	/*
		关于os线程和goroutine的区别
		1. 线程都有固定大小的内存块(一般是2MB)作为栈,而goroutine开始其生命周期的时候一般从2KB开始,可以根据需要动态伸缩,最大会达到1G,从而实现更多的自由度
		2.  os线程是被操作系统内核所调度的,也即是每几毫秒操作系统内核会中端处理器,通过scheduler挂起当前执行的线程并将其寄存器内容保存到内存中,检查线程列表决定下一个运行的线程,
			从内存中回复该线程的寄存器信息,回复此线程的现场并开始执行,这就是线程调度中的上下文切换,简单点就是前线程的寄存器保存到内存,后线程的内存恢复到寄存器,从而开始执行,这个
			上下文切换很慢,涉及到内存访问(计算器>内存>外存(硬盘)),并且增加了cpu的运行周期;
			而GO呢,这个调度器使用了一些技术手段，比如m:n调度，因为其会在n个操作系统线程上多工(调度)m个goroutine。Go调度器的工作和内核的调度是相似的，但是这个调度器只关注单独的
			Go程序中的goroutine（译注：按程序独立）。
			和操作系统的线程调度不同的是，Go调度器并不是用一个硬件定时器，而是被Go语言“建筑”本身进行调度的。例如当一个goroutine调用了time.Sleep，或者被channel调用或者mutex操
			作阻塞时，调度器会使其进入休眠并开始执行另一个goroutine，直到时机到了再去唤醒第一个goroutine。因为这种调度方式不需要进入内核的上下文，所以重新调度一个goroutine比调
			度一个线程代价要低得多。(所以更像是协程咯?)

		GOMAXPROCS
			m:n中的n,一般是系统CPU的核心数
	*/

}

var loadIconsOnce sync.Once
var icons map[string]image.Image

// Concurrency-safe.
func Icon(name string) image.Image {
	loadIconsOnce.Do(loadIcons)
	return icons[name]
}
func loadIcons() {
	icons = map[string]image.Image{
		"spades.png":   loadIcon("spades.png"),
		"hearts.png":   loadIcon("hearts.png"),
		"diamonds.png": loadIcon("diamonds.png"),
		"clubs.png":    loadIcon("clubs.png"),
	}
}
func loadIcon(name string) image.Image { return image.Black }
