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
		观看实例
	*/

	/*
		通过reflect.Value修改值
		 func (reflect.Value).Elem() reflect.Value  对于指针对象的reflect.Value可以使用 Elem() method来解引用,等效于对于指针类型变量作*操作

		Type.Elem()  	Elem returns the value that the interface v contains or that the pointer v points to. It panics if v's Kind is not Interface or Ptr.
		 				It returns the zero Value if v is nil.
		Value.Elem() 	 Elem returns a type's element type.It panics if the type's Kind is not Array, Chan, Map, Ptr, or Slice.

		一个变量是一个可寻址的内存空间,里面存储了一个值，并且存储的值可以通过内存地址来更新。

		下面这个例子中,x是个简单类型而非引用类型,这点要先明确
		a,b不可取址很明显,a中的值仅仅是整数2的拷贝副本,他的内存地址中直接存的应该就是2吧
		c也不可取址,c是一个指针&x的拷贝,当你访问c时获取的是x的值,那么这个值你还能寻址吗,很明显不能,
		d为什么可以寻址,d=*c,它是c的解引用方式生成的，指向另一个变量，因此是可取地址的

		!实际上，所有通过reflect.ValueOf(x)返回的reflect.Value都是不可取地址的
		其实完全可以 用是否从ptr类型获取reflect.Value来看
	*/
	x = 2                                                                        // value   type    variable?
	a := reflect.ValueOf(2)                                                      // 2       int     no
	b := reflect.ValueOf(x)                                                      // 2       int     no
	c := reflect.ValueOf(&x)                                                     // &x      *int    no
	d := c.Elem()                                                                // 2       int     yes (x)
	e := d.Elem()                                                                // 2 		int 	no
	fmt.Println(x, &x, a, b, c, d, e)                                            //x=2 &x=0xc000040280 a=2 b=2 c=0xc000040280 d=2 e=2
	fmt.Println(a.CanAddr(), b.CanAddr(), c.CanAddr(), d.CanAddr(), e.CanAddr()) //false false false true false 判断是否可以取址

	//测试对于 interface 和ptr 以及其他的Elem method
	type elemTst struct{ int }
	cat1 := &elemTst{1}
	ValueOfCat1 := reflect.ValueOf(cat1)
	typeOfCat1 := reflect.TypeOf(cat1)
	fmt.Println(ValueOfCat1, typeOfCat1, ValueOfCat1.Elem(), typeOfCat1.Elem()) // &{1} *main.elemTst {1} main.elemTst

	// var w1 io.Reader
	// v1OfValue := reflect.ValueOf(w1)
	// fmt.Println(v1OfValue, v1OfValue.Elem()) //panic: reflect: call of reflect.Value.Elem on zero Value

	// i1 := 1
	// i1OfValue := reflect.ValueOf(i1)
	// fmt.Println(i1OfValue, i1OfValue.Elem()) //panic: reflect: call of reflect.Value.Elem on int Value

	/*
		从可曲直的reflect.Value来访问变量需要三个步骤
		1.调用Addre()method,返回一个Value,保存了指向变量的指针
		2.调用interface()方法,从而返回一个interface,里面包含指向变量的指针
		3.类型断言

		不对啊,现在断言还想不能够呀
		靠可以,是因为x在上面被我定义为interface{}

		conclude:
		值可修改条件之一:可被寻址
		之二:可导出 大写！
	*/
	x1 := 2
	d = reflect.ValueOf(&x1).Elem()       // d refers to the variable x
	px, ok := d.Addr().Interface().(*int) // px := &x
	fmt.Println(ok)
	if ok {
		*px = 3             // x = 3
		fmt.Println(px, x1) // "3"
	}
	d.SetInt(4) // 或者直接Set family methods
	fmt.Println(x1)

}
