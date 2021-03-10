package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

type iniSet struct {
	x, y int
}

//implent String()string
func (inisetP *iniSet) String() string {
	return fmt.Sprintf("x is %v,y is %v", (*inisetP).x, inisetP.y)
}

func main() {

	/*
		有一点需要注意,receiver 的func我们有语法糖,无论receiver是T 类型或是*T类型,编译器允许我们使用T/*T作为 *T/T的receiver,
		当然*T类型的receiver要求调用对象是一个能够寻址的对象,例如:
		type iniSet struct{...}
		func (*iniSet)String()string
		var _= iniSet{}.String()// compile error: String requires *IntSet receiver

		但是在interface这一章节中,我们就要意识到尽管iniset类型可以调用String()string方法,这不意味着iniSet类型implent这个方法,也就意味着对于
		type  fmt.Stringer{
			String() string
		}
		var t  fmt.Stringer
		var s IntSet
		var _ = s.String() // OK: s is a variable and &s has a String method\

		var _ fmt.Stringer = &s // OK
		var _ fmt.Stringer = s  // compile error: IntSet lacks String method

	*/
	var iniset1 = iniSet{1, 2}
	fmt.Println(iniset1.String())

	// fmt.Println((iniSet{3,4}).String())//compile error:cannot call pointer method on iniSet{...} 使用语法糖失败
	fmt.Println((&iniSet{3, 4}).String()) //what 居然可以 ，说好的不可取址的调用对象会失败呢,和上面一行比较起来,明明就是编译器自己不会主动帮忙取址的意思吗,还是说在编译期,对于literal
	// word是无法取址的意思

	// var _ fmt.Stringer =iniset1 //iniSet does not implement fmt.Stringer (String method has pointer receiver) 可以理解
	var _ fmt.Stringer = &iniset1 //Yes!
	// var _ fmt.Stringer = iniSet{1, 2}//iniSet does not implement fmt.Stringer (String method has pointer receiver) 其实就是leteral word 无法再编译期通过,看起来单纯的字面值无法取址
	var _ fmt.Stringer = &iniSet{1, 3}

	/*
		通过改变接口类型,可以改变暴露的方法,妙
	*/

	os.Stdout.Write([]byte("hello")) // OK: *os.File has Write method
	os.Stdout.Close()                // OK: *os.File has Close method

	var w io.Writer
	w = os.Stdout
	w.Write([]byte("hello")) // OK: io.Writer has Write method
	// w.Close()                // compile error: io.Writer lacks Close method

	/*
		空接口 interface{} 对实现它的类型没有要求，所以我们可以将任意一个值赋给空接口类型
		但是
		对于创建的一个interface{}值持有一个boolean，float，string，map，pointer，或者任意其它的类型；我们当然不能直接对它持有的值做操作，因为interface{}没有任何方法
		那空接口的作用是什么呢,确实好多函数的形参类型被设置为空接口,可以接受任何类型的参数
	*/
	var any interface{}
	any = true
	any = 12.34
	any = "hello"
	any = map[string]int{"one": 1}
	any = new(bytes.Buffer)
	fmt.Println(any)
	/*
		每一个具体类型的组基于它们相同的行为可以表示成一个接口类型。不像基于类的语言，他们一个类实现的接口集合需要进行显式的定义，
		在Go语言中我们可以在需要的时候定义一个新的抽象或者特定特点的组，而不需要修改具体类型的定义。当具体的类型来自不同的作者时
		这种方式会特别有用。当然也确实没有必要在具体的类型中指出这些共性。
		这一段话很有意思,确实其他语言的接口和Go的接口差别较大
	*/
}
