package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
)

func main() {
	/*
		反射是由 reflect 包提供的。 它定义了两个重要的类型, Type 和 Value.
		 一个 Type 表示一个Go类型. 它是一个接口, 有许多方法来区分类型以及检查它们的组成部分, 例如一个结构体的成员或一个函数的参数等. 唯一能反映 reflect.Type 实现的是接口的类型描述信息,
		  也正是这个实体标识了接口值的动态类型.
		一个 reflect.Value 可以装载任意类型的值. 函数 reflect.ValueOf 接受任意的 interface{} 类型, 并返回一个装载着其动态值的 reflect.Value. 和 reflect.TypeOf 类似,
		reflect.ValueOf 返回的结果也是具体的类型, 但是 reflect.Value 也可以持有一个接口值

		两者都实现了fmt.Stringer 接口
	*/
	t := reflect.TypeOf(3)  // a reflect.Type
	fmt.Println(t.String()) // "int"
	fmt.Println(t)          // "int"

	//将一个具体的值转为接口类型会有一个隐式的接口转换操作, 它会创建一个包含两个信息的接口值: 操作数的动态类型(这里是int)和它的动态的值
	var w io.Writer = os.Stdout
	fmt.Println(reflect.TypeOf(w)) // "*os.File"

	//其实Printf 中%T参数正式 使用refleact.TypeOf来实现的,对应的%V就是reflect.ValueOf()实现的
	is := fmt.Stringer(reflect.TypeOf(3))
	fmt.Println(is)

	v := reflect.ValueOf(3)                      // a reflect.Value
	fmt.Println(v)                               // "3"
	fmt.Printf("%v\n", v)                        // "3"
	fmt.Println(v.String())                      // NOTE: "<int Value>" String()返回的依然是类型(除非v是字符串类型)
	fmt.Println(reflect.ValueOf("sss").String()) //sss

	//对 Value 调用 Type 方法将返回具体类型所对应的 reflect.Type:
	t1 := v.Type()           // a reflect.Type
	fmt.Println(t1.String()) // "int"

	fmt.Printf("%T %[1]v \n", v)

	//reflect.ValueOf 的逆操作是 reflect.Value.Interface 方法. 它返回一个 interface{} 类型，装载着与 reflect.Value 相同的具体值:
	v1 := reflect.ValueOf(3) // a reflect.Value
	x := v1.Interface()      // an interface{}
	i := x.(int)             // an int  这个断言是类型2,T是接口类型,所以断言成功后结果将会改变方法集合 这是我们实现从空接口转为具体类型的方法,断言
	fmt.Printf("%d\n", i)    // "3"

	/*
		我们使用 reflect.Value 的 Kind 方法来替代之前的类型 switch. 虽然还是有无穷多的类型, 但是它们的kinds类型却是有限的: Bool, String 和 所有
		数字类型的基础类型; Array 和 Struct 对应的聚合类型; Chan, Func, Ptr, Slice, 和 Map 对应的引用类型; interface 类型; 还有表示空值的Invalid
		类型. (空的 reflect.Value 的 kind 即为 Invalid.)
	*/
	
}
