package main

import (
	"fmt"
	"os"
)

func main() {
	var a int = 5

	var t1 func() func() = func() func() {
		a++
		fmt.Printf("a in t1:%v \n", a)
		return func() {
			a++
			fmt.Printf("a in defer:%v\n", a)
		}
	}

	defer t1() //两者区别,defer后面跟的是func()而非func,这一点需切记,因此可以在defer中实现例子中trace那样在defer时执行首尾记录操作的代码
	//defer 遇到func()func()()()这类当然只会执行最后一个函数作为defer函数入栈,斯巴拉西,有点像python的decorator?记不太清概念了
	defer t1()()
	fmt.Println(a)

	/*
		defer语句中的函数会在return语句更新返回值变量后再执行，又因为在函数中定义的匿名函数可以访问该函数包括返回值变量在内的所有变量，
		所以，对匿名函数采用defer机制，可以使其观察函数的返回值。
	*/

	var double func(x int) (result int) = func(x int) (result int) {
		defer func() { fmt.Printf("double(%d) = %d\n", x, result) }()
		return x + x
	}

	var triple func(x int) (result int) = func(x int) (result int) {
		defer func() { result += x }()
		return double(x)
	}
	double(2)
	res := triple(2)
	fmt.Println(res) //6 我靠,defer确实在return后面才执行的,但不以为这return结果不会被修改,
	//被调用函数内return是没来的及变化的,但是caller会体现出变化

	//defer在当前func return时才触发,所以在循环中若有必要可以添加新的一层 词法域,提前触发
	var filenames []string
	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			fmt.Println(err)
		}
		//do something
		defer f.Close() //观察词法域,f.close()在循环结束之后才会结束,这就意味着句柄一直积压,可能出现系统文件描述符用光的问题
	}
	//优化
	var doFile func(filename string) error = func(filename string) error {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		//do something
		defer f.Close() //在这里添加defer,则每一个句柄都在一个函数内,可以保证外循环体中只有一个句柄存在
	}
	for _, filename := range filenames {
		err := doFile(filename)
		if err != nil {
			fmt.Println(err)
		}
	}

	/*
		panic runtime error
		when it occurs,cutting down current goroutine,implent defer functions,breakdown and log
	*/

}
