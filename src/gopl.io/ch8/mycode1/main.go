package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gopl.io/ch8/thumbnail"
)

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
		不对啊,感觉单纯的一个单项channel,有啥用呢？
		其实在申明的时候还是双向channel,只是在部分函数中只是作为单向channel使用,其实是编译期行为,函数中copy channel产生隐式的类型转换,但切记没有反向的隐式转换,即不能把一个单向channel
		转换成双向channel

		为了表明这种意图并防止被滥用，Go语言的类型系统提供了单方向的channel类型，分别用于只发送或只接收的channel。类型chan<- int表示一个只发送int的channel，
		只能发送不能接收。相反，类型<-chan int表示一个只接收int的channel，只能接收不能发送。（箭头<-和关键字chan的相对位置表明了channel的方向。）这种限制将在编译期检测。
	*/

	/*
		buffered channel 缓存channel 感觉和无缓存channel没什么区别啊,无非是阻塞发生在无缓存时罢了,无缓存channel是带缓存channel的一个特例,即使说其每次都处于缓存队列已满条件下
		缓存队列是FIFO
	*/

	/*
		测试一下双返回值的channel什么时候结束
		有以下代码可知,双参数返回值只会在在channel close关闭之后才会err的,并不是block就会err,需分清

	*/

	chan1 := make(chan int)
	go func() {
		chan1 <- 1
		time.Sleep(time.Second * 3)
		close(chan1)
	}()

	// res, err := <-chan1
	// fmt.Println(res, err)
	// res, err2 := <-chan1
	// fmt.Println(res, err2)
	for res := range chan1 { //range是不会获取到err的,其在真正遇到err的时候就是iterate结束的时候
		fmt.Println(res) //所以在输出 1之后 sleep 3s main goroutine return
	}

	/*
		下面是示例代码 ,分析一下文中的几个观点
		1.为什么wg.add(1)要放在main goroutine而非go routine中
			answer:如果add操作放在 go中,那么会出现什么情况呢,很有可能wg在某一瞬间在range filenames之前就达到了吊诡的wait条件,这么说可以理解吧,add(1) Done add(1) Done 。。。。。add(1)
			在第二次Done时就已经 打到wait()条件了
			当然,我觉得可以直接add(len(filenames))
		2:为什么 close()要放在 go routine中而非main goroutine中
			answer:实际上我们的函数有三类go routine,一个main,一类多个 worker routine,一个 close routine，其中,workers和closer因为wg的缘故是有消息通讯的,那么当workers全部完工之后,workers和
			main则因为unbuffered channel也是和main 保持阻塞等待关系的,也即是说workers结束之时,total计算正确之日,因此此时closer操作是安全且正确的
			那么如果closer在main中,考虑其放在size range之后,Close在range之后，但是range一来close才能结束循环,故close()无法抵达;放在range之前,
			则第一个worker就阻塞在其中了.wait永远达成条件

	*/

	makeThumbnails6 := func(filenames <-chan string) int64 {
		sizes := make(chan int64)
		var wg sync.WaitGroup // number of working goroutines
		for f := range filenames {
			wg.Add(1)
			// worker
			go func(f string) {
				defer wg.Done()
				thumb, err := thumbnail.ImageFile(f)
				if err != nil {
					log.Println(err)
					return
				}
				info, _ := os.Stat(thumb) // OK to ignore error
				sizes <- info.Size()
			}(f)
		}

		// closer
		go func() {
			wg.Wait()
			close(sizes)
		}()

		var total int64
		for size := range sizes {
			total += size
		}
		return total
	}
	fmt.Printf("%T", makeThumbnails6)
}
