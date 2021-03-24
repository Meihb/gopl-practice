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
		另外,bool 1字节,intN、uintN...complex n/8字节
	*/
}
