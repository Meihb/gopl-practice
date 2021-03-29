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
		通过转为新类型指针，我们可以更新浮点数的位模式。通过位模式操作浮点数是可以的，但是更重要的意义是指针转换语法让我们可以在不破坏类型系统的前提下向内存写入任意的值。

		一个unsafe.Pointer指针也可以被转化为uintptr类型，然后保存到指针型数值变量中（译注：这只是和当前指针相同的一个数字值，并不是一个指针），然后用以做必要的指针数值运算。
		（第三章内容，uintptr是一个无符号的整型数，足以保存一个地址）这种转换虽然也是可逆的，但是将uintptr转为unsafe.Pointer指针可能会破坏类型系统，因为并不是所有的数字都
		是有效的内存地址。

		许多将unsafe.Pointer指针转为原生数字，然后再转回为unsafe.Pointer类型指针的操作也是不安全的
		有下面示例就知道原因了,内存不是连续的,有很多空洞,所以不是所有原声数字转为Pointer的操作都是安全的这句话自然是可以理解的 NO!
		不是这样解释的,就算使用offset函数,如果你引入了uintptr的中间变量,也可能出现安全性问题,具体原因如下,GC相关
	*/
	var x struct { //align 8 (64bit system)
		a bool  //align 1 offset 0
		b int16 //align 2 offset 2  因为他和a实在同一行机器字中,顾按照最小公倍数来计算,具体自己再去看一下网页内容
		c []int //align 8 offset 8
	}

	// 和 pb := &x.b 等价
	pb := (*int16)(unsafe.Pointer(
		uintptr(unsafe.Pointer(&x)) +
			unsafe.Offsetof(x.b)))
	*pb = 42
	fmt.Println(x.b) // "42"

	p1 := unsafe.Pointer(&x)
	fmt.Println(p1, unsafe.Alignof(x)) //0xc0000443e0 8
	/*
		uintptr 是一个足够容纳当前环境下指针长度的uint类型，但是他的类型转换是 uintptr,这就很难接受了,不是我喜欢的(uintptr)这样的强制转换
		另外,unsafet.Pointer也是如此类型转换的,挠头
		你说尼玛呢,go的类型转换本就是如此  type(v)
	*/
	fmt.Println(uintptr(p1))

	// NOTE: subtly incorrect!
	/*
		不要试图引入一个uintptr类型的临时变量，因为它可能会破坏代码的安全性 重点！
		微妙的错误,微妙的解释
		产生错误的原因很微妙。有时候垃圾回收器会移动一些变量以降低内存碎片等问题。这类垃圾回收器被称为移动GC。当一个变量被移动，所有的保存该变量旧地址
		的指针必须同时被更新为变量移动后的新地址。从垃圾收集器的视角来看，一个unsafe.Pointer是一个指向变量的指针，因此当变量被移动时对应的指针也必须
		被更新；但是uintptr类型的临时变量只是一个普通的数字，所以其值不应该被改变。上面错误的代码因为引入一个非指针的临时变量tmp，导致垃圾收集器无法
		正确识别这个是一个指向变量x的指针。当第二个语句执行时，变量x可能已经被转移，这时候临时变量tmp也就不再是现在的&x.b地址。第三个向之前无效地址空
		间的赋值语句将彻底摧毁整个程序

		另有一些错误
		pT := uintptr(unsafe.Pointer(new(T)))
		很明显,new(T)并没有赋值给一个变量来应用此地址,如此GC将会回收该内存空间,那么pt将会是无效的地址

		因此目前的指导原则是,原子性! 将涉及到uintptr和unsafe.Pointer的互相转换尽可能减少并放在同一个表达式中

	*/
	tmp := uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)
	pb2 := (*int16)(unsafe.Pointer(tmp))
	*pb2 = 42
	fmt.Println(x.b) // "42"

}
