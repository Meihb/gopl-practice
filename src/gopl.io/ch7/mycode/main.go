package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"syscall"
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
	// os.Stdout.Close()                // OK: *os.File has Close method

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

	/*
		接口值 由两个部分组成,具体的类型和类型值,即动态类型和动态值
		接口值也会被零值化, 其type和value 都是nil
	*/
	var w1 io.Writer
	fmt.Println("w1", w1 == nil) //true
	// w1.Write([]byte("hello")) // panic: nil pointer dereference    panic: nil pointer dereference 零值化的接口是不能直接调用方法的,但是编译期不会报错,而是运行时报错

	//赋值动态类型
	w1 = os.Stdout //此时,w1.type=*os.File,w1.value 是os.Stdout的拷贝

	/*口类型是非常与众不同的。其它类型要么是安全的可比较类型（如基本类型和指针）要么是完全不可比较的类型（如切片，映射类型，和函数），
	但是在比较接口值或者包含了接口值的聚合类型时，我们必须要意识到潜在的panic。同样的风险也存在于使用接口作为map的键或者switch的操作数。
	只能比较你非常确定它们的动态值是可比较类型的接口值。
	*/
	var w3 io.Writer
	fmt.Printf("%T\n", w3) // "<nil>"
	w3 = os.Stdout
	fmt.Printf("%T\n", w3) // "*os.File"
	w3 = new(bytes.Buffer)
	fmt.Printf("%T\n", w3) // "*bytes.Buffer"

	/*
		这个sort.Interface 接口有意思啊,尤其
		package sort

		type reverse struct{ Interface } // that is, sort.Interface

		func (r reverse) Less(i, j int) bool { return r.Interface.Less(j, i) }

		func Reverse(data Interface) Interface { return reverse{data} }

		func Sort(data Interface) {n := data.Len()	quickSort(data, 0, n, maxDepth(n))}
		这个就是 override重写吧 见下面 ti interface
		感觉更像是继承呀
		妙呀妙呀
	*/

	d1 := data1{}
	fmt.Println((ti(d1).string1(2, 3))) // string x is 2,y is 3
	fmt.Printf("%T\n", tiPrime{ti(d1)})
	fmt.Println(tiPrime{ti(d1)}.string1(2, 3)) //string x is 3,y is 2

	/*
		类型断言 alert?
		对于x.(T),我们称之为断言类型
		x是一个接口类型
		f,ok:=x.(T)
		如果T是一个具体类型,则断言检查x的动态类型是否和T相同,如果是,f是x的动态值,当然类型是T，两者相同嘛;如果不是,panic或者false(f为零值)
		如果T是一个接口类型,则断言检查x的动态类型是否满足T接口,如果是,动态值没有获取到；这个结果仍然是一个有相同动态类型和值部分的接口值，
		但是结果为类型T(把x本来暴露的method改成T可以暴露的method);如果不是和第一条一样

		所以我的结论就是,断言 无论T是具体类型还是接口类型,如果成功,返回的依然是自己本身,只是在编译期可能改变了method set,如果失败,
		返回的动态类型值将会是具体类型(具体类型情况)或者nil(接口类型),动态值都是nil

		感觉可以理解,接口类型这玩意好像是compile独有的,runtime还在乎个屁啊
		另外有一个关键的点,如果判断成功的话,返回值是指向同一个地址的,这点非常棒？存疑,也可能因为我的测试数据是pointer类型才如此
	*/

	var w4 io.Writer
	fmt.Printf("w4 :%T %[1]v  \n", w4) //nil nil 符合预期,这是一个动态类型和接口值都为nil的零值,但是runtime其实并没有记录w4是个Writer接口,所以接口判断是编译期行为的吧
	w4 = os.Stdout
	fmt.Printf("w4 :%T %[1]v  \n", w4)  //*os.File &{0xc0000da280} 所以其实并没有记录io.Writer的字段,毕竟interface只包含一个动态类型和接口值,好呀,原来只是编译期语法判断行为
	f, ok := w4.(*os.File)              // success: f == os.Stdout
	fmt.Printf("%T %[1]v %v \n", f, ok) //*os.File &{0xc0000da280} true
	c, ok := w4.(*bytes.Buffer)         // panic: interface holds *os.File, not *bytes.Buffer
	fmt.Printf("%T %[1]v %v \n", c, ok) //*bytes.Buffer <nil> false 对于具体类型的断言,失败后动态类型是T,而接口值是nil

	// w4.Read() 接口未暴露
	rw, ok := w4.(io.ReadWriter)         // success: *os.File has both Read and Write
	fmt.Printf("%T %[1]v %v \n", rw, ok) //*os.File &{0xc0000da280} true 动态类型没有发生变化诶,奇怪奇怪
	p1 := []byte{1, 2}
	rw.Read(p1)
	rw2, ok := w4.(ti)
	fmt.Printf("%T %[1]v %v \n", rw2, ok) //nil nil false
	rw3, ok := w4.(flag.Getter)
	fmt.Printf("rw3 :%T %[1]v %v \n", rw3, ok) //nil nil false 果然,对于接口类型情况而言,失败后类型是nil,接口值也是nil

	var ErrNotExist = errors.New("file does not exist")

	// IsNotExist returns a boolean indicating whether the error is known to
	// report that a file or directory does not exist. It is satisfied by
	// ErrNotExist as well as some syscall errors.
	isNotExist := func(err error) bool {
		pe, ok := err.(*os.PathError)     //这是具体类型而非接口类型,知道了吧
		fmt.Printf("pe :%T %[1]v \n", pe) //果不其然,返回的是 *fs.PathError nil
		if ok {
			err = pe.Err
			fmt.Println("error:", err)
		}
		return err == syscall.ENOENT || err == ErrNotExist
	}
	// _, err := os.Open("/file")
	_, err := os.Hostname()
	fmt.Println("end:", isNotExist(err)) // "true"

	myBytes := []byte("newsd美") //什么呀,这是类型转换,你以为是字面值的第二种写法呢
	fmt.Println(myBytes)
}

//定义一个ti接口
type ti interface {
	string1(x, y int) string //创建一个输出x,y的的接口规范‘
}
type tiPrime struct {
	ti //我们创建一个内嵌ti接口的结构
}
type data1 struct{}

func (tiPrime tiPrime) string1(x, y int) string {
	return tiPrime.ti.string1(y, x)
}

func (data1 data1) string1(x, y int) string {
	return fmt.Sprintf(" string x is %v,y is %v", x, y)
}
