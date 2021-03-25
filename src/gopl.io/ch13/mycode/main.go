package main

import (
	"fmt"
	"unsafe"
)

func main() {
	fmt.Println("unsafe包是一个采用特殊方式实现的包,实际上是有编译器实现的,一般和操作系统密切配合")

	/*unsafe.Sizeof, Alignof 和 Offsetof*/
	//unsafe.Sizeof函数返回操作数在内存中的字节大小，参数可以是任意类型的表达式，但是它并不会对表达式进行求值
	fmt.Println(unsafe.Sizeof(float64(0))) // "8" 8byte
	fmt.Println(unsafe.Sizeof(3 + 2))      //8 可是这里还是计算了呀
	a, b := 3, 2
	fmt.Println(unsafe.Sizeof(a + b)) //8 所以这个不会对表达式求值到底是什么意思呢？
	//考虑到可移植性，引用类型或包含引用类型的大小在32位平台上是4个字节，在64位平台上是8个字节。 这句话很容易理解,记录一下一面忘记

	/*
		Alignof 内存地址对齐 在之前已经了解过相当一些知识
		不过本书说go在其编书的时候struct还没能优化的随意排列接口提内每个字段的内存位置,从而做出了一些举例
		另外,bool 1字节,intN、uintN...complex n/8字节 ,int、uint、uintptr 一个机器字(机器字就是计算机机器字长长度的位)
		*T 一个机器字
		string 2个机器字(data,len)
		[]T 三个机器字(data,len,cap)
		map 一个机器字
		func 1个机器字
		chan 1个机器字
		interface 2个机器字(type,value)
		                               // 64-bit  32-bit
		struct{ bool; float64; int16 } // 3 words 4words
		struct{ float64; int16; bool } // 2 words 3words
		struct{ bool; int16; float64 } // 2 words 3words

		var x struct {
			a bool
			b int16
			c []int
		}
		Sizeof(x)   = 32  Alignof(x)   = 8
		Sizeof(x.a) = 1   Alignof(x.a) = 1 Offsetof(x.a) = 0
		Sizeof(x.b) = 2   Alignof(x.b) = 2 Offsetof(x.b) = 2
		Sizeof(x.c) = 24  Alignof(x.c) = 8 Offsetof(x.c) = 8
	*/

	/*
		unsafe.Pointer是特别定义的一种指针类型（译注：类似C语言中的void*类型的指针），它可以包含任意类型变量的地址。
	*/
	var float64bits = func(f float64) uint64 { return *(*uint64)(unsafe.Pointer(&f)) }
	/*
		一个普通的*T类型指针可以被转化为unsafe.Pointer类型指针，并且一个unsafe.Pointer类型指针也可以被转回普通的指针，被转回普通的指针类型并不需要和原始的*T类型相同
		这不就是转换了TYPE吗,话说这和类型转换有什么区别呀

		看起来类型转换是经过一些逻辑运算的,而通过unsafe.Pointer是保留了原bit位不变直接算成新的类型的
	*/
	fmt.Printf("%#016x\n", 1.0)              //00000x1.0000p+00
	fmt.Printf("%#016x\n", float64bits(1.0)) //0x3ff0000000000000
	fmt.Printf("%#016x\n", uint(1.0))        //0x0000000000000001

	/*
		通过转为新类型指针，我们可以更新浮点数的位模式。通过位模式操作浮点数是可以的，但是更重要的意义是指针转换语法让我们可以在不破坏类型系统的前提下向内存写入任意的值。git
	*/

}
